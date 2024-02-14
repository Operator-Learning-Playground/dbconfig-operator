package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DbConfig
type DbConfig struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DbConfigSpec `json:"spec,omitempty"`
}

type DbConfigSpec struct {
	Dsn         string    `json:"dsn"`
	MaxIdleConn int       `json:"maxIdleConn" default:"10"`
	MaxOpenConn int       `json:"maxOpenConn" default:"100"`
	Services    []Service `json:"services"`
}

type Service struct {
	Dbname   string   `json:"dbname"`
	User     string   `json:"user"`
	Tables   Tables   `json:"tables"`
	Password Password `json:"password"`
	ReBuild  bool     `json:"rebuild" default:"false"`
}

type Tables struct {
	ConfigMapRef string `json:"configMapRef"`
}

type Password struct {
	SecretRef string `json:"secretRef"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DbConfigList
type DbConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []DbConfig `json:"items"`
}
