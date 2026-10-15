package resolver_endpoint

import (
	"context"
	"fmt"
	"time"

	svcapitypes "github.com/aws-controllers-k8s/route53resolver-controller/apis/v1alpha1"
	"github.com/aws-controllers-k8s/route53resolver-controller/pkg/tags"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/route53resolver"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/route53resolver/types"
)

const (
	// minResolverEndpointIPAddresses is the number of IP addresses a Resolver
	// endpoint must always hold. Route 53 documents this on both
	// CreateResolverEndpointInput.IpAddresses ("Even though the minimum is 1,
	// Route 53 requires that you create at least two") and ResolverEndpoint
	// ("An endpoint must always include at least two IP addresses").
	minResolverEndpointIPAddresses = 2

	// resolverEndpointIPRequeueDelay is how long to wait before retrying IP
	// address removals that were deferred while addresses finish attaching.
	resolverEndpointIPRequeueDelay = 30 * time.Second
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
	defer func() {
		exit(err)
	}()

	added, removed := rm.GetIPAddressDifference(desired, latest)

	// A desired set below the service minimum can never converge, so surface it
	// as terminal rather than retrying forever.
	if len(removed) > 0 &&
		len(desired.ko.Spec.IPAddresses) < minResolverEndpointIPAddresses {
		return ackerr.NewTerminalError(fmt.Errorf(
			"spec.ipAddresses declares %d address(es), but a Resolver endpoint requires at least %d",
			len(desired.ko.Spec.IPAddresses), minResolverEndpointIPAddresses,
		))
	}

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
		setIPAddressCount(latest, resp.ResolverEndpoint)
	}

	// An address associated above is still CREATING/ATTACHING, so it does not
	// yet count toward the endpoint's attached-address floor. Removing an old
	// address now can leave the endpoint with no attached address, which AWS
	// rejects (RSLVR-00506) and which previously stranded that address
	// indefinitely. Requeue instead and remove once the new addresses attach.
	if len(added) > 0 && len(removed) > 0 {
		return ackrequeue.NeededAfter(
			fmt.Errorf(
				"waiting for %d newly associated IP address(es) to attach before removing %d",
				len(added), len(removed),
			),
			resolverEndpointIPRequeueDelay,
		)
	}

	attached := countAttachedIPAddresses(latest)
	deferred := 0
	for _, ipa := range removed {
		if ipa == nil || ipa.IPID == nil {
			deferred++
			continue
		}
		isAttached := isIPAddressAttached(ipa)

		// Only an ATTACHED address holds a slot against the endpoint's
		// two-address floor, so gate just those on remaining headroom.
		if isAttached && attached-1 < minResolverEndpointIPAddresses {
			deferred++
			continue
		}
		// AWS is already removing a DETACHING/DELETING address; a second call
		// would only add noise.
		if isIPAddressLeaving(ipa) {
			deferred++
			continue
		}

		resp, err := rm.sdkapi.DisassociateResolverEndpointIpAddress(
			ctx,
			&svcsdk.DisassociateResolverEndpointIpAddressInput{
				IpAddress: &svcsdktypes.IpAddressUpdate{
					IpId: ipa.IPID,
				},
				ResolverEndpointId: latest.ko.Status.ID,
			},
		)
		rm.metrics.RecordAPICall("UPDATE", "DisassociateResolverEndpointIpAddress", err)
		if err != nil {
			if isAttached {
				return err
			}
			// A never-attached address may be refused while it is still
			// settling. Defer rather than fail, so a user who removed an
			// address that failed to attach still converges.
			deferred++
			continue
		}
		setIPAddressCount(latest, resp.ResolverEndpoint)
		if isAttached {
			attached--
		}
	}

	if deferred > 0 {
		return ackrequeue.NeededAfter(
			fmt.Errorf(
				"deferred %d IP address removal(s) until endpoint IP address transitions settle",
				deferred,
			),
			resolverEndpointIPRequeueDelay,
		)
	}

	return nil
}

// setIPAddressCount records the endpoint's address count reported by an
// associate/disassociate response.
func setIPAddressCount(
	latest *resource,
	endpoint *svcsdktypes.ResolverEndpoint,
) {
	if endpoint == nil || endpoint.IpAddressCount == nil {
		return
	}
	countCopy := int64(*endpoint.IpAddressCount)
	latest.ko.Status.IPAddressCount = &countCopy
}

// isIPAddressAttached reports whether an observed address is ATTACHED, the only
// state in which it counts toward the endpoint's attached-address floor.
func isIPAddressAttached(ipa *svcapitypes.IPAddressResponse) bool {
	return ipa != nil && ipa.Status != nil &&
		*ipa.Status == string(svcsdktypes.IpAddressStatusAttached)
}

// isIPAddressLeaving reports whether AWS is already removing an address, in
// which case re-issuing the disassociate achieves nothing.
func isIPAddressLeaving(ipa *svcapitypes.IPAddressResponse) bool {
	if ipa == nil || ipa.Status == nil {
		return false
	}
	switch *ipa.Status {
	case string(svcsdktypes.IpAddressStatusDetaching),
		string(svcsdktypes.IpAddressStatusDeleting):
		return true
	}
	return false
}

// countAttachedIPAddresses counts the endpoint's observed ATTACHED addresses.
func countAttachedIPAddresses(r *resource) int {
	count := 0
	for _, ipa := range r.ko.Status.IPAddresses {
		if isIPAddressAttached(ipa) {
			count++
		}
	}
	return count
}

func (rm *resourceManager) GetIPAddressDifference(
	desired, latest *resource,
) (added []*svcapitypes.IPAddressRequest, removed []*svcapitypes.IPAddressResponse) {

	for _, ipa := range desired.ko.Spec.IPAddresses {
		if ipa == nil || ipa.SubnetID == nil {
			continue
		}
		if !inIpAddress(*ipa.SubnetID, latest.ko.Spec.IPAddresses) {
			added = append(added, ipa)
		}
	}

	// ListAttachedIPAddresses appends to Spec.IPAddresses and
	// Status.IPAddresses together, so entry i of each describes the same
	// address. Bound the walk anyway: nothing in the schema enforces that, and
	// an unpaired entry would otherwise index out of range.
	paired := len(latest.ko.Spec.IPAddresses)
	if n := len(latest.ko.Status.IPAddresses); n < paired {
		paired = n
	}
	for i := 0; i < paired; i++ {
		ipa := latest.ko.Spec.IPAddresses[i]
		if ipa == nil || ipa.SubnetID == nil {
			continue
		}
		if !inIpAddress(*ipa.SubnetID, desired.ko.Spec.IPAddresses) {
			removed = append(removed, latest.ko.Status.IPAddresses[i])
		}
	}

	return added, removed
}

func inIpAddress(
	subnetId string,
	ipAddresses []*svcapitypes.IPAddressRequest,
) bool {

	for _, ipa := range ipAddresses {
		if ipa == nil || ipa.SubnetID == nil {
			continue
		}
		if *ipa.SubnetID == subnetId {
			return true
		}
	}
	return false
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
