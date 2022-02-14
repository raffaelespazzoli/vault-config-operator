/*
Copyright 2021.

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

package v1alpha1

import (
	"context"
	"reflect"
	"strings"

	"github.com/redhat-cop/operator-utils/pkg/util/apis"
	vaultutils "github.com/redhat-cop/vault-config-operator/api/v1alpha1/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ArtifactorySecretEngineRoleSpec defines the desired state of ArtifactorySecretEngineRole
type ArtifactorySecretEngineRoleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Authentication is the kube aoth configuraiton to be used to execute this request
	// +kubebuilder:validation:Required
	Authentication KubeAuthConfiguration `json:"authentication,omitempty"`

	// Path at which to create the role.
	// The final path will be {[spec.authentication.namespace]}/{spec.path}/roles/{metadata.name}.
	// The authentication role must have the following capabilities = [ "create", "read", "update", "delete"] on that path.
	// +kubebuilder:validation:Required
	Path Path `json:"path,omitempty"`

	ArtifactoryRole `json:",inline"`
}

type ArtifactoryRole struct {
	// GrantType Optional. Defaults to 'client_credentials' when creating the access token. You likely don't need to change this.
	// +kubebuilder:validation:Optional
	GrantType string `json:"grantType,omitempty"`

	// Username Optional The username for which the access token is created. If the user does not exist, Artifactory will create a transient user. Note that non-admininstrative access tokens can only create tokens for themselves.
	// +kubebuilder:validation:Optional
	Username string `json:"username,omitempty"`

	// Scope Required. See the JFrog Artifactory REST documentation on "Create Token" for a full and up to date description (https://www.jfrog.com/confluence/display/JFROG/Artifactory+REST+API#ArtifactoryRESTAPI-CreateToken.1).
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems:=1
	// +listType:=set
	Scopes []string `json:"scope,omitempty"`

	// Audience Optional. See the JFrog Artifactory REST documentation on "Create Token" for a full and up to date description(https://www.jfrog.com/confluence/display/JFROG/Artifactory+REST+API#ArtifactoryRESTAPI-CreateToken.1).
	// +kubebuilder:validation:Optional
	Audience string `json:"audience,omitempty"`

	// DefaultTTL Default TTL for issued access tokens. If unset, uses the backend's default_ttl. Cannot exceed max_ttl.
	// +kubebuilder:validation:Optional
	DefaultTTL metav1.Duration `json:"defaultTTL,omitempty"`

	// MaxTTL Maximum TTL that an access token can be renewed for. If unset, uses the backend's max_ttl. Cannot exceed backend's max_ttl.
	// +kubebuilder:validation:Optional
	MaxTTL metav1.Duration `json:"maxTTL,omitempty"`
}

func (i *ArtifactoryRole) toMap() map[string]interface{} {
	payload := map[string]interface{}{}
	payload["grant_type"] = i.GrantType
	payload["username"] = i.Username
	payload["scope"] = strings.Join(i.Scopes, ",")
	payload["audience"] = i.Audience
	payload["default_ttl"] = i.DefaultTTL
	payload["max_ttl"] = i.MaxTTL
	return payload
}

func (d *ArtifactorySecretEngineRole) GetPath() string {
	return string(d.Spec.Path) + "/" + "roles" + "/" + d.Name
}
func (d *ArtifactorySecretEngineRole) GetPayload() map[string]interface{} {
	return d.Spec.ArtifactoryRole.toMap()
}
func (d *ArtifactorySecretEngineRole) IsEquivalentToDesiredState(payload map[string]interface{}) bool {
	desiredState := d.Spec.ArtifactoryRole.toMap()
	return reflect.DeepEqual(desiredState, payload)
}

func (d *ArtifactorySecretEngineRole) IsInitialized() bool {
	return true
}

func (d *ArtifactorySecretEngineRole) PrepareInternalValues(context context.Context, object client.Object) error {
	return nil
}

func (r *ArtifactorySecretEngineRole) IsValid() (bool, error) {
	return true, nil
}

var _ vaultutils.VaultObject = &ArtifactorySecretEngineRole{}
var _ apis.ConditionsAware = &ArtifactorySecretEngineRole{}

// ArtifactorySecretEngineRoleStatus defines the observed state of ArtifactorySecretEngineRole
type ArtifactorySecretEngineRoleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

func (m *ArtifactorySecretEngineRole) GetConditions() []metav1.Condition {
	return m.Status.Conditions
}

func (m *ArtifactorySecretEngineRole) SetConditions(conditions []metav1.Condition) {
	m.Status.Conditions = conditions
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ArtifactorySecretEngineRole is the Schema for the artifactorysecretengineroles API
type ArtifactorySecretEngineRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArtifactorySecretEngineRoleSpec   `json:"spec,omitempty"`
	Status ArtifactorySecretEngineRoleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ArtifactorySecretEngineRoleList contains a list of ArtifactorySecretEngineRole
type ArtifactorySecretEngineRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ArtifactorySecretEngineRole `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ArtifactorySecretEngineRole{}, &ArtifactorySecretEngineRoleList{})
}
