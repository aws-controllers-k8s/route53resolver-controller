// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package resolver_endpoint

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	svcapitypes "github.com/aws-controllers-k8s/route53resolver-controller/apis/v1alpha1"
	ackmetrics "github.com/aws-controllers-k8s/runtime/pkg/metrics"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/route53resolver"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/route53resolver/types"
	"github.com/go-logr/logr"
)

const (
	statusAttached  = string(svcsdktypes.IpAddressStatusAttached)
	statusAttaching = string(svcsdktypes.IpAddressStatusAttaching)
	statusDetaching = string(svcsdktypes.IpAddressStatusDetaching)
	statusCreating  = string(svcsdktypes.IpAddressStatusCreating)
	statusFailed    = string(svcsdktypes.IpAddressStatusFailedCreation)
)

// specIPs builds a Spec.IPAddresses list from subnet IDs.
func specIPs(subnetIDs ...string) []*svcapitypes.IPAddressRequest {
	out := make([]*svcapitypes.IPAddressRequest, 0, len(subnetIDs))
	for _, s := range subnetIDs {
		out = append(out, &svcapitypes.IPAddressRequest{SubnetID: aws.String(s)})
	}
	return out
}

// observedIP is one observed address: subnet, IP-address ID, and status.
type observedIP struct {
	subnetID string
	ipID     string
	status   string
}

// observed builds a resource whose Spec and Status IP lists are populated the
// way ListAttachedIPAddresses populates them (paired, same order).
func observed(ips ...observedIP) *resource {
	ko := &svcapitypes.ResolverEndpoint{}
	ko.Status.ID = aws.String("rslvr-out-test")
	for _, ip := range ips {
		ko.Spec.IPAddresses = append(ko.Spec.IPAddresses,
			&svcapitypes.IPAddressRequest{SubnetID: aws.String(ip.subnetID)})
		entry := &svcapitypes.IPAddressResponse{
			SubnetID: aws.String(ip.subnetID),
			Status:   aws.String(ip.status),
		}
		if ip.ipID != "" {
			entry.IPID = aws.String(ip.ipID)
		}
		ko.Status.IPAddresses = append(ko.Status.IPAddresses, entry)
	}
	return &resource{ko}
}

// desiredWith builds a desired resource declaring the given subnets.
func desiredWith(subnetIDs ...string) *resource {
	ko := &svcapitypes.ResolverEndpoint{}
	ko.Spec.IPAddresses = specIPs(subnetIDs...)
	return &resource{ko}
}

