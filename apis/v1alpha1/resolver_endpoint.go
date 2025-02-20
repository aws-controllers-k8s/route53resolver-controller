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

package v1alpha1

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResolverEndpointSpec defines the desired state of ResolverEndpoint.
//
// In the response to a CreateResolverEndpoint (https://docs.aws.amazon.com/Route53/latest/APIReference/API_route53resolver_CreateResolverEndpoint.html),
// DeleteResolverEndpoint (https://docs.aws.amazon.com/Route53/latest/APIReference/API_route53resolver_DeleteResolverEndpoint.html),
// GetResolverEndpoint (https://docs.aws.amazon.com/Route53/latest/APIReference/API_route53resolver_GetResolverEndpoint.html),
// Updates the name, or ResolverEndpointType for an endpoint, or UpdateResolverEndpoint
// (https://docs.aws.amazon.com/Route53/latest/APIReference/API_route53resolver_UpdateResolverEndpoint.html)
// request, a complex type that contains settings for an existing inbound or
// outbound Resolver endpoint.
type ResolverEndpointSpec struct {

	// Specify the applicable value:
	//
	//   - INBOUND: Resolver forwards DNS queries to the DNS service for a VPC
	//     from your network
	//
	//   - OUTBOUND: Resolver forwards DNS queries from the DNS service for a VPC
	//     to your network
	//
	// +kubebuilder:validation:Required
	Direction *string `json:"direction"`
	// The subnets and IP addresses in your VPC that DNS queries originate from
	// (for outbound endpoints) or that you forward DNS queries to (for inbound
	// endpoints). The subnet ID uniquely identifies a VPC.
	//
	// Even though the minimum is 1, Route 53 requires that you create at least
	// two.
	// +kubebuilder:validation:Required
	IPAddresses []*IPAddressRequest `json:"ipAddresses,omitempty"`
	// A friendly name that lets you easily find a configuration in the Resolver
	// dashboard in the Route 53 console.
	Name *string `json:"name,omitempty"`
	// For the endpoint type you can choose either IPv4, IPv6, or dual-stack. A
	// dual-stack endpoint means that it will resolve via both IPv4 and IPv6. This
	// endpoint type is applied to all IP addresses.
	ResolverEndpointType *string `json:"resolverEndpointType,omitempty"`
	// The ID of one or more security groups that you want to use to control access
	// to this VPC. The security group that you specify must include one or more
	// inbound rules (for inbound Resolver endpoints) or outbound rules (for outbound
	// Resolver endpoints). Inbound and outbound rules must allow TCP and UDP access.
	// For inbound access, open port 53. For outbound access, open the port that
	// you're using for DNS queries on your network.
	//
	// Some security group rules will cause your connection to be tracked. For outbound
	// resolver endpoint, it can potentially impact the maximum queries per second
	// from outbound endpoint to your target name server. For inbound resolver endpoint,
	// it can bring down the overall maximum queries per second per IP address to
	// as low as 1500. To avoid connection tracking caused by security group, see
	// Untracked connections (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/security-group-connection-tracking.html#untracked-connectionsl).
	SecurityGroupIDs  []*string                                  `json:"securityGroupIDs,omitempty"`
	SecurityGroupRefs []*ackv1alpha1.AWSResourceReferenceWrapper `json:"securityGroupRefs,omitempty"`
	// A list of the tag keys and values that you want to associate with the endpoint.
	Tags []*Tag `json:"tags,omitempty"`
}

// ResolverEndpointStatus defines the observed state of ResolverEndpoint
type ResolverEndpointStatus struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	// +kubebuilder:validation:Optional
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRs managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	// +kubebuilder:validation:Optional
	Conditions []*ackv1alpha1.Condition `json:"conditions"`
	// The date and time that the endpoint was created, in Unix time format and
	// Coordinated Universal Time (UTC).
	// +kubebuilder:validation:Optional
	CreationTime *string `json:"creationTime,omitempty"`
	// A unique string that identifies the request that created the Resolver endpoint.
	// The CreatorRequestId allows failed requests to be retried without the risk
	// of running the operation twice.
	// +kubebuilder:validation:Optional
	CreatorRequestID *string `json:"creatorRequestID"`
	// The ID of the VPC that you want to create the Resolver endpoint in.
	// +kubebuilder:validation:Optional
	HostVPCID *string `json:"hostVPCID,omitempty"`
	// +kubebuilder:validation:Optional
	IPAddresses []*IPAddressResponse `json:"ipAddresses,omitempty"`
	// The ID of the Resolver endpoint.
	// +kubebuilder:validation:Optional
	ID *string `json:"id,omitempty"`
	// The number of IP addresses that the Resolver endpoint can use for DNS queries.
	// +kubebuilder:validation:Optional
	IPAddressCount *int64 `json:"ipAddressCount,omitempty"`
	// The date and time that the endpoint was last modified, in Unix time format
	// and Coordinated Universal Time (UTC).
	// +kubebuilder:validation:Optional
	ModificationTime *string `json:"modificationTime,omitempty"`
	// A code that specifies the current status of the Resolver endpoint. Valid
	// values include the following:
	//
	//    * CREATING: Resolver is creating and configuring one or more Amazon VPC
	//    network interfaces for this endpoint.
	//
	//    * OPERATIONAL: The Amazon VPC network interfaces for this endpoint are
	//    correctly configured and able to pass inbound or outbound DNS queries
	//    between your network and Resolver.
	//
	//    * UPDATING: Resolver is associating or disassociating one or more network
	//    interfaces with this endpoint.
	//
	//    * AUTO_RECOVERING: Resolver is trying to recover one or more of the network
	//    interfaces that are associated with this endpoint. During the recovery
	//    process, the endpoint functions with limited capacity because of the limit
	//    on the number of DNS queries per IP address (per network interface). For
	//    the current limit, see Limits on Route 53 Resolver (https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/DNSLimitations.html#limits-api-entities-resolver).
	//
	//    * ACTION_NEEDED: This endpoint is unhealthy, and Resolver can't automatically
	//    recover it. To resolve the problem, we recommend that you check each IP
	//    address that you associated with the endpoint. For each IP address that
	//    isn't available, add another IP address and then delete the IP address
	//    that isn't available. (An endpoint must always include at least two IP
	//    addresses.) A status of ACTION_NEEDED can have a variety of causes. Here
	//    are two common causes: One or more of the network interfaces that are
	//    associated with the endpoint were deleted using Amazon VPC. The network
	//    interface couldn't be created for some reason that's outside the control
	//    of Resolver.
	//
	//    * DELETING: Resolver is deleting this endpoint and the associated network
	//    interfaces.
	// +kubebuilder:validation:Optional
	Status *string `json:"status,omitempty"`
	// A detailed description of the status of the Resolver endpoint.
	// +kubebuilder:validation:Optional
	StatusMessage *string `json:"statusMessage,omitempty"`
}

// ResolverEndpoint is the Schema for the ResolverEndpoints API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ID",type=string,priority=0,JSONPath=`.status.id`
type ResolverEndpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ResolverEndpointSpec   `json:"spec,omitempty"`
	Status            ResolverEndpointStatus `json:"status,omitempty"`
}

// ResolverEndpointList contains a list of ResolverEndpoint
// +kubebuilder:object:root=true
type ResolverEndpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResolverEndpoint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResolverEndpoint{}, &ResolverEndpointList{})
}
