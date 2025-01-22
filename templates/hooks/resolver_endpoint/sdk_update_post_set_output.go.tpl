	if delta.DifferentAt("Spec.IPAddresses") {
		rm.SyncIPAddresses(ctx, desired, latest)
		ko.Status.IPAddressCount = latest.ko.Status.IPAddressCount
	}