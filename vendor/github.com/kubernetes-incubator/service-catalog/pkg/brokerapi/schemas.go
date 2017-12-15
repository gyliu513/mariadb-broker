/*
Copyright 2016 The Kubernetes Authors.

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

package brokerapi

// Schemas represents a broker's schemas for both service instances and service
// bindings
type Schemas struct {
	Instance Schema `json:"instance"`
	Binding  Schema `json:"binding"`
}

// Schema consists of the schema for inputs and the schema for outputs.
// Schemas are in the form of JSON Schema v4 (http://json-schema.org/).
type Schema struct {
	Inputs  string `json:"inputs"`
	Outputs string `json:"outputs"`
}
