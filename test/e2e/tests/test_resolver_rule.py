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

"""Integration tests for the Route53 ResolverEndpoint resource
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
from e2e.tests.test_resolver_endpoint import resolver_endpoint

RESOURCE_PLURAL = "resolverrules"

# Time to wait after modifying the CR for the status to change
MODIFY_WAIT_AFTER_SECONDS = 10

# Time to wait after the zone has changed status, for the CR to update
CHECK_STATUS_WAIT_SECONDS = 10


def create_resolver_endpoint():
    resolver_endpoint = random_suffix_name("resolver-endpoint", 32)
    security_group_id = get_security_group(get_bootstrap_resources().ResolverEndpointVPC.vpc_id)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["RESOLVER_NAME"] = resolver_endpoint
    replacements["DIRECTION"] = "OUTBOUND"
    replacements["SUBNET_1"] = get_bootstrap_resources().ResolverEndpointVPC.private_subnets.subnet_ids[0]
    replacements["SUBNET_2"] = get_bootstrap_resources().ResolverEndpointVPC.private_subnets.subnet_ids[1]
    replacements["SECURITY_GROUP"] = security_group_id


    resource_data = load_route53resolver_resource(
        "resolver_endpoint",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    # Create the k8s resource
    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, "resolverendpoints",
        resolver_endpoint, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield ref, cr

    # Try to delete, if doesn't already exist
    try:
        _, deleted = k8s.delete_custom_resource(ref, 3, 10)
        assert deleted
    except:
        pass


def get_security_group(vpc_id: str) -> str:
    ec2_client = boto3.client("ec2")
    filters = [{'Name': 'vpc-id', 'Values': [vpc_id]}]
    response = ec2_client.describe_security_groups(Filters=filters)
    return response['SecurityGroups'][0]['GroupId']


@pytest.fixture
def resolver_rule():
    resolver_rule = random_suffix_name("resolver-rule", 32)
    vpc_id = get_bootstrap_resources().ResolverEndpointVPC.vpc_id

    res_end = create_resolver_endpoint()
    for i in res_end:
        (ref, cr) = i

    resolver_endpoint_id = cr["status"]["id"]
    replacements = REPLACEMENT_VALUES.copy()
    replacements["RESOLVER_RULE_NAME"] = resolver_rule
    replacements["RESOLVER_RULE_DOMAIN"] = "abc.xyz1"
    replacements["RESOLVER_ENDPOINT_ID"] = resolver_endpoint_id
    replacements["RESOLVER_RULE_TYPE"] = "FORWARD"
    replacements["VPC_ID"] = vpc_id
    replacements["IP"] = "1.2.3.4"
    replacements["PORT"] = "53"


    resource_data = load_route53resolver_resource(
        "resolver_rule",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    # Create the k8s resource
    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
        resolver_rule, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr)

    # Try to delete, if doesn't already exist
    try:
        _, deleted = k8s.delete_custom_resource(ref, 3, 10)
        assert deleted
    except:
        pass

@service_marker
@pytest.mark.canary
class TestResolverRule:
    def test_create_delete_public(self, route53resolver_client, resolver_rule):
        (ref, cr) = resolver_rule

        resolver_rule_id = cr["status"]["id"]

        assert resolver_rule_id

        try:
            aws_res = route53resolver_client.get_resolver_rule(ResolverRuleId=resolver_rule_id)
            assert aws_res is not None
        except route53resolver_client.exceptions.ResourceNotFoundException:
            pytest.fail(f"Could not find Resolver Rule with ID '{resolver_rule_id}' in Route53")

