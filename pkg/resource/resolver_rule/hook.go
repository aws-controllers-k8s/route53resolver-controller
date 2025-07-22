package resolver_rule

import (
	"context"
	"fmt"
	"math"
	"time"

	svcapitypes "github.com/aws-controllers-k8s/route53resolver-controller/apis/v1alpha1"
	"github.com/aws-controllers-k8s/route53resolver-controller/pkg/tags"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/route53resolver"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/route53resolver/types"
	"github.com/samber/lo"
)

var (
	TypeVPCId       = "VPCID"
	RequeueOnUpdate = requeue.Needed(fmt.Errorf("requeing to sync resource status"))
)

// getCreatorRequestId will generate a CreatorRequestId for a given resolver endpoint
// using the name of the endpoint and the current timestamp, so that it produces a
// unique value
func getCreatorRequestId(rule *svcapitypes.ResolverRule) *string {
	requestId := fmt.Sprintf("%s-%d", *rule.Spec.Name, time.Now().UnixMilli())
	return &requestId
}

// addRulesToSpec updates a resource's Spec EgressRules and IngressRules
// using data from a DescribeSecurityGroups response
func (rm *resourceManager) getAttachedVPC(
	ctx context.Context,
	latest *resource,
) (associationList []*svcapitypes.ResolverRuleAssociation, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.getAttachedVPC")
	defer func(err error) {
		exit(err)
	}(err)

	// ko.Spec.Associations = nil
	input := &svcsdk.ListResolverRuleAssociationsInput{
		Filters: []svcsdktypes.Filter{
			{
				Name:   lo.ToPtr("ResolverRuleId"),
				Values: []string{*latest.ko.Status.ID},
			},
		},
	}
	resolverRuleList, err := rm.sdkapi.ListResolverRuleAssociations(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "ListResolverRuleAssociations", err)
	for _, association := range resolverRuleList.ResolverRuleAssociations {
		if association.Status != svcsdktypes.ResolverRuleAssociationStatusDeleting {
			var svcassociation svcapitypes.ResolverRuleAssociation
			svcassociation.VPCID = association.VPCId
			associationList = append(associationList, &svcassociation)
		}
	}
	return associationList, nil
}

func (rm *resourceManager) customUpdateResolverRule(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.customUpdateResolverRule")
	defer exit(err)

	if delta.DifferentAt("Spec.Tags") {
		if err = rm.syncTags(ctx, desired, latest); err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.Associations") {
		if err := rm.syncAssociation(ctx, desired, latest); err != nil {
			return nil, err
		}
	}

	if !delta.DifferentExcept("Spec.Tags", "Spec.Associations") {
		return desired, nil
	}

	// Default `updated` to `desired` because it is likely
	// EC2 `modify` APIs do NOT return output, only errors.
	// If the `modify` calls (i.e. `sync`) do NOT return
	// an error, then the update was successful and desired.Spec
	// (now updated.Spec) reflects the latest resource state.
	updated = rm.concreteResource(desired.DeepCopy())

	if delta.DifferentAt("Spec.TargetIPs") || delta.DifferentAt("Spec.Name") {
		if err := rm.syncResolverRuleConfig(ctx, desired, latest); err != nil {
			return nil, err
		}
	}

	return updated, RequeueOnUpdate
}

func (rm *resourceManager) createAssociation(
	ctx context.Context,
	r *resource,
) error {
	if r.ko.Spec.Associations != nil {
		if err := rm.syncAssociation(ctx, r, nil); err != nil {
			return err
		}
	}
	return nil
}

