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

"""Cleans up the resources created by the bootstrapping process.

The generic teardown in `Resources.cleanup()` deletes the bootstrapped VPCs and
bucket in reverse creation order, which is correct but not sufficient here: a
Resolver endpoint owns service-managed ENIs in the VPC's subnets, so any endpoint
that outlives its test blocks subnet deletion, which blocks VPC deletion
(DependencyViolation). The generic cleanup logs that failure and continues, so
the VPC leaks silently and the account eventually hits VpcLimitExceeded.

This module therefore sweeps the Resolver resources in the bootstrapped VPCs
first, in dependency order, and only then runs the generic teardown.
"""

import logging
import time

import boto3
from botocore.exceptions import ClientError

from acktest.bootstrapping import Resources

from e2e import bootstrap_directory
from e2e.bootstrap_resources import BootstrapResources

# An endpoint's ENIs are released as it is deleted, so the subnets it occupies
# stay undeletable until that finishes.
ENDPOINT_DELETE_TIMEOUT_SECONDS = 600
ENDPOINT_DELETE_INTERVAL_SECONDS = 15


def route53resolver_client():
    return boto3.client("route53resolver")


def _list_all(list_fn, result_key: str, **kwargs):
    """Yield every item from a paginated Resolver list call.

    Paged explicitly rather than through a botocore paginator, so this does not
    depend on a paginator being configured for each operation.
    """
    next_token = None
    while True:
        if next_token:
            kwargs["NextToken"] = next_token
        response = list_fn(**kwargs)
        for item in response.get(result_key, []):
            yield item
        next_token = response.get("NextToken")
        if not next_token:
            return


def list_endpoints_in_vpcs(vpc_ids: list) -> list:
    """Return the IDs of every Resolver endpoint hosted in the given VPCs."""
    r53r = route53resolver_client()
    endpoint_ids = []
    for vpc_id in vpc_ids:
        endpoints = _list_all(
            r53r.list_resolver_endpoints,
            "ResolverEndpoints",
            Filters=[{"Name": "HostVPCId", "Values": [vpc_id]}],
        )
        for endpoint in endpoints:
            endpoint_ids.append(endpoint["Id"])
    return endpoint_ids


def rules_for_endpoints(endpoint_ids: list) -> list:
    """Return the IDs of every Resolver rule that forwards to one of the given
    endpoints.

    Rules with no endpoint are skipped, which excludes the RECURSIVE rule AWS
    auto-defines per VPC. That rule cannot be disassociated
    ([RSLVR-00713]) and is not ours to remove.
    """
    if not endpoint_ids:
        return []

    r53r = route53resolver_client()
    rule_ids = []
    for rule in _list_all(r53r.list_resolver_rules, "ResolverRules"):
        if rule.get("ResolverEndpointId") in endpoint_ids:
            rule_ids.append(rule["Id"])
    return rule_ids


def disassociate_rules(rule_ids: list):
    """Disassociate the given rules from every VPC they are associated with.

    A rule cannot be deleted while it is associated with a VPC.
    """
    r53r = route53resolver_client()
    for rule_id in rule_ids:
        associations = _list_all(
            r53r.list_resolver_rule_associations,
            "ResolverRuleAssociations",
            Filters=[{"Name": "ResolverRuleId", "Values": [rule_id]}],
        )
        for association in associations:
            r53r.disassociate_resolver_rule(
                VPCId=association["VPCId"],
                ResolverRuleId=association["ResolverRuleId"],
            )
            logging.info(
                f"Disassociated Resolver rule {association['ResolverRuleId']} "
                f"from VPC {association['VPCId']}"
            )


def delete_rules(rule_ids: list):
    """Delete the given Resolver rules.

    An endpoint cannot be deleted while a rule still targets it.
    """
    r53r = route53resolver_client()
    for rule_id in rule_ids:
        r53r.delete_resolver_rule(ResolverRuleId=rule_id)
        logging.info(f"Deleted Resolver rule {rule_id}")


def disassociate_query_log_configs_from_vpcs(vpc_ids: list):
    """Disassociate every query log config associated with the given VPCs."""
    r53r = route53resolver_client()
    for vpc_id in vpc_ids:
        associations = _list_all(
            r53r.list_resolver_query_log_config_associations,
            "ResolverQueryLogConfigAssociations",
            Filters=[{"Name": "ResourceId", "Values": [vpc_id]}],
        )
        for association in associations:
            r53r.disassociate_resolver_query_log_config(
                ResolverQueryLogConfigId=association["ResolverQueryLogConfigId"],
                ResourceId=association["ResourceId"],
            )
            logging.info(
                "Disassociated query log config "
                f"{association['ResolverQueryLogConfigId']} from "
                f"{association['ResourceId']}"
            )


