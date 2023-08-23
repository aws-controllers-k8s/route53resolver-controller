	if len(desired.ko.Spec.Associations) > 0 {
		ko.Spec.Associations = desired.ko.Spec.Associations
		if err := rm.createAssociation(ctx, &resource{ko}); err != nil {
			rlog.Debug("Error while syncing Association", err)
		}
	}