func (rm *resourceManager) syncResolverRuleConfig(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncResolverRuleConfig")
	defer exit(err)
	input := &svcsdk.UpdateResolverRuleInput{}
	if latest.ko.Status.ID != nil {
		input.ResolverRuleId = latest.ko.Status.ID
	}
	resconf := &svcsdktypes.ResolverRuleConfig{}
	resconf.Name = desired.ko.Spec.Name
	resconf.ResolverEndpointId = desired.ko.Spec.ResolverEndpointID
	var targip []svcsdktypes.TargetAddress
	for _, tip := range desired.ko.Spec.TargetIPs {
		targipelem := svcsdktypes.TargetAddress{}
		if tip.IP != nil {
			targipelem.Ip = tip.IP
		}
		if tip.IPv6 != nil {
			targipelem.Ipv6 = tip.IPv6
		}
		if tip.Port != nil {
			if *tip.Port > math.MaxInt32 || *tip.Port < math.MinInt32 {
				return fmt.Errorf("error: field TargetAddress.Port is of type int32")
			}
			portCopy := int32(*tip.Port)
			targipelem.Port = &portCopy
		}
		targip = append(targip, targipelem)
	}
	resconf.TargetIps = targip
	resconf.Name = desired.ko.Spec.Name
	input.Config = resconf
	var resp *svcsdk.UpdateResolverRuleOutput
	_ = resp
	_, err = rm.sdkapi.UpdateResolverRule(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateResolverRule", err)
	if err != nil {
		return err
	}
	return nil
}

func (rm *resourceManager) syncAssociation(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncAssociation")
	defer exit(err)

	latestAssociations := make(map[string]string)
	desiredAssociations := make(map[string]string)
	associationidVpc := make(map[string]string)

	if latest != nil {
		for _, association := range latest.ko.Spec.Associations {
			if association.VPCID != nil {
				latestAssociations[*association.VPCID] = TypeVPCId
			}
		}
	}
	if desired != nil {
		for _, association := range desired.ko.Spec.Associations {
			if association.VPCID != nil {
				desiredAssociations[*association.VPCID] = TypeVPCId
			}
		}
	}
	// Determining the associations to be added and deleted by comparing associations of latest and desired.
	toAdd := lo.OmitByKeys(desiredAssociations, lo.Keys(latestAssociations))
	includedVpcs := lo.PickByKeys(associationidVpc, lo.Keys(desiredAssociations))
	associations_diff := lo.OmitByKeys(latestAssociations, lo.Keys(desiredAssociations))
	toDelete := lo.OmitByKeys(associations_diff, lo.Values(includedVpcs))

	upsertErr := rm.upsertNewAssociations(ctx, desired, latest, toAdd)
	if upsertErr != nil {
		return upsertErr
	}
	deletErr := rm.deleteOldAssociations(ctx, desired, latest, toDelete)
	if deletErr != nil {
		return deletErr
	}
	return nil

}

func (rm *resourceManager) deleteOldAssociations(
	ctx context.Context,
	desired *resource,
	latest *resource,
	toDelete map[string]string,
) (err error) {
	for rid, rtype := range toDelete {
		input := &svcsdk.DisassociateResolverRuleInput{}
		if rtype == TypeVPCId {
			input.ResolverRuleId = desired.ko.Status.ID
			input.VPCId = &rid
			_, err = rm.sdkapi.DisassociateResolverRule(ctx, input)
			rm.metrics.RecordAPICall("UPDATE", "DisassociateResolverRule", err)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (rm *resourceManager) upsertNewAssociations(
	ctx context.Context,
	desired *resource,
	latest *resource,
	toAdd map[string]string,
) (err error) {
	for rid, rtype := range toAdd {
		input := &svcsdk.AssociateResolverRuleInput{}
		if rtype == TypeVPCId {
			input.ResolverRuleId = desired.ko.Status.ID
			input.VPCId = &rid
			_, err = rm.sdkapi.AssociateResolverRule(ctx, input)
			rm.metrics.RecordAPICall("UPDATE", "AssociateResolverRule", err)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// getTags retrieves the resource's associated tags.
func (rm *resourceManager) getTags(
	ctx context.Context,
	resourceARN string,
) ([]*svcapitypes.Tag, error) {
	return tags.GetTags(ctx, rm.sdkapi, rm.metrics, resourceARN)
}

// syncTags keeps the resource's tags in sync.
func (rm *resourceManager) syncTags(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (err error) {
	return tags.SyncTags(ctx, desired.ko.Spec.Tags, latest.ko.Spec.Tags, latest.ko.Status.ACKResourceMetadata, convertToOrderedACKTags, rm.sdkapi, rm.metrics)
}
