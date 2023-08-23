package resolver_rule

import (
	"context"
	"fmt"
	"time"

	svcapitypes "github.com/aws-controllers-k8s/route53resolver-controller/apis/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/samber/lo"
)

var TypeVPCId = "VPCID"

// getCreatorRequestId will generate a CreatorRequestId for a given resolver endpoint
// using the name of the endpoint and the current timestamp, so that it produces a
// unique value
func getCreatorRequestId(rule *svcapitypes.ResolverRule) string {
	return fmt.Sprintf("%s-%d", *rule.Spec.Name, time.Now().UnixMilli())
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
		Filters: []*svcsdk.Filter{
			{
				Name:   lo.ToPtr("ResolverRuleId"),
				Values: []*string{latest.ko.Status.ID},
			},
		},
	}
	resolverRuleList, err := rm.sdkapi.ListResolverRuleAssociationsWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "ListResolverRuleAssociations", err)
	for _, association := range resolverRuleList.ResolverRuleAssociations {
		if *association.Status != "DELETING" {
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

	// Default `updated` to `desired` because it is likely
	// EC2 `modify` APIs do NOT return output, only errors.
	// If the `modify` calls (i.e. `sync`) do NOT return
	// an error, then the update was successful and desired.Spec
	// (now updated.Spec) reflects the latest resource state.
	updated = rm.concreteResource(desired.DeepCopy())

	if delta.DifferentAt("Spec.Associations") {

		if err := rm.syncAssociation(ctx, desired, latest); err != nil {
			return nil, err
		}
		latest.ko.Spec.Associations = desired.ko.Spec.Associations
	}

	if delta.DifferentAt("Spec.TargetIPs") {
		if err := rm.syncResolverRuleConfig(ctx, desired, latest); err != nil {
			return nil, err
		}
	}
	updated, err = rm.sdkFind(ctx, desired)
	if err != nil {
		return nil, err
	}

	return updated, nil
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
		input.SetResolverRuleId(*latest.ko.Status.ID)
	}
	resconf := &svcsdk.ResolverRuleConfig{}
	resconf.SetName(*desired.ko.Spec.Name)
	resconf.SetResolverEndpointId(*desired.ko.Spec.ResolverEndpointID)
	var targip []*svcsdk.TargetAddress
	for _, tip := range desired.ko.Spec.TargetIPs {
		targipelem := &svcsdk.TargetAddress{}
		if tip.IP != nil {
			targipelem.Ip = tip.IP
		}
		if tip.IPv6 != nil {
			targipelem.Ipv6 = tip.IPv6
		}
		if tip.Port != nil {
			targipelem.Port = tip.Port
		}
		targip = append(targip, targipelem)
	}
	resconf.TargetIps = targip
	input.SetConfig(resconf)
	var resp *svcsdk.UpdateResolverRuleOutput
	_ = resp
	_, err = rm.sdkapi.UpdateResolverRuleWithContext(ctx, input)
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
			_, err = rm.sdkapi.DisassociateResolverRuleWithContext(ctx, input)
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
			_, err = rm.sdkapi.AssociateResolverRuleWithContext(ctx, input)
			rm.metrics.RecordAPICall("UPDATE", "AssociateResolverRule", err)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
