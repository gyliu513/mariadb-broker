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

package validation

import (
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"

	sc "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog"
)

// ValidateBrokerName is the validation function for Broker names.
var ValidateBrokerName = apivalidation.NameIsDNSSubdomain

// ValidateBroker implements the validation rules for a BrokerResource.
func ValidateBroker(broker *sc.Broker) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs,
		apivalidation.ValidateObjectMeta(&broker.ObjectMeta,
			false, /* namespace required */
			ValidateBrokerName,
			field.NewPath("metadata"))...)

	allErrs = append(allErrs, validateBrokerSpec(&broker.Spec, field.NewPath("spec"))...)
	return allErrs
}

func validateBrokerSpec(spec *sc.BrokerSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if "" == spec.URL {
		allErrs = append(allErrs,
			field.Required(fldPath.Child("url"),
				"brokers must have a remote url to contact"))
	}

	if spec.AuthSecret != nil {
		for _, msg := range apivalidation.ValidateNamespaceName(spec.AuthSecret.Namespace, false /* prefix */) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("authSecret", "namespace"), spec.AuthSecret.Namespace, msg))
		}

		for _, msg := range apivalidation.NameIsDNSSubdomain(spec.AuthSecret.Name, false /* prefix */) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("authSecret", "name"), spec.AuthSecret.Name, msg))
		}
	}

	return allErrs
}

// ValidateBrokerUpdate checks that when changing from an older broker to a newer broker is okay ?
func ValidateBrokerUpdate(new *sc.Broker, old *sc.Broker) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, ValidateBroker(new)...)
	allErrs = append(allErrs, ValidateBroker(old)...)
	return allErrs
}

// ValidateBrokerStatusUpdate checks that when changing from an older broker to a newer broker is okay.
func ValidateBrokerStatusUpdate(new *sc.Broker, old *sc.Broker) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, ValidateBrokerUpdate(new, old)...)
	return allErrs
}
