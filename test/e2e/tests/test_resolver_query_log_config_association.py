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

import logging
import time

import pytest

from acktest.k8s import resource as k8s
from acktest.k8s import condition
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_route53resolver_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources

RESOURCE_PLURAL = "resolverquerylogconfigassociations"

ASSOCIATION_TIMEOUT_SECONDS = 120


def create_query_log_config(route53resolver_client, bucket_name) -> str:
    config_name = random_suffix_name("qlca-test", 32)

    response = route53resolver_client.create_resolver_query_log_config(
        Name=config_name,
        DestinationArn=f"arn:aws:s3:::{bucket_name}",
    )
    config_id = response["ResolverQueryLogConfig"]["Id"]

    deadline = time.time() + 60
    while time.time() < deadline:
        res = route53resolver_client.get_resolver_query_log_config(
            ResolverQueryLogConfigId=config_id
        )
        if res["ResolverQueryLogConfig"]["Status"] == "CREATED":
            break
        time.sleep(5)

    return config_id


def delete_query_log_config(route53resolver_client, config_id):
    try:
        route53resolver_client.delete_resolver_query_log_config(
            ResolverQueryLogConfigId=config_id
        )
    except Exception as e:
        logging.warning(f"Failed to delete query log config {config_id}: {e}")


def wait_for_active(ref, timeout=ASSOCIATION_TIMEOUT_SECONDS):
    deadline = time.time() + timeout
    while time.time() < deadline:
        cr = k8s.get_resource(ref)
        if cr and cr.get("status", {}).get("status") == "ACTIVE":
            return cr
        time.sleep(10)
    pytest.fail(f"ResolverQueryLogConfigAssociation {ref.name} did not reach ACTIVE within {timeout}s")


@pytest.fixture
def resolver_query_log_config_association(route53resolver_client):
    resources = get_bootstrap_resources()
    bucket_name = resources.QueryLogBucket.name
    vpc_id = resources.ResolverEndpointVPC.vpc_id

    config_id = create_query_log_config(route53resolver_client, bucket_name)

    association_name = random_suffix_name("qlc-assoc", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["RESOLVER_QUERY_LOG_CONFIG_ASSOCIATION_NAME"] = association_name
    replacements["RESOLVER_QUERY_LOG_CONFIG_ID"] = config_id
    replacements["VPC_ID"] = vpc_id

    resource_data = load_route53resolver_resource(
        "resolver_query_log_config_association",
        additional_replacements=replacements,
    )

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
        association_name, namespace="default",
    )

    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)
    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr, config_id)

    try:
        if k8s.get_resource_exists(ref):
            k8s.delete_custom_resource(ref, 3, 10)
    except Exception as e:
        logging.warning(f"Cleanup failed for {association_name}: {e}")

    delete_query_log_config(route53resolver_client, config_id)


@service_marker
class TestResolverQueryLogConfigAssociation:
    @pytest.mark.canary
    def test_create_delete(self, route53resolver_client, resolver_query_log_config_association):
        (ref, cr, config_id) = resolver_query_log_config_association

        cr = wait_for_active(ref)

        association_id = cr["status"]["id"]
        assert association_id is not None

        condition.assert_synced(ref)

        aws_res = route53resolver_client.get_resolver_query_log_config_association(
            ResolverQueryLogConfigAssociationId=association_id
        )
        assoc = aws_res["ResolverQueryLogConfigAssociation"]
        assert assoc["ResolverQueryLogConfigId"] == config_id
        assert assoc["Status"] == "ACTIVE"

        _, deleted = k8s.delete_custom_resource(ref, 12, 10)
        assert deleted

        deleted_in_aws = False
        for _ in range(9):
            try:
                route53resolver_client.get_resolver_query_log_config_association(
                    ResolverQueryLogConfigAssociationId=association_id
                )
                time.sleep(10)
            except route53resolver_client.exceptions.ResourceNotFoundException:
                deleted_in_aws = True
                break

        assert deleted_in_aws
