	ko.Spec.Associations, err = rm.getAttachedVPC(ctx,&resource{ko})
	if err != nil {
		return nil, err
	}
	