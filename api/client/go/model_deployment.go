/*
 * SignalCD
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"time"
)

// Deployment struct for Deployment
type Deployment struct {
	Number   int64     `json:"number"`
	Created  time.Time `json:"created,omitempty"`
	Pipeline Pipeline  `json:"pipeline,omitempty"`
}