func TestGetIPAddressDifference(t *testing.T) {
	tests := []struct {
		name        string
		desired     *resource
		latest      *resource
		wantAdded   []string
		wantRemoved []string
	}{
		{
			name:        "full subnet swap",
			desired:     desiredWith("subnet-b1", "subnet-b2", "subnet-b3"),
			latest:      observed(observedIP{"subnet-a1", "ip-a1", statusAttached}, observedIP{"subnet-a2", "ip-a2", statusAttached}, observedIP{"subnet-a3", "ip-a3", statusAttached}),
			wantAdded:   []string{"subnet-b1", "subnet-b2", "subnet-b3"},
			wantRemoved: []string{"ip-a1", "ip-a2", "ip-a3"},
		},
		{
			name:        "partial overlap keeps shared subnet",
			desired:     desiredWith("subnet-a1", "subnet-b1"),
			latest:      observed(observedIP{"subnet-a1", "ip-a1", statusAttached}, observedIP{"subnet-a2", "ip-a2", statusAttached}),
			wantAdded:   []string{"subnet-b1"},
			wantRemoved: []string{"ip-a2"},
		},
		{
			name:        "no change",
			desired:     desiredWith("subnet-a1", "subnet-a2"),
			latest:      observed(observedIP{"subnet-a1", "ip-a1", statusAttached}, observedIP{"subnet-a2", "ip-a2", statusAttached}),
			wantAdded:   nil,
			wantRemoved: nil,
		},
		{
			name:        "stale extra address is selected for removal",
			desired:     desiredWith("subnet-b1", "subnet-b2", "subnet-b3"),
			latest:      observed(observedIP{"subnet-a3", "ip-a3", statusAttached}, observedIP{"subnet-b1", "ip-b1", statusAttached}, observedIP{"subnet-b2", "ip-b2", statusAttached}, observedIP{"subnet-b3", "ip-b3", statusAttached}),
			wantAdded:   nil,
			wantRemoved: []string{"ip-a3"},
		},
		{
			name:        "desired entry with nil subnet is skipped",
			desired:     &resource{&svcapitypes.ResolverEndpoint{Spec: svcapitypes.ResolverEndpointSpec{IPAddresses: []*svcapitypes.IPAddressRequest{nil, {SubnetID: nil}, {SubnetID: aws.String("subnet-b1")}}}}},
			latest:      observed(observedIP{"subnet-b1", "ip-b1", statusAttached}),
			wantAdded:   nil,
			wantRemoved: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm := &resourceManager{}
			added, removed := rm.GetIPAddressDifference(tt.desired, tt.latest)

			gotAdded := make([]string, 0, len(added))
			for _, a := range added {
				gotAdded = append(gotAdded, *a.SubnetID)
			}
			gotRemoved := make([]string, 0, len(removed))
			for _, r := range removed {
				gotRemoved = append(gotRemoved, *r.IPID)
			}
			assertStrings(t, "added", tt.wantAdded, gotAdded)
			assertStrings(t, "removed", tt.wantRemoved, gotRemoved)
		})
	}
}

// TestGetIPAddressDifferenceUnpairedSlices guards the positional Status lookup:
// nothing in the schema enforces that Spec and Status IP lists are equal length.
func TestGetIPAddressDifferenceUnpairedSlices(t *testing.T) {
	latest := observed(observedIP{"subnet-a1", "ip-a1", statusAttached}, observedIP{"subnet-a2", "ip-a2", statusAttached})
	// Truncate Status so it is shorter than Spec.
	latest.ko.Status.IPAddresses = latest.ko.Status.IPAddresses[:1]

	rm := &resourceManager{}
	_, removed := rm.GetIPAddressDifference(desiredWith("subnet-b1"), latest)

	if len(removed) != 1 || *removed[0].IPID != "ip-a1" {
		t.Fatalf("expected only the paired entry, got %d entries", len(removed))
	}
}

func TestIsIPAddressAttachedAndCount(t *testing.T) {
	if isIPAddressAttached(nil) {
		t.Error("nil entry must not count as attached")
	}
	if isIPAddressAttached(&svcapitypes.IPAddressResponse{Status: nil}) {
		t.Error("nil status must not count as attached")
	}
	if !isIPAddressAttached(&svcapitypes.IPAddressResponse{Status: aws.String(statusAttached)}) {
		t.Error("ATTACHED must count as attached")
	}
	for _, s := range []string{statusAttaching, statusDetaching, statusCreating, statusFailed} {
		if isIPAddressAttached(&svcapitypes.IPAddressResponse{Status: aws.String(s)}) {
			t.Errorf("%s must not count as attached", s)
		}
	}

	r := observed(
		observedIP{"subnet-a1", "ip-a1", statusAttached},
		observedIP{"subnet-a2", "ip-a2", statusAttached},
		observedIP{"subnet-b1", "ip-b1", statusAttaching},
	)
	if got := countAttachedIPAddresses(r); got != 2 {
		t.Errorf("expected 2 attached, got %d", got)
	}
}

// recordingTransport serves canned Route 53 Resolver responses and records
// which operation each request targeted (via the X-Amz-Target header).
type recordingTransport struct {
	calls []string
	err   error
}

