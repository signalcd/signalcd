/*
 * SignalCD
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"time"
)

type DeploymentStepStatus struct {
	Phase string `json:"phase"`

	Started time.Time `json:"started"`

	Stopped time.Time `json:"stopped,omitempty"`
}