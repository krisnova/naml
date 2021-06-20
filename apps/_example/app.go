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
//    ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗
//    ████╗  ██║██╔═████╗██║   ██║██╔══██╗
//    ██╔██╗ ██║██║██╔██║██║   ██║███████║
//    ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║
//    ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║
//    ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝

package myapplication

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type MyApplication struct {
	metav1.ObjectMeta
}

func New() *MyApplication {
	return &MyApplication{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "my-application",
			ResourceVersion: "v1.0.0",
			Namespace:       "default",
			Labels: map[string]string{
				"k8s-app": "myapp",
				"app":     "myapp",
			},
			Annotations: map[string]string{
				"beeps": "boops",
			},
		},
	}
}

func (v *MyApplication) Install(client *kubernetes.Clientset) error {
	return fmt.Errorf("[install] for %s not yet implemented", v.Name)
}

func (v *MyApplication) Uninstall(client *kubernetes.Clientset) error {
	return fmt.Errorf("[uninstall] for %s not yet implemented", v.Name)
}

func (v *MyApplication) Meta() *metav1.ObjectMeta {
	return &v.ObjectMeta
}
