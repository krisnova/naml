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
	"fmt"
	"github.com/hexops/valast"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"strings"
)

func Literal(kubeobject interface{}) string {
	l := valast.String(kubeobject)
	return l
}

// cleanObjectMeta helps us get rid of things like timestamps
// by only "opting in" to certain fields.
func cleanObjectMeta(m metav1.ObjectMeta) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:                       m.Name,
		Namespace:                  m.Namespace,
		Labels:                     m.Labels,
		Annotations:                m.Annotations,
		ClusterName:                m.ClusterName,
		ResourceVersion:            m.ResourceVersion,
		Finalizers:                 m.Finalizers,
		Generation:                 m.Generation,
		GenerateName:               m.GenerateName,
		UID:                        m.UID,
		ManagedFields:              m.ManagedFields,
		OwnerReferences:            m.OwnerReferences,
		DeletionGracePeriodSeconds: m.DeletionGracePeriodSeconds,
	}
}

func alias(generated, defaultalias string) string {
	aliased := generated

	// default "corev1"
	aliased = strings.Replace(aliased, "v1", defaultalias, -1)

	// ------------------------------
	// [ appsv1 ]
	appsv1types := []string{}

	for _, t := range appsv1types {
		aliased = strings.Replace(aliased,
			fmt.Sprintf("%s.%s", defaultalias, t),
			fmt.Sprintf("appsv1.%s", t),
			-1)
	}

	// ------------------------------
	// [ metav1 ]
	metav1Types := []string{
		"APIGroup",
		"ObjectMeta",
		"Time",
		"TypeMeta",
		"ManagedFieldsEntry",
		"OwnerReference",
		"CreateOptions",
		"DeleteOptions",
		"LabelSelector",
	}

	for _, t := range metav1Types {
		aliased = strings.Replace(aliased,
			fmt.Sprintf("%s.%s", defaultalias, t),
			fmt.Sprintf("metav1.%s", t),
			-1)
	}

	// ------------------------------
	// [ corev1 ]
	corev1types := []string{
		"Volume",
		"PodTemplateSpec",
		"PodSpec",
		"ResourceRequirements",
		"ResourceList",
		"VolumeDevice",
		"Probe",
		"Container",
		"EnvFromSource",
		"EnvVar",
		"VolumeMount",
		"Lifecycle",
		"SecurityContext",
		"EphemeralContainer",
		"LocalObjectReference",
		"Affinity",
		"Toleration",
		"HostAlias",
		"PodDNSConfig",
		"PodReadinessGate",
		"PreemptionPolicy",
		"TopologySpreadConstraint",
		"TerminationMessagePolicy",
		"PullPolicy",
		"RestartPolicy",
		"DNSPolicy",
		"PodSecurityContext",
	}
	// ------------------------------

	for _, t := range corev1types {
		aliased = strings.Replace(aliased,
			fmt.Sprintf("%s.%s", defaultalias, t),
			fmt.Sprintf("corev1.%s", t),
			-1)
	}

	return aliased
}

func sanitizeK8sObjectName(name string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9 \\-]+")
	return reg.ReplaceAllString(name, "")
}

func newl(n string) string {
	n = strings.ReplaceAll(n, "{", `{
`)
	return strings.ReplaceAll(n, ",", `,
`)

}
