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
	"errors"
	"reflect"

	"github.com/redhat-cop/operator-utils/pkg/util/apis"
	vaultutils "github.com/redhat-cop/vault-config-operator/api/v1alpha1/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ArtifactorySecretEngineConfigSpec defines the desired state of ArtifactorySecretEngineConfig
type ArtifactorySecretEngineConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Authentication is the kube aoth configuraiton to be used to execute this request
	// +kubebuilder:validation:Required
	Authentication KubeAuthConfiguration `json:"authentication,omitempty"`

	// Path at which to create the role.
	// The final path will be {[.spec.authentication.namespace]}/{.spec.path}/{.metadata.name}/config/admin.
	// The authentication role must have the following capabilities = [ "create", "read", "update", "delete"] on that path.
	// +kubebuilder:validation:Required
	Path Path `json:"path,omitempty"`

	ArtifactoryConfig `json:",inline"`

	// RootCredentials specifies how to retrieve the admin token for this Artifactory connection.
	// Only Secret or Vault Secret can be used.
	// +kubebuilder:validation:Required
	RootCredentials RootCredentialConfig `json:"rootCredentials,omitempty"`
}

var _ vaultutils.VaultObject = &ArtifactorySecretEngineConfig{}
var _ apis.ConditionsAware = &ArtifactorySecretEngineConfig{}

type ArtifactoryConfig struct {
	// URL is the Artifactry URL
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^(http|https):\/\/.+$`
	URL string `json:"url,omitempty"`

	retrievedToken string `json:"-"`
}

func (d *ArtifactorySecretEngineConfig) GetPath() string {
	return string(d.Spec.Path) + "/" + d.Name + "/" + "config/admin"
}
func (d *ArtifactorySecretEngineConfig) GetPayload() map[string]interface{} {
	return d.Spec.ArtifactoryConfig.toMap()
}
func (d *ArtifactorySecretEngineConfig) IsEquivalentToDesiredState(payload map[string]interface{}) bool {
	desiredState := d.Spec.ArtifactoryConfig.toMap()
	delete(desiredState, "access_token")
	return reflect.DeepEqual(desiredState, payload)
}

func (i *ArtifactoryConfig) toMap() map[string]interface{} {
	payload := map[string]interface{}{}
	payload["url"] = i.URL
	payload["access_token"] = i.retrievedToken
	return payload
}

func (d *ArtifactorySecretEngineConfig) IsInitialized() bool {
	return true
}

func (d *ArtifactorySecretEngineConfig) PrepareInternalValues(context context.Context, object client.Object) error {
	return d.setInternalCredentials(context)
}

func (r *ArtifactorySecretEngineConfig) IsValid() (bool, error) {
	err := r.isValid()
	return err == nil, err
}

func (r *ArtifactorySecretEngineConfig) isValid() error {
	return r.Spec.RootCredentials.validateEitherFromVaultSecretOrFromSecret()
}

func (r *ArtifactorySecretEngineConfig) SetToken(token string) {
	r.Spec.retrievedToken = token
}

func (r *ArtifactorySecretEngineConfig) setInternalCredentials(context context.Context) error {
	log := log.FromContext(context)
	k8sClient := context.Value("kubeClient").(client.Client)
	if r.Spec.RootCredentials.Secret != nil {
		secret := &corev1.Secret{}
		err := k8sClient.Get(context, types.NamespacedName{
			Namespace: r.Namespace,
			Name:      r.Spec.RootCredentials.Secret.Name,
		}, secret)
		if err != nil {
			log.Error(err, "unable to retrieve Secret", "instance", r)
			return err
		}
		r.SetToken(string(secret.Data[r.Spec.RootCredentials.PasswordKey]))
		return nil
	}
	if r.Spec.RootCredentials.VaultSecret != nil {
		secret, exists, err := vaultutils.ReadSecret(context, string(r.Spec.RootCredentials.VaultSecret.Path))
		if err != nil {
			return err
		}
		if !exists {
			err = errors.New("secret not found")
			log.Error(err, "unable to retrieve vault secret", "instance", r)
			return err
		}

		r.SetToken(secret.Data[r.Spec.RootCredentials.PasswordKey].(string))
		return nil
	}
	return errors.New("no means of retrieving a secret was specified")
}

// ArtifactorySecretEngineConfigStatus defines the observed state of ArtifactorySecretEngineConfig
type ArtifactorySecretEngineConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

func (m *ArtifactorySecretEngineConfig) GetConditions() []metav1.Condition {
	return m.Status.Conditions
}

func (m *ArtifactorySecretEngineConfig) SetConditions(conditions []metav1.Condition) {
	m.Status.Conditions = conditions
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ArtifactorySecretEngineConfig is the Schema for the artifactorysecretengineconfigs API
type ArtifactorySecretEngineConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArtifactorySecretEngineConfigSpec   `json:"spec,omitempty"`
	Status ArtifactorySecretEngineConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ArtifactorySecretEngineConfigList contains a list of ArtifactorySecretEngineConfig
type ArtifactorySecretEngineConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ArtifactorySecretEngineConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ArtifactorySecretEngineConfig{}, &ArtifactorySecretEngineConfigList{})
}
