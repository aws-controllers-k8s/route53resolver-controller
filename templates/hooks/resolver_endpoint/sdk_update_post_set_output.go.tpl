	if delta.DifferentAt("Spec.IPAddresses") {
		err = rm.SyncIPAddresses(ctx, desired, latest)
		ko.Status.IPAddressCount = latest.ko.Status.IPAddressCount
		if err != nil {
			// Return the resource so its status is still persisted, and let the
			// error reach the reconciler: a requeue completes deferred removals,
			// and a real failure is no longer silently dropped.
			return &resource{ko}, err
		}
	}
