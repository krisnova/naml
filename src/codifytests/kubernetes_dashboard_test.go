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

package codifytests

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/kris-nova/naml"
)

func TestKubernetesDashboard_v2_0_0(t *testing.T) {
	filename := "test_kubernetes_dashboard.yaml"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("unable to read file: %s: %v", filename, err)
	}
	buffer := bytes.Buffer{}
	buffer.Write(data)
	output, err := naml.Codify(&buffer, MainGoValues())
	if err != nil {
		t.Errorf("unable to codify: %s: %v", filename, err)
	}
	program, err := naml.Compile(output)
	if err != nil {
		t.Errorf("unable to compile: %s: %v", filename, err)
	}
	stdout, stderr, err := program.Execute([]string{""})
	if stderr.Len() > 0 {
		t.Errorf("failed executing binary: %s: %v", filename, err)
		t.Errorf(stderr.String())
	}
	if err != nil {
		t.Errorf("failed executing binary: %s: %v", filename, err)
	}
	t.Logf(stdout.String())
}
