	if ko.Status.Status != nil && *ko.Status.Status != "COMPLETE" {
		return &resource{ko}, requeueWhileCreating
	}
