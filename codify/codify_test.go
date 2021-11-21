//
// Copyright © 2021 Kris Nóva <kris@nivenly.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//   ███╗   ██╗ █████╗ ███╗   ███╗██╗
//   ████╗  ██║██╔══██╗████╗ ████║██║
//   ██╔██╗ ██║███████║██╔████╔██║██║
//   ██║╚██╗██║██╔══██║██║╚██╔╝██║██║
//   ██║ ╚████║██║  ██║██║ ╚═╝ ██║███████╗
//   ╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝
//

package codify

import (
	"strings"
	"testing"
)

// TestAliasSubstitutionCorrectDefault is the simplest test.
// This will check that a 1 to 1 substitution occurs, while
// we are using the correct default.
func TestAliasSubstitutionCorrectDefault(t *testing.T) {

	generated := "v1.Volume" // We know Volume is a corev1.Volume
	result := alias(generated, "corev1")
	if !strings.Contains(result, "corev1.Volume") {
		t.Errorf("missing expected corev1.Volume: %s", result)
	}
	generated = "SomethingSomething v1.Volume SomethingSomething"
	result = alias(generated, "corev1")
	if !strings.Contains(result, "corev1.Volume") {
		t.Errorf("missing expected corev1.Volume: %s", result)
	}

	generated = "v1.SadVolume" // We know SadVolume is unknown, but we expect a default anyway
	result = alias(generated, "corev1")
	if !strings.Contains(result, "corev1.SadVolume") {
		t.Errorf("unexpected substitution: %s", result)
	}

	generated = "" // Check empty string
	result = alias(generated, "corev1")
	if strings.Contains(result, "corev1") {
		t.Errorf("unexpected substitution: %s", result)
	}
}

// TestAliasEdgeCases will exercise some (but not all) cases
// where the host object has "sub objects" of different types.
func TestAliasEdgeCases(t *testing.T) {

	// We can assume
	// v1.Volume      should be   corev1.Volume
	// v1.ObjectMeta  should be   metav1.ObjectMeta

	generated := "Something v1.Volume Else v1.ObjectMeta Other v1.ThisShouldDefault"
	result := alias(generated, "appsv1")
	// here we build the EXACT string we are expecting!
	expected := "Something corev1.Volume Else metav1.ObjectMeta Other appsv1.ThisShouldDefault"
	if result != expected {
		t.Errorf("unexpected result")
		t.Errorf("expected: %s", expected)
		t.Errorf("result:   %s", result)
	}
}

// TestPolicyV1Beta1 will check for the special type "policyv1beta1"
//
// This is unique to RBAC, and is slightly different than the other
// alias checks within the Codify function.
//
// PolicyV1Interface
// policyv1beta1
//
func TestPolicyV1Beta1(t *testing.T) {
	generated := "v1.PolicyV1Interface"
	result := alias(generated, "appsv1")
	if !strings.Contains(result, "policyv1beta1.PolicyV1Interface") {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestCleanValast20open(t *testing.T) {
	input := `something{{`
	expected := `something{
{`
	actual := cleanValast20(input)
	if actual != expected {
		t.Errorf("unexpected cleanValast20 output: %s", actual)
	}
}

func TestCleanValast20close(t *testing.T) {
	input := `}}`
	expected := `},
}`
	actual := cleanValast20(input)
	if actual != expected {
		t.Errorf("unexpected cleanValast20 output: %s", actual)
	}
}
