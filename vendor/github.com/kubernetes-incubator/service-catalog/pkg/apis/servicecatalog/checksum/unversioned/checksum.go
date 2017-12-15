/*
Copyright 2017 The Kubernetes Authors.

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

package unversioned

import (
	"crypto/sha256"
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog"
)

// InstanceSpecChecksum calculates a checksum of the given InstanceSpec based
// on the following fields;
//
// - ServiceClassName
// - PlanName
// - Parameters
// - OSBGUID
func InstanceSpecChecksum(spec servicecatalog.InstanceSpec) string {
	specString := ""
	specString += fmt.Sprintf("serviceClassName: %v\n", spec.ServiceClassName)
	specString += fmt.Sprintf("planName: %v\n", spec.PlanName)

	if spec.Parameters != nil {
		specString += fmt.Sprintf("parameters:\n\n%v\n\n", string(spec.Parameters.Raw))
	}

	specString += fmt.Sprintf("osbGuid: %v\n", spec.OSBGUID)

	sum := sha256.Sum256([]byte(specString))
	return fmt.Sprintf("%x", sum)
}

// BindingSpecChecksum calculates a checksum of the given BindingSpec based on
// the following fields:
//
// - InstanceRef.Name
// - Parameters
// - OSBGUID
func BindingSpecChecksum(spec servicecatalog.BindingSpec) string {
	specString := ""
	specString += fmt.Sprintf("instanceRef: %v\n", spec.InstanceRef.Name)

	if spec.Parameters != nil {
		specString += fmt.Sprintf("parameters:\n\n%v\n\n", string(spec.Parameters.Raw))
	}

	specString += fmt.Sprintf("osbGuid: %v\n", spec.OSBGUID)

	sum := sha256.Sum256([]byte(specString))
	return fmt.Sprintf("%x", sum)
}
