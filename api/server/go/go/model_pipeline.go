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

type Pipeline struct {
	Id string `json:"id"`

	Name string `json:"name,omitempty"`

	Created time.Time `json:"created,omitempty"`
}