func (rt *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	target := req.Header.Get("X-Amz-Target")
	if i := strings.LastIndex(target, "."); i >= 0 {
		target = target[i+1:]
	}
	rt.calls = append(rt.calls, target)
	if rt.err != nil {
		return nil, rt.err
	}
	body := `{"ResolverEndpoint":{"IpAddressCount":3}}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/x-amz-json-1.1"}},
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Request:    req,
	}, nil
}

func (rt *recordingTransport) count(op string) int {
	n := 0
	for _, c := range rt.calls {
		if c == op {
			n++
		}
	}
	return n
}

func newTestResourceManager(rt *recordingTransport) *resourceManager {
	cfg := aws.Config{
		Region:      "us-west-2",
		Credentials: staticCreds{},
		HTTPClient:  &http.Client{Transport: rt},
		// One attempt per call, so recorded call counts reflect what the hook
		// asked for rather than SDK retry behaviour.
		RetryMaxAttempts: 1,
	}
	return &resourceManager{
		log:     logr.Discard(),
		metrics: ackmetrics.NewMetrics("route53resolver"),
		sdkapi:  svcsdk.NewFromConfig(cfg),
	}
}

type staticCreds struct{}

func (staticCreds) Retrieve(context.Context) (aws.Credentials, error) {
	return aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET"}, nil
}

func TestSyncIPAddresses(t *testing.T) {
	tests := []struct {
		name            string
		desired         *resource
		latest          *resource
		transportErr    error
		wantAssociate   int
		wantDisassociat int
		wantRequeue     bool
		wantTerminal    bool
		wantRawErr      bool
	}{
		{
			// The reported bug: a full swap must NOT disassociate in the same
			// pass, because the new addresses are not attached yet.
			name:            "full swap associates then requeues without removing",
			desired:         desiredWith("subnet-b1", "subnet-b2", "subnet-b3"),
			latest:          observed(observedIP{"subnet-a1", "ip-a1", statusAttached}, observedIP{"subnet-a2", "ip-a2", statusAttached}, observedIP{"subnet-a3", "ip-a3", statusAttached}),
			wantAssociate:   3,
			wantDisassociat: 0,
			wantRequeue:     true,
		},
		{
			// The requeued pass: new addresses are attached, so the stale one
			// can now be removed while three attached addresses remain.
			name:            "stale address removed once headroom exists",
			desired:         desiredWith("subnet-b1", "subnet-b2", "subnet-b3"),
			latest:          observed(observedIP{"subnet-a3", "ip-a3", statusAttached}, observedIP{"subnet-b1", "ip-b1", statusAttached}, observedIP{"subnet-b2", "ip-b2", statusAttached}, observedIP{"subnet-b3", "ip-b3", statusAttached}),
			wantAssociate:   0,
			wantDisassociat: 1,
		},
		{
			name:            "removal deferred when it would breach the two-address floor",
			desired:         desiredWith("subnet-a2", "subnet-b1"),
			latest:          observed(observedIP{"subnet-a1", "ip-a1", statusAttached}, observedIP{"subnet-a2", "ip-a2", statusAttached}, observedIP{"subnet-b1", "ip-b1", statusAttaching}),
			wantAssociate:   0,
			wantDisassociat: 0,
			wantRequeue:     true,
		},
		{
			name:            "non-attached candidate is deferred, never double-removed",
			desired:         desiredWith("subnet-a1", "subnet-a2"),
			latest:          observed(observedIP{"subnet-a1", "ip-a1", statusAttached}, observedIP{"subnet-a2", "ip-a2", statusAttached}, observedIP{"subnet-a3", "ip-a3", statusDetaching}),
			wantAssociate:   0,
			wantDisassociat: 0,
			wantRequeue:     true,
		},
		{
			name:            "candidate without an IP address ID is deferred, not reported synced",
			desired:         desiredWith("subnet-a1", "subnet-a2"),
			latest:          observed(observedIP{"subnet-a1", "ip-a1", statusAttached}, observedIP{"subnet-a2", "ip-a2", statusAttached}, observedIP{"subnet-a3", "", statusAttached}),
			wantAssociate:   0,
			wantDisassociat: 0,
			wantRequeue:     true,
		},
		{
			name:            "desired set below the service minimum is terminal",
			desired:         desiredWith("subnet-a1"),
			latest:          observed(observedIP{"subnet-a1", "ip-a1", statusAttached}, observedIP{"subnet-a2", "ip-a2", statusAttached}),
			wantAssociate:   0,
			wantDisassociat: 0,
			wantTerminal:    true,
		},
		{
			// knottnt on PR #80: a user who removes an address that failed to
			// attach must be able to converge. A never-attached address holds
			// no slot against the floor, so attempt it.
			name:            "failed address is removed so the user converges",
			desired:         desiredWith("subnet-a1", "subnet-a2"),
			latest:          observed(observedIP{"subnet-a1", "ip-a1", statusAttached}, observedIP{"subnet-a2", "ip-a2", statusAttached}, observedIP{"subnet-b1", "ip-b1", statusFailed}),
			wantAssociate:   0,
			wantDisassociat: 1,
		},
		{
			// Attempting it is safe because AWS refusing is not fatal: defer
			// and retry rather than failing the reconcile.
			name:            "rejected removal of a failed address defers, does not error",
			desired:         desiredWith("subnet-a1", "subnet-a2"),
			latest:          observed(observedIP{"subnet-a1", "ip-a1", statusAttached}, observedIP{"subnet-a2", "ip-a2", statusAttached}, observedIP{"subnet-b1", "ip-b1", statusFailed}),
			transportErr:    errors.New("boom"),
			wantAssociate:   0,
			wantDisassociat: 1,
			wantRequeue:     true,
		},
		{
			name:            "no drift makes no calls",
			desired:         desiredWith("subnet-a1", "subnet-a2"),
			latest:          observed(observedIP{"subnet-a1", "ip-a1", statusAttached}, observedIP{"subnet-a2", "ip-a2", statusAttached}),
			wantAssociate:   0,
			wantDisassociat: 0,
		},
		{
			// A real API failure must surface raw so controller-runtime backs
			// off, rather than becoming a fixed-interval requeue.
			name:          "associate failure propagates as a raw error",
			desired:       desiredWith("subnet-b1", "subnet-b2"),
			latest:        observed(observedIP{"subnet-a1", "ip-a1", statusAttached}, observedIP{"subnet-a2", "ip-a2", statusAttached}),
			transportErr:  errors.New("boom"),
			wantAssociate: 1,
			wantRawErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &recordingTransport{err: tt.transportErr}
			rm := newTestResourceManager(rt)

			err := rm.SyncIPAddresses(context.Background(), tt.desired, tt.latest)

			if got := rt.count("AssociateResolverEndpointIpAddress"); got != tt.wantAssociate {
				t.Errorf("associate calls: want %d, got %d", tt.wantAssociate, got)
			}
			if got := rt.count("DisassociateResolverEndpointIpAddress"); got != tt.wantDisassociat {
				t.Errorf("disassociate calls: want %d, got %d", tt.wantDisassociat, got)
			}

			var requeue *ackrequeue.RequeueNeededAfter
			isRequeue := errors.As(err, &requeue)
			switch {
			case tt.wantRequeue:
				if !isRequeue {
					t.Fatalf("expected a requeue error, got %v", err)
				}
				if requeue.Duration() != resolverEndpointIPRequeueDelay {
					t.Errorf("requeue delay: want %s, got %s", resolverEndpointIPRequeueDelay, requeue.Duration())
				}
			case tt.wantTerminal:
				if err == nil || !strings.Contains(err.Error(), "at least 2") {
					t.Fatalf("expected a terminal minimum-address error, got %v", err)
				}
				if isRequeue {
					t.Error("a non-convergent desired set must not be a requeue")
				}
			case tt.wantRawErr:
				if err == nil {
					t.Fatal("expected an error")
				}
				if isRequeue {
					t.Error("a real API failure must not be wrapped as a requeue")
				}
			default:
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func assertStrings(t *testing.T, label string, want, got []string) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("%s: want %v, got %v", label, want, got)
	}
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("%s: want %v, got %v", label, want, got)
		}
	}
}
