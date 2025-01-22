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

RESOURCE_PLURAL = "resolverendpoints"

# Time to wait after modifying the CR for the status to change
MODIFY_WAIT_AFTER_SECONDS = 10

# Time to wait after the zone has changed status, for the CR to update
CHECK_STATUS_WAIT_SECONDS = 10

@pytest.fixture
def resolver_endpoint():
    resolver_endpoint = random_suffix_name("resolver-endpoint", 32)
    security_group_id = get_security_group(get_bootstrap_resources().ResolverEndpointVPC.vpc_id)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["RESOLVER_NAME"] = resolver_endpoint
    replacements["DIRECTION"] = "OUTBOUND"
    replacements["SUBNET_1"] = get_bootstrap_resources().ResolverEndpointVPC.private_subnets.subnet_ids[0]
    replacements["SUBNET_2"] = get_bootstrap_resources().ResolverEndpointVPC.private_subnets.subnet_ids[1]
    replacements["SECURITY_GROUP"] = security_group_id
    replacements["DELETION_POLICY"] = "delete"


    resource_data = load_route53resolver_resource(
        "resolver_endpoint",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    # Create the k8s resource
    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
        resolver_endpoint, namespace="default",
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

@pytest.fixture
def resolver_endpoint_adopt():
    resolver_endpoint = random_suffix_name("resolver-endpoint", 32)
    security_group_id = get_security_group(get_bootstrap_resources().ResolverEndpointVPC.vpc_id)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["RESOLVER_NAME"] = resolver_endpoint
    replacements["DIRECTION"] = "OUTBOUND"
    replacements["SUBNET_1"] = get_bootstrap_resources().ResolverEndpointVPC.private_subnets.subnet_ids[0]
    replacements["SUBNET_2"] = get_bootstrap_resources().ResolverEndpointVPC.private_subnets.subnet_ids[1]
    replacements["SECURITY_GROUP"] = security_group_id
    replacements["DELETION_POLICY"] = "retain"


    resource_data = load_route53resolver_resource(
        "resolver_endpoint",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    # Create the k8s resource
    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
        resolver_endpoint, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    resolver_id = cr["status"]["id"]

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    #Delete with retain policy allows us to adopt same resource
    _, deleted = k8s.delete_custom_resource(ref, 3, 10)
    assert deleted  


    replacements["ADOPTION_POLICY"] = "adopt"
    replacements["ADOPTION_FIELDS"] = f"{{\\\"id\\\": \\\"{resolver_id}\\\"}}"
    resource_data = load_route53resolver_resource(
        "resolver_endpoint_adoption",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)
     # Adopt the k8s resource
    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
        resolver_endpoint, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)
    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr)

    # Try to delete, if doesn't already exist
    _, deleted = k8s.delete_custom_resource(ref, 3, 10)
    assert deleted


def get_security_group(vpc_id: str) -> str:
    ec2_client = boto3.client("ec2")
    filters = [{'Name': 'vpc-id', 'Values': [vpc_id]}]
    response = ec2_client.describe_security_groups(Filters=filters)
    return response['SecurityGroups'][0]['GroupId']

@service_marker
@pytest.mark.canary
class TestResolverEndpoint:
    def test_create_delete_public(self, route53resolver_client, resolver_endpoint):
        (ref, cr) = resolver_endpoint

        resolver_endpoint_id = cr["status"]["id"]

        assert resolver_endpoint_id

        try:
            aws_res = route53resolver_client.get_resolver_endpoint(ResolverEndpointId=resolver_endpoint_id)
            assert aws_res is not None
        except route53resolver_client.exceptions.ResourceNotFoundException:
            pytest.fail(f"Could not find Resolver Endpoint with ID '{resolver_endpoint_id}' in Route53")
    
    def test_adopt_delete(self, route53resolver_client, resolver_endpoint_adopt):
        (ref, cr) = resolver_endpoint_adopt

        assert 'status' in cr
        assert 'id' in cr['status']

        resolver_endpoint_id = cr["status"]["id"]

        assert resolver_endpoint_id

        try:
            aws_res = route53resolver_client.get_resolver_endpoint(ResolverEndpointId=resolver_endpoint_id)
            assert aws_res is not None
        except route53resolver_client.exceptions.ResourceNotFoundException:
            pytest.fail(f"Could not find Resolver Endpoint with ID '{resolver_endpoint_id}' in Route53")

