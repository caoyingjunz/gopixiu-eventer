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
