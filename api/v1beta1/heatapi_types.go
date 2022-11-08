/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"fmt"
	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	"github.com/openstack-k8s-operators/lib-common/modules/common/endpoint"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HeatApiSpec defines the desired state of HeatApi
type HeatApiSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

        // +kubebuilder:validation:Optional
        // +kubebuilder:default=nova
        // APIMessageBusUser - username to use when accessing the API message bus
        MessageBusUser string `json:"messageBusUser"`

        // +kubebuilder:validation:Optional
        // APIMessageBusHostname - hostname to use when accessing the API message
        // bus. If not provided then upcalls will be disabled. This filed is
        // Required for cell0.
        // TODO(gibi): Add a webhook to validate cell0 constraint.
        MessageBusHostname string `json:"messageBusHostname"`

	// +kubebuilder:validation:Required
	// Secret is the name of the Secret instance containing password
	// information for the heat-api sevice.
	Secret string `json:"secret"`

	// +kubebuilder:validation:Optional
	// PasswordSelectors - Field names to identify the passwords from the
	// Secret
	PasswordSelectors PasswordSelector `json:"passwordSelectors,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="heat"
	// ServiceUser - optional username used for this service to register in
	// keystone
	ServiceUser string `json:"serviceUser,omitempty"`

	// +kubebuilder:validation:Required
	KeystoneAuthURL string `json:"keystoneAuthURL"`

	// +kubebuilder:validation:Required
	// NovaServiceBase specifies the generic fields of the service
	HeatServiceBase `json:",inline"`
}

// HeatApiStatus defines the observed state of HeatApi
type HeatApiStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ReadyCount of heat API instances
	ReadyCount int32 `json:"readyCount,omitempty"`

	// Map of hashes to track e.g. job status
	Hash map[string]string `json:"hash,omitempty"`

	// API endpoint
	APIEndpoints map[string]string `json:"apiEndpoint,omitempty"`

	// Conditions
	Conditions condition.Conditions `json:"conditions,omitempty" optional:"true"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HeatApi is the Schema for the heatapis API
type HeatApi struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HeatApiSpec   `json:"spec,omitempty"`
	Status HeatApiStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HeatApiList contains a list of HeatApi
type HeatApiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HeatApi `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HeatApi{}, &HeatApiList{})
}

// GetEndpoint - returns OpenStack endpoint url for type
func (instance HeatAPI) GetEndpoint(endpointType endpoint.Endpoint) (string, error) {
	if url, found := instance.Status.APIEndpoints[string(endpointType)]; found {
		return url, nil
	}
	return "", fmt.Errorf("%s endpoint not found", string(endpointType))
}

// IsReady - returns true if service is ready to server requests
func (instance HeatAPI) IsReady() bool {
	return instance.Status.Conditions.IsTrue(condition.ExposeServiceReadyCondition) &&
		instance.Status.Conditions.IsTrue(condition.DeploymentReadyCondition)
}
