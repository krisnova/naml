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

import v1 "k8s.io/api/core/v1"

type Pod struct {

}

func NewPod(pod *v1.Pod) *Pod {
	return &Pod{}
}

func (p Pod) Install() string {
	return `
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"": "",
			},
		},
	}
	_, err = client.CoreV1().Pods("").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}
`
}

func (p Pod) Uninstall() string {
	return `
	err = client.CoreV1().Pods("").Delete(context.TODO(), "", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
 `
}
