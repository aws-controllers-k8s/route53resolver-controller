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

RESOURCE_PLURAL = "resolverquerylogconfigs"

CREATE_TIMEOUT_SECONDS = 120


def wait_for_created(ref, timeout=CREATE_TIMEOUT_SECONDS):
    deadline = time.time() + timeout
    while time.time() < deadline:
        cr = k8s.get_resource(ref)
        if cr and cr.get("status", {}).get("status") == "CREATED":
            return cr
        time.sleep(10)
    pytest.fail(f"ResolverQueryLogConfig {ref.name} did not reach CREATED within {timeout}s")


@pytest.fixture
def resolver_query_log_config(route53resolver_client):
    bucket_name = get_bootstrap_resources().QueryLogBucket.name
    destination_arn = f"arn:aws:s3:::{bucket_name}"

    config_name = random_suffix_name("qlc-test", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["RESOLVER_QUERY_LOG_CONFIG_NAME"] = config_name
    replacements["DESTINATION_ARN"] = destination_arn

    resource_data = load_route53resolver_resource(
        "resolver_query_log_config",
        additional_replacements=replacements,
    )

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
        config_name, namespace="default",
    )

    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)
    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr, bucket_name)

    try:
        if k8s.get_resource_exists(ref):
            k8s.delete_custom_resource(ref, 3, 10)
    except Exception as e:
        logging.warning(f"Cleanup failed for {config_name}: {e}")


@service_marker
class TestResolverQueryLogConfig:
    @pytest.mark.canary
    def test_create_delete(self, route53resolver_client, resolver_query_log_config):
        (ref, cr, bucket_name) = resolver_query_log_config

        cr = wait_for_created(ref)

        config_id = cr["status"]["id"]
        assert config_id is not None

        condition.assert_synced(ref)

        aws_res = route53resolver_client.get_resolver_query_log_config(
            ResolverQueryLogConfigId=config_id
        )
        config = aws_res["ResolverQueryLogConfig"]
        assert config["Status"] == "CREATED"
        assert config["DestinationArn"] == f"arn:aws:s3:::{bucket_name}"

        _, deleted = k8s.delete_custom_resource(ref, 12, 10)
        assert deleted

        deleted_in_aws = False
        for _ in range(9):
            try:
                route53resolver_client.get_resolver_query_log_config(
                    ResolverQueryLogConfigId=config_id
                )
                time.sleep(10)
            except route53resolver_client.exceptions.ResourceNotFoundException:
                deleted_in_aws = True
                break

        assert deleted_in_aws

    @pytest.mark.canary
    def test_tags(self, route53resolver_client, resolver_query_log_config):
        (ref, cr, _) = resolver_query_log_config

        cr = wait_for_created(ref)
        config_id = cr["status"]["id"]
        arn = cr["status"]["ackResourceMetadata"]["arn"]

        tags_res = route53resolver_client.list_tags_for_resource(ResourceArn=arn)
        tag_map = {t["Key"]: t["Value"] for t in tags_res["Tags"]}
        assert tag_map.get("managed-by") == "ack-e2e-test"

        updates = {
            "spec": {
                "tags": [
                    {"key": "managed-by", "value": "ack-e2e-test"},
                    {"key": "env", "value": "testing"},
                ]
            }
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(15)

        tags_res = route53resolver_client.list_tags_for_resource(ResourceArn=arn)
        tag_map = {t["Key"]: t["Value"] for t in tags_res["Tags"]}
        assert tag_map.get("env") == "testing"
        assert tag_map.get("managed-by") == "ack-e2e-test"
