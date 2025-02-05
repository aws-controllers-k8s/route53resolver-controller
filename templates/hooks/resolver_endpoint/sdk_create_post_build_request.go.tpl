    // A unique string that identifies the request and that allows failed requests to be
    // retried without the risk of running the operation twice.
    // CreatorRequestId can be any unique string, for example, a date/time stamp.
    // TODO: Name is not sufficient, since a failed request cannot be retried.
    // We might need to import the `time` package into `sdk.go`
	input.CreatorRequestId = getCreatorRequestId(desired.ko)