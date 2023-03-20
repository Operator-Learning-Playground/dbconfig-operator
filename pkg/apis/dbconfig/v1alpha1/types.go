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
	Dsn    		string     `json:"dsn"`
	MaxIdleConn int  	   `json:"maxIdleConn"`
	Services    []Services  `json:"services"`
}

type Services struct {
	Service Service   `json:"service"`
}

type Service struct {
	Dbname   string   `json:"dbname"`
	Tables   string   `json:"tables"`
	User     string	  `json:"user"`	
	Password string   `json:"password"`  
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DbConfigList
type DbConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []DbConfig `json:"items"`
}


