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

package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/kris-nova/naml"
)

func TestManifests(t *testing.T) {
	files, err := ioutil.ReadDir("manifests")
	if err != nil {
		t.Errorf("unable to list manifest directory: %v", err)
	}
	wg := sync.WaitGroup{}
	for _, file := range files {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			t.Logf("testing [%s]", name)
			err := generateCompileRunYAML(filepath.Join("manifests", file.Name()))
			if err != nil {
				t.Errorf(err.Error())
			}
		}(file.Name())
	}
	wg.Wait()
	t.Logf("Manifest tests complete")
}

//
//func TestPrometheusNodeExportHelmChart(t *testing.T) {
//	err := generateCompileRunYAML("test_prometheus_helm_node_exporter.yaml")
//	if err != nil {
//		t.Errorf(err.Error())
//	}
//}
//
//func TestPrometheusHelmChart(t *testing.T) {
//	err := generateCompileRunYAML("test_prometheus_helm_chart.yaml")
//	if err != nil {
//		t.Errorf(err.Error())
//	}
//}
//
//func TestNginx(t *testing.T) {
//	err := generateCompileRunYAML("test_nginx.yaml")
//	if err != nil {
//		t.Errorf(err.Error())
//	}
//}
//
//func TestNginxIngressController(t *testing.T) {
//	err := generateCompileRunYAML("test_nginx_ingress_controller.yaml")
//	if err != nil {
//		t.Errorf(err.Error())
//	}
//}
//
//func TestKubernetesDashboard_v2_0_0(t *testing.T) {
//	err := generateCompileRunYAML("test_kubernetes_dashboard.yaml")
//	if err != nil {
//		t.Errorf(err.Error())
//	}
//}

func mainGoValues(name string) *naml.MainGoValues {
	return &naml.MainGoValues{
		AppNameLower:  strings.ToLower(name),
		AppNameTitle:  strings.ToUpper(name),
		AuthorName:    "Björn Nóva",
		AuthorEmail:   "barnaby@nivenly.com",
		CopyrightYear: "1999",
		Description:   "Test Description.",
	}
}

// generateCompileRunYAML will build a Go program from YAML and try to compile and run it :)
func generateCompileRunYAML(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("unable to read file: %s: %v", filename, err)
	}
	buffer := bytes.Buffer{}
	buffer.Write(data)
	output, err := naml.Codify(&buffer, mainGoValues(filename))
	if err != nil {
		return fmt.Errorf("unable to codify: %s: %v", filename, err)
	}
	program, err := naml.Compile(output)
	if program != nil {
		defer program.Remove()
	}
	if err != nil {
		return fmt.Errorf("unable to compile: %s: %v", filename, err)
	}
	_, stderr, err := program.Execute([]string{""})
	if stderr.Len() > 0 {
		return fmt.Errorf("failed executing binary: %s: %v: %s", filename, err, stderr.String())
	}
	if err != nil {
		return fmt.Errorf("failed executing binary: %s: %v", filename, err)
	}
	return nil
}