def delete_query_log_configs_for_bucket(bucket_name: str):
    """Delete every query log config whose destination is the given bucket.

    The config holds a reference to the bootstrapped bucket, so it must go before
    the bucket can be emptied and removed.
    """
    r53r = route53resolver_client()
    configs = _list_all(
        r53r.list_resolver_query_log_configs, "ResolverQueryLogConfigs"
    )
    for config in configs:
        if bucket_name in config.get("DestinationArn", ""):
            r53r.delete_resolver_query_log_config(
                ResolverQueryLogConfigId=config["Id"],
            )
            logging.info(f"Deleted query log config {config['Id']}")


def delete_endpoints(endpoint_ids: list):
    """Delete the given endpoints and wait for them to disappear.

    Waiting matters: the ENIs are released as part of the delete, and the generic
    VPC teardown that follows cannot succeed until they are gone.
    """
    r53r = route53resolver_client()
    for endpoint_id in endpoint_ids:
        r53r.delete_resolver_endpoint(ResolverEndpointId=endpoint_id)
        logging.info(f"Deleting Resolver endpoint {endpoint_id}")

    deadline = time.time() + ENDPOINT_DELETE_TIMEOUT_SECONDS
    pending = list(endpoint_ids)
    while pending and time.time() < deadline:
        still_pending = []
        for endpoint_id in pending:
            try:
                r53r.get_resolver_endpoint(ResolverEndpointId=endpoint_id)
                still_pending.append(endpoint_id)
            except ClientError as ex:
                if ex.response["Error"]["Code"] == "ResourceNotFoundException":
                    logging.info(f"Deleted Resolver endpoint {endpoint_id}")
                else:
                    raise
        pending = still_pending
        if pending:
            time.sleep(ENDPOINT_DELETE_INTERVAL_SECONDS)

    if pending:
        logging.warning(
            f"Resolver endpoints still present after {ENDPOINT_DELETE_TIMEOUT_SECONDS}s, "
            f"VPC teardown may fail: {pending}"
        )


def service_cleanup():
    logging.getLogger().setLevel(logging.INFO)

    resources = BootstrapResources.deserialize(bootstrap_directory)

    vpc_ids = [
        resources.ResolverEndpointVPC.vpc_id,
        resources.AssociationTestVPC.vpc_id,
    ]

    # Each step is independent: log and carry on, so one failure cannot stop the
    # rest of the sweep or the generic teardown below.
    endpoint_ids = []
    try:
        endpoint_ids = list_endpoints_in_vpcs(vpc_ids)
        logging.info(f"Found {len(endpoint_ids)} Resolver endpoint(s) to clean up")
    except Exception:
        logging.exception(f"Unable to list Resolver endpoints in {vpc_ids}")

    try:
        disassociate_query_log_configs_from_vpcs(vpc_ids)
    except Exception:
        logging.exception(f"Unable to disassociate query log configs from {vpc_ids}")

    try:
        delete_query_log_configs_for_bucket(resources.QueryLogBucket.name)
    except Exception:
        logging.exception(
            f"Unable to delete query log configs for bucket {resources.QueryLogBucket.name}"
        )

    rule_ids = []
    try:
        rule_ids = rules_for_endpoints(endpoint_ids)
    except Exception:
        logging.exception(f"Unable to list Resolver rules for {endpoint_ids}")

    try:
        disassociate_rules(rule_ids)
    except Exception:
        logging.exception(f"Unable to disassociate Resolver rules {rule_ids}")

    try:
        delete_rules(rule_ids)
    except Exception:
        logging.exception(f"Unable to delete Resolver rules {rule_ids}")

    try:
        delete_endpoints(endpoint_ids)
    except Exception:
        logging.exception(f"Unable to delete Resolver endpoints {endpoint_ids}")

    # Now that nothing holds an ENI in the subnets, the generic teardown can
    # remove the bootstrapped VPCs and bucket.
    resources.cleanup()


if __name__ == "__main__":
    service_cleanup()
