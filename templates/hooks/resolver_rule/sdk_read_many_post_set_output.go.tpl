	ko.Spec.Associations, err = rm.getAttachedVPC(ctx,&resource{ko})
	if err != nil {
		return nil, err
	}
	
	tags, err := rm.getTags(ctx, string(*ko.Status.ACKResourceMetadata.ARN))
	if err != nil {
		return nil, err
	}
	ko.Spec.Tags = tags
