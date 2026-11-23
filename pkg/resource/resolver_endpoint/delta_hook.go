package resolver_endpoint

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

// customPostCompare reports a Spec.IPAddresses difference exactly when
// SyncIPAddresses has addresses to associate or disassociate, by asking the
// same differ the sync path uses. The generated order-sensitive DeepEqual
// reported a difference no update could resolve, because AWS returns the
// addresses in an order of its choosing, fills in an IP the user may have left
// unset, and never echoes SubnetRef (community#3029).
//
// Sharing one differ is what keeps the two sides from drifting apart: each used
// to identify an address by its subnet alone, so a changed ip was invisible
// here and in the sync path both, and the endpoint silently kept the old
// address (community#3090).
func customPostCompare(
	delta *ackcompare.Delta,
	a *resource,
	b *resource,
) {
	if a == nil || b == nil || a.ko == nil || b.ko == nil {
		return
	}
	added, _, unmatchedObserved := ipAddressDifference(a, b)
	if len(added) > 0 || unmatchedObserved > 0 {
		delta.Add("Spec.IPAddresses", a.ko.Spec.IPAddresses, b.ko.Spec.IPAddresses)
	}
}
