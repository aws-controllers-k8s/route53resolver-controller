# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Integration tests for the Route53 ResolverRuleAssociation resource
"""

import boto3
import logging
import time
from typing import Dict

import pytest

from acktest.k8s import resource as k8s
from acktest.k8s import condition
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_route53resolver_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources

RESOURCE_PLURAL = "resolverruleassociations"

# Time to wait after modifying the CR for the status to change
MODIFY_WAIT_AFTER_SECONDS = 10

# Time to wait for association to become COMPLETE
ASSOCIATION_WAIT_SECONDS = 30


def create_forward_rule(route53resolver_client) -> str:
    """Create a FORWARD resolver rule via AWS API for testing associations."""
    rule_name = random_suffix_name("assoc-test-rule", 32)
    response = route53resolver_client.create_resolver_rule(
        CreatorRequestId=rule_name,
        DomainName="test-assoc.example.com",
        Name=rule_name,
        RuleType="SYSTEM",
    )
    rule_id = response["ResolverRule"]["Id"]
    return rule_id


def delete_rule(route53resolver_client, rule_id: str):
    """Delete the resolver rule created for testing."""
    try:
        route53resolver_client.delete_resolver_rule(ResolverRuleId=rule_id)
    except route53resolver_client.exceptions.ResourceNotFoundException:
        pass


@pytest.fixture
def resolver_rule_association(route53resolver_client):
    """Create a single ResolverRuleAssociation CR and yield it for testing."""
    rule_id = create_forward_rule(route53resolver_client)
    vpc_id = get_bootstrap_resources().ResolverEndpointVPC.vpc_id

    association_name = random_suffix_name("rule-assoc", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["RESOLVER_RULE_ASSOCIATION_NAME"] = association_name
    replacements["RESOLVER_RULE_ID"] = rule_id
    replacements["VPC_ID"] = vpc_id

    resource_data = load_route53resolver_resource(
        "resolver_rule_association",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
        association_name, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr, rule_id)

    # Cleanup
    try:
        _, deleted = k8s.delete_custom_resource(ref, 3, 10)
        assert deleted
    except:
        pass

    time.sleep(MODIFY_WAIT_AFTER_SECONDS)
    delete_rule(route53resolver_client, rule_id)


@pytest.fixture
def two_vpc_associations(route53resolver_client):
    """Create two ResolverRuleAssociation CRs for the same rule but different VPCs."""
    rule_id = create_forward_rule(route53resolver_client)

    # Use the bootstrap VPC's subnets to derive two VPC IDs
    # The bootstrap creates one VPC; for the second we create a temporary one
    vpc_id_1 = get_bootstrap_resources().ResolverEndpointVPC.vpc_id

    ec2_client = boto3.client("ec2")
    vpc_response = ec2_client.create_vpc(CidrBlock="10.99.0.0/16")
    vpc_id_2 = vpc_response["Vpc"]["VpcId"]

    # Wait for VPC to be available
    ec2_client.get_waiter("vpc_available").wait(VpcIds=[vpc_id_2])

    association_name_1 = random_suffix_name("rule-assoc-1", 32)
    association_name_2 = random_suffix_name("rule-assoc-2", 32)

    refs = []
    for assoc_name, vpc_id in [(association_name_1, vpc_id_1), (association_name_2, vpc_id_2)]:
        replacements = REPLACEMENT_VALUES.copy()
        replacements["RESOLVER_RULE_ASSOCIATION_NAME"] = assoc_name
        replacements["RESOLVER_RULE_ID"] = rule_id
        replacements["VPC_ID"] = vpc_id

        resource_data = load_route53resolver_resource(
            "resolver_rule_association",
            additional_replacements=replacements,
        )

        ref = k8s.CustomResourceReference(
            CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
            assoc_name, namespace="default",
        )
        k8s.create_custom_resource(ref, resource_data)
        cr = k8s.wait_resource_consumed_by_controller(ref)
        assert cr is not None
        refs.append((ref, cr))

    yield (refs, rule_id, vpc_id_2)

    # Cleanup
    for (ref, _) in refs:
        try:
            _, deleted = k8s.delete_custom_resource(ref, 3, 10)
        except:
            pass

    time.sleep(MODIFY_WAIT_AFTER_SECONDS)
    delete_rule(route53resolver_client, rule_id)

    # Delete the temporary VPC
    try:
        ec2_client.delete_vpc(VpcId=vpc_id_2)
    except:
        pass


@service_marker
@pytest.mark.canary
class TestResolverRuleAssociation:
    def test_create_delete(self, route53resolver_client, resolver_rule_association):
        """Test basic create and delete of a ResolverRuleAssociation."""
        (ref, cr, rule_id) = resolver_rule_association

        # Verify the association ID was assigned
        association_id = cr["status"]["id"]
        assert association_id is not None

        # Wait for association to become COMPLETE
        time.sleep(ASSOCIATION_WAIT_SECONDS)
        cr = k8s.get_resource(ref)

        # Verify Synced condition
        assert condition.assert_synced(ref)

        # Verify the association exists in AWS
        aws_res = route53resolver_client.get_resolver_rule_association(
            ResolverRuleAssociationId=association_id
        )
        assert aws_res is not None
        assoc = aws_res["ResolverRuleAssociation"]
        assert assoc["ResolverRuleId"] == rule_id
        assert assoc["Status"] == "COMPLETE"

        # Delete the CR
        _, deleted = k8s.delete_custom_resource(ref, 3, 10)
        assert deleted

        # Wait for AWS to process deletion
        time.sleep(ASSOCIATION_WAIT_SECONDS)

        # Verify the association no longer exists in AWS
        try:
            route53resolver_client.get_resolver_rule_association(
                ResolverRuleAssociationId=association_id
            )
            pytest.fail(f"Association {association_id} should have been deleted")
        except route53resolver_client.exceptions.ResourceNotFoundException:
            pass

    def test_multiple_vpcs_same_rule(self, route53resolver_client, two_vpc_associations):
        """Test that the same resolver rule can be associated with multiple VPCs."""
        (refs, rule_id, vpc_id_2) = two_vpc_associations

        time.sleep(ASSOCIATION_WAIT_SECONDS)

        # Verify both associations exist and are COMPLETE
        for (ref, _) in refs:
            cr = k8s.get_resource(ref)
            association_id = cr["status"]["id"]
            assert association_id is not None

            aws_res = route53resolver_client.get_resolver_rule_association(
                ResolverRuleAssociationId=association_id
            )
            assoc = aws_res["ResolverRuleAssociation"]
            assert assoc["ResolverRuleId"] == rule_id
            assert assoc["Status"] == "COMPLETE"

        # Verify the two associations have different IDs and VPCs
        cr_1 = k8s.get_resource(refs[0][0])
        cr_2 = k8s.get_resource(refs[1][0])
        assert cr_1["status"]["id"] != cr_2["status"]["id"]
        assert cr_1["spec"]["vpcID"] != cr_2["spec"]["vpcID"]
