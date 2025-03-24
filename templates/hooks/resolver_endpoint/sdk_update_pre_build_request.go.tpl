	if delta.DifferentAt("Spec.Tags") {
		if err = rm.syncTags(ctx, desired, latest); err != nil {
			return nil, err
		}
	} else if !delta.DifferentExcept("Spec.Tags") {
		return desired, nil
	}
