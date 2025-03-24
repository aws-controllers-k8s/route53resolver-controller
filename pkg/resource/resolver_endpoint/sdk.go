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

// Code generated by ack-generate. DO NOT EDIT.

package resolver_endpoint

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/route53resolver"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/route53resolver/types"
	smithy "github.com/aws/smithy-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/route53resolver-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &svcsdk.Client{}
	_ = &svcapitypes.ResolverEndpoint{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
	_ = &ackrequeue.NoRequeue{}
	_ = &aws.Config{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.GetResolverEndpointOutput
	resp, err = rm.sdkapi.GetResolverEndpoint(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "GetResolverEndpoint", err)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "ResourceNotFoundException" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.ResolverEndpoint.Arn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.ResolverEndpoint.Arn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.ResolverEndpoint.CreationTime != nil {
		ko.Status.CreationTime = resp.ResolverEndpoint.CreationTime
	} else {
		ko.Status.CreationTime = nil
	}
	if resp.ResolverEndpoint.CreatorRequestId != nil {
		ko.Status.CreatorRequestID = resp.ResolverEndpoint.CreatorRequestId
	} else {
		ko.Status.CreatorRequestID = nil
	}
	if resp.ResolverEndpoint.Direction != "" {
		ko.Spec.Direction = aws.String(string(resp.ResolverEndpoint.Direction))
	} else {
		ko.Spec.Direction = nil
	}
	if resp.ResolverEndpoint.HostVPCId != nil {
		ko.Status.HostVPCID = resp.ResolverEndpoint.HostVPCId
	} else {
		ko.Status.HostVPCID = nil
	}
	if resp.ResolverEndpoint.Id != nil {
		ko.Status.ID = resp.ResolverEndpoint.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.ResolverEndpoint.IpAddressCount != nil {
		ipAddressCountCopy := int64(*resp.ResolverEndpoint.IpAddressCount)
		ko.Status.IPAddressCount = &ipAddressCountCopy
	} else {
		ko.Status.IPAddressCount = nil
	}
	if resp.ResolverEndpoint.ModificationTime != nil {
		ko.Status.ModificationTime = resp.ResolverEndpoint.ModificationTime
	} else {
		ko.Status.ModificationTime = nil
	}
	if resp.ResolverEndpoint.Name != nil {
		ko.Spec.Name = resp.ResolverEndpoint.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.ResolverEndpoint.ResolverEndpointType != "" {
		ko.Spec.ResolverEndpointType = aws.String(string(resp.ResolverEndpoint.ResolverEndpointType))
	} else {
		ko.Spec.ResolverEndpointType = nil
	}
	if resp.ResolverEndpoint.SecurityGroupIds != nil {
		ko.Spec.SecurityGroupIDs = aws.StringSlice(resp.ResolverEndpoint.SecurityGroupIds)
	} else {
		ko.Spec.SecurityGroupIDs = nil
	}
	if resp.ResolverEndpoint.Status != "" {
		ko.Status.Status = aws.String(string(resp.ResolverEndpoint.Status))
	} else {
		ko.Status.Status = nil
	}
	if resp.ResolverEndpoint.StatusMessage != nil {
		ko.Status.StatusMessage = resp.ResolverEndpoint.StatusMessage
	} else {
		ko.Status.StatusMessage = nil
	}

	rm.setStatusDefaults(ko)
	rm.ListAttachedIPAddresses(ctx, ko)

	tags, err := rm.getTags(ctx, string(*ko.Status.ACKResourceMetadata.ARN))
	if err != nil {
		return nil, err
	}
	ko.Spec.Tags = tags

	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return r.ko.Status.ID == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.GetResolverEndpointInput, error) {
	res := &svcsdk.GetResolverEndpointInput{}

	if r.ko.Status.ID != nil {
		res.ResolverEndpointId = r.ko.Status.ID
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer func() {
		exit(err)
	}()
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}
	// A unique string that identifies the request and that allows failed requests to be
	// retried without the risk of running the operation twice.
	// CreatorRequestId can be any unique string, for example, a date/time stamp.
	// TODO: Name is not sufficient, since a failed request cannot be retried.
	// We might need to import the `time` package into `sdk.go`
	input.CreatorRequestId = getCreatorRequestId(desired.ko)

	var resp *svcsdk.CreateResolverEndpointOutput
	_ = resp
	resp, err = rm.sdkapi.CreateResolverEndpoint(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateResolverEndpoint", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.ResolverEndpoint.Arn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.ResolverEndpoint.Arn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.ResolverEndpoint.CreationTime != nil {
		ko.Status.CreationTime = resp.ResolverEndpoint.CreationTime
	} else {
		ko.Status.CreationTime = nil
	}
	if resp.ResolverEndpoint.CreatorRequestId != nil {
		ko.Status.CreatorRequestID = resp.ResolverEndpoint.CreatorRequestId
	} else {
		ko.Status.CreatorRequestID = nil
	}
	if resp.ResolverEndpoint.Direction != "" {
		ko.Spec.Direction = aws.String(string(resp.ResolverEndpoint.Direction))
	} else {
		ko.Spec.Direction = nil
	}
	if resp.ResolverEndpoint.HostVPCId != nil {
		ko.Status.HostVPCID = resp.ResolverEndpoint.HostVPCId
	} else {
		ko.Status.HostVPCID = nil
	}
	if resp.ResolverEndpoint.Id != nil {
		ko.Status.ID = resp.ResolverEndpoint.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.ResolverEndpoint.IpAddressCount != nil {
		ipAddressCountCopy := int64(*resp.ResolverEndpoint.IpAddressCount)
		ko.Status.IPAddressCount = &ipAddressCountCopy
	} else {
		ko.Status.IPAddressCount = nil
	}
	if resp.ResolverEndpoint.ModificationTime != nil {
		ko.Status.ModificationTime = resp.ResolverEndpoint.ModificationTime
	} else {
		ko.Status.ModificationTime = nil
	}
	if resp.ResolverEndpoint.Name != nil {
		ko.Spec.Name = resp.ResolverEndpoint.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.ResolverEndpoint.ResolverEndpointType != "" {
		ko.Spec.ResolverEndpointType = aws.String(string(resp.ResolverEndpoint.ResolverEndpointType))
	} else {
		ko.Spec.ResolverEndpointType = nil
	}
	if resp.ResolverEndpoint.SecurityGroupIds != nil {
		ko.Spec.SecurityGroupIDs = aws.StringSlice(resp.ResolverEndpoint.SecurityGroupIds)
	} else {
		ko.Spec.SecurityGroupIDs = nil
	}
	if resp.ResolverEndpoint.Status != "" {
		ko.Status.Status = aws.String(string(resp.ResolverEndpoint.Status))
	} else {
		ko.Status.Status = nil
	}
	if resp.ResolverEndpoint.StatusMessage != nil {
		ko.Status.StatusMessage = resp.ResolverEndpoint.StatusMessage
	} else {
		ko.Status.StatusMessage = nil
	}

	rm.setStatusDefaults(ko)
	rm.ListAttachedIPAddresses(ctx, ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateResolverEndpointInput, error) {
	res := &svcsdk.CreateResolverEndpointInput{}

	if r.ko.Spec.Direction != nil {
		res.Direction = svcsdktypes.ResolverEndpointDirection(*r.ko.Spec.Direction)
	}
	if r.ko.Spec.IPAddresses != nil {
		f1 := []svcsdktypes.IpAddressRequest{}
		for _, f1iter := range r.ko.Spec.IPAddresses {
			f1elem := &svcsdktypes.IpAddressRequest{}
			if f1iter.IP != nil {
				f1elem.Ip = f1iter.IP
			}
			if f1iter.IPv6 != nil {
				f1elem.Ipv6 = f1iter.IPv6
			}
			if f1iter.SubnetID != nil {
				f1elem.SubnetId = f1iter.SubnetID
			}
			f1 = append(f1, *f1elem)
		}
		res.IpAddresses = f1
	}
	if r.ko.Spec.Name != nil {
		res.Name = r.ko.Spec.Name
	}
	if r.ko.Spec.ResolverEndpointType != nil {
		res.ResolverEndpointType = svcsdktypes.ResolverEndpointType(*r.ko.Spec.ResolverEndpointType)
	}
	if r.ko.Spec.SecurityGroupIDs != nil {
		res.SecurityGroupIds = aws.ToStringSlice(r.ko.Spec.SecurityGroupIDs)
	}
	if r.ko.Spec.Tags != nil {
		f5 := []svcsdktypes.Tag{}
		for _, f5iter := range r.ko.Spec.Tags {
			f5elem := &svcsdktypes.Tag{}
			if f5iter.Key != nil {
				f5elem.Key = f5iter.Key
			}
			if f5iter.Value != nil {
				f5elem.Value = f5iter.Value
			}
			f5 = append(f5, *f5elem)
		}
		res.Tags = f5
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkUpdate")
	defer func() {
		exit(err)
	}()
	if delta.DifferentAt("Spec.Tags") {
		if err = rm.syncTags(ctx, desired, latest); err != nil {
			return nil, err
		}
	} else if !delta.DifferentExcept("Spec.Tags") {
		return desired, nil
	}

	input, err := rm.newUpdateRequestPayload(ctx, desired, delta)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.UpdateResolverEndpointOutput
	_ = resp
	resp, err = rm.sdkapi.UpdateResolverEndpoint(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateResolverEndpoint", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.ResolverEndpoint.Arn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.ResolverEndpoint.Arn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.ResolverEndpoint.CreationTime != nil {
		ko.Status.CreationTime = resp.ResolverEndpoint.CreationTime
	} else {
		ko.Status.CreationTime = nil
	}
	if resp.ResolverEndpoint.CreatorRequestId != nil {
		ko.Status.CreatorRequestID = resp.ResolverEndpoint.CreatorRequestId
	} else {
		ko.Status.CreatorRequestID = nil
	}
	if resp.ResolverEndpoint.Direction != "" {
		ko.Spec.Direction = aws.String(string(resp.ResolverEndpoint.Direction))
	} else {
		ko.Spec.Direction = nil
	}
	if resp.ResolverEndpoint.HostVPCId != nil {
		ko.Status.HostVPCID = resp.ResolverEndpoint.HostVPCId
	} else {
		ko.Status.HostVPCID = nil
	}
	if resp.ResolverEndpoint.Id != nil {
		ko.Status.ID = resp.ResolverEndpoint.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.ResolverEndpoint.IpAddressCount != nil {
		ipAddressCountCopy := int64(*resp.ResolverEndpoint.IpAddressCount)
		ko.Status.IPAddressCount = &ipAddressCountCopy
	} else {
		ko.Status.IPAddressCount = nil
	}
	if resp.ResolverEndpoint.ModificationTime != nil {
		ko.Status.ModificationTime = resp.ResolverEndpoint.ModificationTime
	} else {
		ko.Status.ModificationTime = nil
	}
	if resp.ResolverEndpoint.Name != nil {
		ko.Spec.Name = resp.ResolverEndpoint.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.ResolverEndpoint.ResolverEndpointType != "" {
		ko.Spec.ResolverEndpointType = aws.String(string(resp.ResolverEndpoint.ResolverEndpointType))
	} else {
		ko.Spec.ResolverEndpointType = nil
	}
	if resp.ResolverEndpoint.SecurityGroupIds != nil {
		ko.Spec.SecurityGroupIDs = aws.StringSlice(resp.ResolverEndpoint.SecurityGroupIds)
	} else {
		ko.Spec.SecurityGroupIDs = nil
	}
	if resp.ResolverEndpoint.Status != "" {
		ko.Status.Status = aws.String(string(resp.ResolverEndpoint.Status))
	} else {
		ko.Status.Status = nil
	}
	if resp.ResolverEndpoint.StatusMessage != nil {
		ko.Status.StatusMessage = resp.ResolverEndpoint.StatusMessage
	} else {
		ko.Status.StatusMessage = nil
	}

	rm.setStatusDefaults(ko)
	if delta.DifferentAt("Spec.IPAddresses") {
		rm.SyncIPAddresses(ctx, desired, latest)
		ko.Status.IPAddressCount = latest.ko.Status.IPAddressCount
	}
	return &resource{ko}, nil
}

// newUpdateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Update API call for the resource
func (rm *resourceManager) newUpdateRequestPayload(
	ctx context.Context,
	r *resource,
	delta *ackcompare.Delta,
) (*svcsdk.UpdateResolverEndpointInput, error) {
	res := &svcsdk.UpdateResolverEndpointInput{}

	if r.ko.Spec.Name != nil {
		res.Name = r.ko.Spec.Name
	}
	if r.ko.Status.ID != nil {
		res.ResolverEndpointId = r.ko.Status.ID
	}
	if r.ko.Spec.ResolverEndpointType != nil {
		res.ResolverEndpointType = svcsdktypes.ResolverEndpointType(*r.ko.Spec.ResolverEndpointType)
	}

	return res, nil
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteResolverEndpointOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteResolverEndpoint(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteResolverEndpoint", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteResolverEndpointInput, error) {
	res := &svcsdk.DeleteResolverEndpointInput{}

	if r.ko.Status.ID != nil {
		res.ResolverEndpointId = r.ko.Status.ID
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.ResolverEndpoint,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	// No terminal_errors specified for this resource in generator config
	return false
}
