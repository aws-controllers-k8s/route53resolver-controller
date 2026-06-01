package resolver_rule_association

import (
	"errors"
	"time"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
)

var requeueWhileCreating = ackrequeue.NeededAfter(
	errors.New("association is not COMPLETE yet"),
	5*time.Second,
)
