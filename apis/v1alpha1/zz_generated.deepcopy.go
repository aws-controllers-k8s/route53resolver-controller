//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	corev1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Filter) DeepCopyInto(out *Filter) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]*string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(string)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Filter.
func (in *Filter) DeepCopy() *Filter {
	if in == nil {
		return nil
	}
	out := new(Filter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FirewallConfig) DeepCopyInto(out *FirewallConfig) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.OwnerID != nil {
		in, out := &in.OwnerID, &out.OwnerID
		*out = new(string)
		**out = **in
	}
	if in.ResourceID != nil {
		in, out := &in.ResourceID, &out.ResourceID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirewallConfig.
func (in *FirewallConfig) DeepCopy() *FirewallConfig {
	if in == nil {
		return nil
	}
	out := new(FirewallConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FirewallDomainList) DeepCopyInto(out *FirewallDomainList) {
	*out = *in
	if in.ARN != nil {
		in, out := &in.ARN, &out.ARN
		*out = new(string)
		**out = **in
	}
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = new(string)
		**out = **in
	}
	if in.CreatorRequestID != nil {
		in, out := &in.CreatorRequestID, &out.CreatorRequestID
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.ModificationTime != nil {
		in, out := &in.ModificationTime, &out.ModificationTime
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.StatusMessage != nil {
		in, out := &in.StatusMessage, &out.StatusMessage
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirewallDomainList.
func (in *FirewallDomainList) DeepCopy() *FirewallDomainList {
	if in == nil {
		return nil
	}
	out := new(FirewallDomainList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FirewallDomainListMetadata) DeepCopyInto(out *FirewallDomainListMetadata) {
	*out = *in
	if in.ARN != nil {
		in, out := &in.ARN, &out.ARN
		*out = new(string)
		**out = **in
	}
	if in.CreatorRequestID != nil {
		in, out := &in.CreatorRequestID, &out.CreatorRequestID
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirewallDomainListMetadata.
func (in *FirewallDomainListMetadata) DeepCopy() *FirewallDomainListMetadata {
	if in == nil {
		return nil
	}
	out := new(FirewallDomainListMetadata)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FirewallRule) DeepCopyInto(out *FirewallRule) {
	*out = *in
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = new(string)
		**out = **in
	}
	if in.CreatorRequestID != nil {
		in, out := &in.CreatorRequestID, &out.CreatorRequestID
		*out = new(string)
		**out = **in
	}
	if in.FirewallDomainListID != nil {
		in, out := &in.FirewallDomainListID, &out.FirewallDomainListID
		*out = new(string)
		**out = **in
	}
	if in.FirewallRuleGroupID != nil {
		in, out := &in.FirewallRuleGroupID, &out.FirewallRuleGroupID
		*out = new(string)
		**out = **in
	}
	if in.ModificationTime != nil {
		in, out := &in.ModificationTime, &out.ModificationTime
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirewallRule.
func (in *FirewallRule) DeepCopy() *FirewallRule {
	if in == nil {
		return nil
	}
	out := new(FirewallRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FirewallRuleGroup) DeepCopyInto(out *FirewallRuleGroup) {
	*out = *in
	if in.ARN != nil {
		in, out := &in.ARN, &out.ARN
		*out = new(string)
		**out = **in
	}
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = new(string)
		**out = **in
	}
	if in.CreatorRequestID != nil {
		in, out := &in.CreatorRequestID, &out.CreatorRequestID
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.ModificationTime != nil {
		in, out := &in.ModificationTime, &out.ModificationTime
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.OwnerID != nil {
		in, out := &in.OwnerID, &out.OwnerID
		*out = new(string)
		**out = **in
	}
	if in.ShareStatus != nil {
		in, out := &in.ShareStatus, &out.ShareStatus
		*out = new(string)
		**out = **in
	}
	if in.StatusMessage != nil {
		in, out := &in.StatusMessage, &out.StatusMessage
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirewallRuleGroup.
func (in *FirewallRuleGroup) DeepCopy() *FirewallRuleGroup {
	if in == nil {
		return nil
	}
	out := new(FirewallRuleGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FirewallRuleGroupAssociation) DeepCopyInto(out *FirewallRuleGroupAssociation) {
	*out = *in
	if in.ARN != nil {
		in, out := &in.ARN, &out.ARN
		*out = new(string)
		**out = **in
	}
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = new(string)
		**out = **in
	}
	if in.CreatorRequestID != nil {
		in, out := &in.CreatorRequestID, &out.CreatorRequestID
		*out = new(string)
		**out = **in
	}
	if in.FirewallRuleGroupID != nil {
		in, out := &in.FirewallRuleGroupID, &out.FirewallRuleGroupID
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.ModificationTime != nil {
		in, out := &in.ModificationTime, &out.ModificationTime
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.StatusMessage != nil {
		in, out := &in.StatusMessage, &out.StatusMessage
		*out = new(string)
		**out = **in
	}
	if in.VPCID != nil {
		in, out := &in.VPCID, &out.VPCID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirewallRuleGroupAssociation.
func (in *FirewallRuleGroupAssociation) DeepCopy() *FirewallRuleGroupAssociation {
	if in == nil {
		return nil
	}
	out := new(FirewallRuleGroupAssociation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FirewallRuleGroupMetadata) DeepCopyInto(out *FirewallRuleGroupMetadata) {
	*out = *in
	if in.ARN != nil {
		in, out := &in.ARN, &out.ARN
		*out = new(string)
		**out = **in
	}
	if in.CreatorRequestID != nil {
		in, out := &in.CreatorRequestID, &out.CreatorRequestID
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.OwnerID != nil {
		in, out := &in.OwnerID, &out.OwnerID
		*out = new(string)
		**out = **in
	}
	if in.ShareStatus != nil {
		in, out := &in.ShareStatus, &out.ShareStatus
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirewallRuleGroupMetadata.
func (in *FirewallRuleGroupMetadata) DeepCopy() *FirewallRuleGroupMetadata {
	if in == nil {
		return nil
	}
	out := new(FirewallRuleGroupMetadata)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPAddressRequest) DeepCopyInto(out *IPAddressRequest) {
	*out = *in
	if in.IP != nil {
		in, out := &in.IP, &out.IP
		*out = new(string)
		**out = **in
	}
	if in.IPv6 != nil {
		in, out := &in.IPv6, &out.IPv6
		*out = new(string)
		**out = **in
	}
	if in.SubnetID != nil {
		in, out := &in.SubnetID, &out.SubnetID
		*out = new(string)
		**out = **in
	}
	if in.SubnetRef != nil {
		in, out := &in.SubnetRef, &out.SubnetRef
		*out = new(corev1alpha1.AWSResourceReferenceWrapper)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPAddressRequest.
func (in *IPAddressRequest) DeepCopy() *IPAddressRequest {
	if in == nil {
		return nil
	}
	out := new(IPAddressRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPAddressResponse) DeepCopyInto(out *IPAddressResponse) {
	*out = *in
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = new(string)
		**out = **in
	}
	if in.IP != nil {
		in, out := &in.IP, &out.IP
		*out = new(string)
		**out = **in
	}
	if in.IPID != nil {
		in, out := &in.IPID, &out.IPID
		*out = new(string)
		**out = **in
	}
	if in.IPv6 != nil {
		in, out := &in.IPv6, &out.IPv6
		*out = new(string)
		**out = **in
	}
	if in.ModificationTime != nil {
		in, out := &in.ModificationTime, &out.ModificationTime
		*out = new(string)
		**out = **in
	}
	if in.StatusMessage != nil {
		in, out := &in.StatusMessage, &out.StatusMessage
		*out = new(string)
		**out = **in
	}
	if in.SubnetID != nil {
		in, out := &in.SubnetID, &out.SubnetID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPAddressResponse.
func (in *IPAddressResponse) DeepCopy() *IPAddressResponse {
	if in == nil {
		return nil
	}
	out := new(IPAddressResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPAddressUpdate) DeepCopyInto(out *IPAddressUpdate) {
	*out = *in
	if in.IP != nil {
		in, out := &in.IP, &out.IP
		*out = new(string)
		**out = **in
	}
	if in.IPID != nil {
		in, out := &in.IPID, &out.IPID
		*out = new(string)
		**out = **in
	}
	if in.IPv6 != nil {
		in, out := &in.IPv6, &out.IPv6
		*out = new(string)
		**out = **in
	}
	if in.SubnetID != nil {
		in, out := &in.SubnetID, &out.SubnetID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPAddressUpdate.
func (in *IPAddressUpdate) DeepCopy() *IPAddressUpdate {
	if in == nil {
		return nil
	}
	out := new(IPAddressUpdate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverConfig) DeepCopyInto(out *ResolverConfig) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.OwnerID != nil {
		in, out := &in.OwnerID, &out.OwnerID
		*out = new(string)
		**out = **in
	}
	if in.ResourceID != nil {
		in, out := &in.ResourceID, &out.ResourceID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverConfig.
func (in *ResolverConfig) DeepCopy() *ResolverConfig {
	if in == nil {
		return nil
	}
	out := new(ResolverConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverDNSSEcConfig) DeepCopyInto(out *ResolverDNSSEcConfig) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.OwnerID != nil {
		in, out := &in.OwnerID, &out.OwnerID
		*out = new(string)
		**out = **in
	}
	if in.ResourceID != nil {
		in, out := &in.ResourceID, &out.ResourceID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverDNSSEcConfig.
func (in *ResolverDNSSEcConfig) DeepCopy() *ResolverDNSSEcConfig {
	if in == nil {
		return nil
	}
	out := new(ResolverDNSSEcConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverEndpoint) DeepCopyInto(out *ResolverEndpoint) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverEndpoint.
func (in *ResolverEndpoint) DeepCopy() *ResolverEndpoint {
	if in == nil {
		return nil
	}
	out := new(ResolverEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResolverEndpoint) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverEndpointList) DeepCopyInto(out *ResolverEndpointList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResolverEndpoint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverEndpointList.
func (in *ResolverEndpointList) DeepCopy() *ResolverEndpointList {
	if in == nil {
		return nil
	}
	out := new(ResolverEndpointList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResolverEndpointList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverEndpointSpec) DeepCopyInto(out *ResolverEndpointSpec) {
	*out = *in
	if in.Direction != nil {
		in, out := &in.Direction, &out.Direction
		*out = new(string)
		**out = **in
	}
	if in.IPAddresses != nil {
		in, out := &in.IPAddresses, &out.IPAddresses
		*out = make([]*IPAddressRequest, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(IPAddressRequest)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.ResolverEndpointType != nil {
		in, out := &in.ResolverEndpointType, &out.ResolverEndpointType
		*out = new(string)
		**out = **in
	}
	if in.SecurityGroupIDs != nil {
		in, out := &in.SecurityGroupIDs, &out.SecurityGroupIDs
		*out = make([]*string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(string)
				**out = **in
			}
		}
	}
	if in.SecurityGroupRefs != nil {
		in, out := &in.SecurityGroupRefs, &out.SecurityGroupRefs
		*out = make([]*corev1alpha1.AWSResourceReferenceWrapper, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(corev1alpha1.AWSResourceReferenceWrapper)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make([]*Tag, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Tag)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverEndpointSpec.
func (in *ResolverEndpointSpec) DeepCopy() *ResolverEndpointSpec {
	if in == nil {
		return nil
	}
	out := new(ResolverEndpointSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverEndpointStatus) DeepCopyInto(out *ResolverEndpointStatus) {
	*out = *in
	if in.ACKResourceMetadata != nil {
		in, out := &in.ACKResourceMetadata, &out.ACKResourceMetadata
		*out = new(corev1alpha1.ResourceMetadata)
		(*in).DeepCopyInto(*out)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]*corev1alpha1.Condition, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(corev1alpha1.Condition)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = new(string)
		**out = **in
	}
	if in.CreatorRequestID != nil {
		in, out := &in.CreatorRequestID, &out.CreatorRequestID
		*out = new(string)
		**out = **in
	}
	if in.HostVPCID != nil {
		in, out := &in.HostVPCID, &out.HostVPCID
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.IPAddressCount != nil {
		in, out := &in.IPAddressCount, &out.IPAddressCount
		*out = new(int64)
		**out = **in
	}
	if in.ModificationTime != nil {
		in, out := &in.ModificationTime, &out.ModificationTime
		*out = new(string)
		**out = **in
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(string)
		**out = **in
	}
	if in.StatusMessage != nil {
		in, out := &in.StatusMessage, &out.StatusMessage
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverEndpointStatus.
func (in *ResolverEndpointStatus) DeepCopy() *ResolverEndpointStatus {
	if in == nil {
		return nil
	}
	out := new(ResolverEndpointStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverEndpoint_SDK) DeepCopyInto(out *ResolverEndpoint_SDK) {
	*out = *in
	if in.ARN != nil {
		in, out := &in.ARN, &out.ARN
		*out = new(string)
		**out = **in
	}
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = new(string)
		**out = **in
	}
	if in.CreatorRequestID != nil {
		in, out := &in.CreatorRequestID, &out.CreatorRequestID
		*out = new(string)
		**out = **in
	}
	if in.Direction != nil {
		in, out := &in.Direction, &out.Direction
		*out = new(string)
		**out = **in
	}
	if in.HostVPCID != nil {
		in, out := &in.HostVPCID, &out.HostVPCID
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.IPAddressCount != nil {
		in, out := &in.IPAddressCount, &out.IPAddressCount
		*out = new(int64)
		**out = **in
	}
	if in.ModificationTime != nil {
		in, out := &in.ModificationTime, &out.ModificationTime
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.ResolverEndpointType != nil {
		in, out := &in.ResolverEndpointType, &out.ResolverEndpointType
		*out = new(string)
		**out = **in
	}
	if in.SecurityGroupIDs != nil {
		in, out := &in.SecurityGroupIDs, &out.SecurityGroupIDs
		*out = make([]*string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(string)
				**out = **in
			}
		}
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(string)
		**out = **in
	}
	if in.StatusMessage != nil {
		in, out := &in.StatusMessage, &out.StatusMessage
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverEndpoint_SDK.
func (in *ResolverEndpoint_SDK) DeepCopy() *ResolverEndpoint_SDK {
	if in == nil {
		return nil
	}
	out := new(ResolverEndpoint_SDK)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverQueryLogConfig) DeepCopyInto(out *ResolverQueryLogConfig) {
	*out = *in
	if in.ARN != nil {
		in, out := &in.ARN, &out.ARN
		*out = new(string)
		**out = **in
	}
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = new(string)
		**out = **in
	}
	if in.CreatorRequestID != nil {
		in, out := &in.CreatorRequestID, &out.CreatorRequestID
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.OwnerID != nil {
		in, out := &in.OwnerID, &out.OwnerID
		*out = new(string)
		**out = **in
	}
	if in.ShareStatus != nil {
		in, out := &in.ShareStatus, &out.ShareStatus
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverQueryLogConfig.
func (in *ResolverQueryLogConfig) DeepCopy() *ResolverQueryLogConfig {
	if in == nil {
		return nil
	}
	out := new(ResolverQueryLogConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverQueryLogConfigAssociation) DeepCopyInto(out *ResolverQueryLogConfigAssociation) {
	*out = *in
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.ResolverQueryLogConfigID != nil {
		in, out := &in.ResolverQueryLogConfigID, &out.ResolverQueryLogConfigID
		*out = new(string)
		**out = **in
	}
	if in.ResourceID != nil {
		in, out := &in.ResourceID, &out.ResourceID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverQueryLogConfigAssociation.
func (in *ResolverQueryLogConfigAssociation) DeepCopy() *ResolverQueryLogConfigAssociation {
	if in == nil {
		return nil
	}
	out := new(ResolverQueryLogConfigAssociation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverRule) DeepCopyInto(out *ResolverRule) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverRule.
func (in *ResolverRule) DeepCopy() *ResolverRule {
	if in == nil {
		return nil
	}
	out := new(ResolverRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResolverRule) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverRuleAssociation) DeepCopyInto(out *ResolverRuleAssociation) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.ResolverRuleID != nil {
		in, out := &in.ResolverRuleID, &out.ResolverRuleID
		*out = new(string)
		**out = **in
	}
	if in.StatusMessage != nil {
		in, out := &in.StatusMessage, &out.StatusMessage
		*out = new(string)
		**out = **in
	}
	if in.VPCID != nil {
		in, out := &in.VPCID, &out.VPCID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverRuleAssociation.
func (in *ResolverRuleAssociation) DeepCopy() *ResolverRuleAssociation {
	if in == nil {
		return nil
	}
	out := new(ResolverRuleAssociation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverRuleConfig) DeepCopyInto(out *ResolverRuleConfig) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.ResolverEndpointID != nil {
		in, out := &in.ResolverEndpointID, &out.ResolverEndpointID
		*out = new(string)
		**out = **in
	}
	if in.TargetIPs != nil {
		in, out := &in.TargetIPs, &out.TargetIPs
		*out = make([]*TargetAddress, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(TargetAddress)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverRuleConfig.
func (in *ResolverRuleConfig) DeepCopy() *ResolverRuleConfig {
	if in == nil {
		return nil
	}
	out := new(ResolverRuleConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverRuleList) DeepCopyInto(out *ResolverRuleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResolverRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverRuleList.
func (in *ResolverRuleList) DeepCopy() *ResolverRuleList {
	if in == nil {
		return nil
	}
	out := new(ResolverRuleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResolverRuleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverRuleSpec) DeepCopyInto(out *ResolverRuleSpec) {
	*out = *in
	if in.DomainName != nil {
		in, out := &in.DomainName, &out.DomainName
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.ResolverEndpointID != nil {
		in, out := &in.ResolverEndpointID, &out.ResolverEndpointID
		*out = new(string)
		**out = **in
	}
	if in.RuleType != nil {
		in, out := &in.RuleType, &out.RuleType
		*out = new(string)
		**out = **in
	}
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make([]*Tag, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Tag)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.TargetIPs != nil {
		in, out := &in.TargetIPs, &out.TargetIPs
		*out = make([]*TargetAddress, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(TargetAddress)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverRuleSpec.
func (in *ResolverRuleSpec) DeepCopy() *ResolverRuleSpec {
	if in == nil {
		return nil
	}
	out := new(ResolverRuleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverRuleStatus) DeepCopyInto(out *ResolverRuleStatus) {
	*out = *in
	if in.ACKResourceMetadata != nil {
		in, out := &in.ACKResourceMetadata, &out.ACKResourceMetadata
		*out = new(corev1alpha1.ResourceMetadata)
		(*in).DeepCopyInto(*out)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]*corev1alpha1.Condition, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(corev1alpha1.Condition)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = new(string)
		**out = **in
	}
	if in.CreatorRequestID != nil {
		in, out := &in.CreatorRequestID, &out.CreatorRequestID
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.ModificationTime != nil {
		in, out := &in.ModificationTime, &out.ModificationTime
		*out = new(string)
		**out = **in
	}
	if in.OwnerID != nil {
		in, out := &in.OwnerID, &out.OwnerID
		*out = new(string)
		**out = **in
	}
	if in.ShareStatus != nil {
		in, out := &in.ShareStatus, &out.ShareStatus
		*out = new(string)
		**out = **in
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(string)
		**out = **in
	}
	if in.StatusMessage != nil {
		in, out := &in.StatusMessage, &out.StatusMessage
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverRuleStatus.
func (in *ResolverRuleStatus) DeepCopy() *ResolverRuleStatus {
	if in == nil {
		return nil
	}
	out := new(ResolverRuleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResolverRule_SDK) DeepCopyInto(out *ResolverRule_SDK) {
	*out = *in
	if in.ARN != nil {
		in, out := &in.ARN, &out.ARN
		*out = new(string)
		**out = **in
	}
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = new(string)
		**out = **in
	}
	if in.CreatorRequestID != nil {
		in, out := &in.CreatorRequestID, &out.CreatorRequestID
		*out = new(string)
		**out = **in
	}
	if in.DomainName != nil {
		in, out := &in.DomainName, &out.DomainName
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.ModificationTime != nil {
		in, out := &in.ModificationTime, &out.ModificationTime
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.OwnerID != nil {
		in, out := &in.OwnerID, &out.OwnerID
		*out = new(string)
		**out = **in
	}
	if in.ResolverEndpointID != nil {
		in, out := &in.ResolverEndpointID, &out.ResolverEndpointID
		*out = new(string)
		**out = **in
	}
	if in.RuleType != nil {
		in, out := &in.RuleType, &out.RuleType
		*out = new(string)
		**out = **in
	}
	if in.ShareStatus != nil {
		in, out := &in.ShareStatus, &out.ShareStatus
		*out = new(string)
		**out = **in
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(string)
		**out = **in
	}
	if in.StatusMessage != nil {
		in, out := &in.StatusMessage, &out.StatusMessage
		*out = new(string)
		**out = **in
	}
	if in.TargetIPs != nil {
		in, out := &in.TargetIPs, &out.TargetIPs
		*out = make([]*TargetAddress, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(TargetAddress)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolverRule_SDK.
func (in *ResolverRule_SDK) DeepCopy() *ResolverRule_SDK {
	if in == nil {
		return nil
	}
	out := new(ResolverRule_SDK)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Tag) DeepCopyInto(out *Tag) {
	*out = *in
	if in.Key != nil {
		in, out := &in.Key, &out.Key
		*out = new(string)
		**out = **in
	}
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Tag.
func (in *Tag) DeepCopy() *Tag {
	if in == nil {
		return nil
	}
	out := new(Tag)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetAddress) DeepCopyInto(out *TargetAddress) {
	*out = *in
	if in.IP != nil {
		in, out := &in.IP, &out.IP
		*out = new(string)
		**out = **in
	}
	if in.IPv6 != nil {
		in, out := &in.IPv6, &out.IPv6
		*out = new(string)
		**out = **in
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetAddress.
func (in *TargetAddress) DeepCopy() *TargetAddress {
	if in == nil {
		return nil
	}
	out := new(TargetAddress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UpdateIPAddress) DeepCopyInto(out *UpdateIPAddress) {
	*out = *in
	if in.IPID != nil {
		in, out := &in.IPID, &out.IPID
		*out = new(string)
		**out = **in
	}
	if in.IPv6 != nil {
		in, out := &in.IPv6, &out.IPv6
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UpdateIPAddress.
func (in *UpdateIPAddress) DeepCopy() *UpdateIPAddress {
	if in == nil {
		return nil
	}
	out := new(UpdateIPAddress)
	in.DeepCopyInto(out)
	return out
}
