package cmd

import "time"

type EventResp struct {
	NAMESPACE string `json:"namespace"`
	LastTimestamp time.Time `json:"lastTimestamp"`
	Type string `json:"type"`
	Reason string `json:"reason"`
	Object string `json:"object"`
	Message string `json:"message"`
}

type EventFlags struct {
	K8S
	ElasticSearch

}

type K8S struct {
	KubeConfig string
}

type ElasticSearch struct {
	Address string
	UserName string
	Password string
	CaCert string
}