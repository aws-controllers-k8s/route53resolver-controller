package resolver_endpoint

import (
	"fmt"
	"time"

	svcapitypes "github.com/aws-controllers-k8s/route53resolver-controller/apis/v1alpha1"
)

// getCreatorRequestId will generate a CreatorRequestId for a given resolver endpoint
// using the name of the endpoint and the current timestamp, so that it produces a
// unique value
func getCreatorRequestId(endpoint *svcapitypes.ResolverEndpoint) string {
	return fmt.Sprintf("%s-%d", *endpoint.Spec.Name, time.Now().UnixMilli())
}
