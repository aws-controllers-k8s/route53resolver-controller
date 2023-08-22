	if r.ko.Spec.Associations != nil && r.ko.Status.ID != nil {
		desired := rm.concreteResource(r.DeepCopy())
		desired.ko.Spec.Associations = nil
		if err = rm.syncAssociation(ctx, desired, r); err != nil {
			return nil, err
		}
	}