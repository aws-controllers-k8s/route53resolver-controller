package resolver_rule

import (
	"testing"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"

	svcapitypes "github.com/aws-controllers-k8s/route53resolver-controller/apis/v1alpha1"
)

func TestCompareDomainName(t *testing.T) {
	tests := []struct {
		name      string
		a, b      *string
		wantDelta bool
	}{
		{"both nil", nil, nil, false},
		{"equal", aws.String("example.com"), aws.String("example.com"), false},
		{"trailing dot on latest", aws.String("example.com"), aws.String("example.com."), false},
		{"trailing dot on desired", aws.String("example.com."), aws.String("example.com"), false},
		{"different names", aws.String("example.com"), aws.String("example.org."), true},
		{"nil difference", aws.String("example.com"), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &resource{ko: &svcapitypes.ResolverRule{}}
			b := &resource{ko: &svcapitypes.ResolverRule{}}
			a.ko.Spec.DomainName = tt.a
			b.ko.Spec.DomainName = tt.b
			delta := ackcompare.NewDelta()
			compareDomainName(delta, a, b)
			if got := delta.DifferentAt("Spec.DomainName"); got != tt.wantDelta {
				t.Errorf("DifferentAt(Spec.DomainName) = %v, want %v", got, tt.wantDelta)
			}
		})
	}
}
