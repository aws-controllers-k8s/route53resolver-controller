package resolver_endpoint

import (
	"context"
	"fmt"
	"time"

	svcapitypes "github.com/aws-controllers-k8s/route53resolver-controller/apis/v1alpha1"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/route53resolver"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/route53resolver/types"
)

// getCreatorRequestId will generate a CreatorRequestId for a given resolver endpoint
// using the name of the endpoint and the current timestamp, so that it produces a
// unique value
func getCreatorRequestId(endpoint *svcapitypes.ResolverEndpoint) *string {
	requestId := fmt.Sprintf("%s-%d", *endpoint.Spec.Name, time.Now().UnixMilli())
	return &requestId
}

func (rm *resourceManager) ListAttachedIPAddresses(
	ctx context.Context,
	resource *svcapitypes.ResolverEndpoint,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.SyncAttachedIPAddresses")
	defer exit(err)

	var nextToken *string

	f0 := []*svcapitypes.IPAddressRequest{}
	f2 := []*svcapitypes.IPAddressResponse{}
	for {
		resp, err := rm.sdkapi.ListResolverEndpointIpAddresses(
			ctx,
			&svcsdk.ListResolverEndpointIpAddressesInput{
				ResolverEndpointId: resource.Status.ID,
				NextToken:          nextToken,
			},
		)
		rm.metrics.RecordAPICall("READ_MANY", "ListResolverEndpointIpAddresses", err)
		if err != nil {
			return err
		}

		for _, elem := range resp.IpAddresses {
			f1 := &svcapitypes.IPAddressRequest{}
			f3 := &svcapitypes.IPAddressResponse{}
			if elem.Ip != nil {
				f1.IP = elem.Ip
			}
			if elem.Ipv6 != nil {
				f1.IPv6 = elem.Ipv6
			}
			if elem.SubnetId != nil {
				f1.SubnetID = elem.SubnetId
			}
			if elem.CreationTime != nil {
				f3.CreationTime = elem.CreationTime
			}
			if elem.ModificationTime != nil {
				f3.ModificationTime = elem.ModificationTime
			}
			if elem.Status != "" {
				f3.Status = aws.String(string(elem.Status))
			}
			if elem.StatusMessage != nil {
				f3.StatusMessage = elem.StatusMessage
			}
			if elem.IpId != nil {
				f3.IPID = elem.IpId
			}
			f0 = append(f0, f1)
			f2 = append(f2, f3)
		}
		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}
	resource.Spec.IPAddresses = f0
	resource.Status.IPAddresses = f2

	return err
}

func (rm *resourceManager) SyncIPAddresses(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.SyncIPAddresses")
	defer exit(err)

	added, removed := rm.GetIPAddressDifference(desired, latest)

	if len(added) > 0 {
		for _, ipa := range added {
			resp, err := rm.sdkapi.AssociateResolverEndpointIpAddress(
				ctx,
				&svcsdk.AssociateResolverEndpointIpAddressInput{
					IpAddress: &svcsdktypes.IpAddressUpdate{
						Ip:       ipa.IP,
						Ipv6:     ipa.IPv6,
						SubnetId: ipa.SubnetID,
					},
					ResolverEndpointId: latest.ko.Status.ID,
				},
			)
			rm.metrics.RecordAPICall("UPDATE", "AssociateResolverEndpointIpAddress", err)
			if err != nil {
				return err
			}
			countCopy := int64(*resp.ResolverEndpoint.IpAddressCount)
			latest.ko.Status.IPAddressCount = &countCopy
		}
	}

	if len(removed) > 0 {
		for _, ipid := range removed {
			resp, err := rm.sdkapi.DisassociateResolverEndpointIpAddress(
				ctx,
				&svcsdk.DisassociateResolverEndpointIpAddressInput{
					IpAddress: &svcsdktypes.IpAddressUpdate{
						IpId: ipid,
					},
					ResolverEndpointId: latest.ko.Status.ID,
				},
			)
			rm.metrics.RecordAPICall("UPDATE", "DisassociateResolverEndpointIpAddress", err)
			if err != nil {
				return err
			}
			countCopy := int64(*resp.ResolverEndpoint.IpAddressCount)
			latest.ko.Status.IPAddressCount = &countCopy
		}
	}

	return err
}

func (rm *resourceManager) GetIPAddressDifference(
	desired, latest *resource,
) (added []*svcapitypes.IPAddressRequest, removed []*string) {

	for _, ipa := range desired.ko.Spec.IPAddresses {
		if !inIpAddress(*ipa.SubnetID, latest.ko.Spec.IPAddresses) {
			added = append(added, ipa)
		}
	}

	for i, ipa := range latest.ko.Spec.IPAddresses {
		if !inIpAddress(*ipa.SubnetID, desired.ko.Spec.IPAddresses) {
			removed = append(removed, latest.ko.Status.IPAddresses[i].IPID)
		}
	}

	return added, removed
}

func inIpAddress(
	subnetId string,
	ipAddresses []*svcapitypes.IPAddressRequest,
) bool {

	for _, ipa := range ipAddresses {
		if *ipa.SubnetID == subnetId {
			return true
		}
	}
	return false
}
