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

package naml

import (
	"bytes"
	"strings"
	"testing"
)

// TestYAMLToGo will generate source code from YAML, and then compile the source, and then run the new program
func TestYAMLToGo(t *testing.T) {
	buf := bytes.Buffer{}
	buf.Write([]byte(TestNginxDeploymentYAML))
	v := &MainGoValues{
		AppNameLower:  "testnginx",
		AppNameTitle:  "TestNginx",
		AuthorName:    "Björn Nóva",
		AuthorEmail:   "barnaby@nivenly.com",
		CopyrightYear: "1999",
		Description:   "Test nginx deployment.",
	}
	code, err := Codify(&buf, v)
	if err != nil {
		t.Errorf("Unable to codify test YAML: %v", err)
	}
	codeStr := string(code)

	// Check year
	if !strings.Contains(codeStr, "1999") {
		t.Errorf("Failure interpolating values into codify.")
	}

	// Check Install
	if !strings.Contains(codeStr, "Install") {
		t.Errorf("Missing Install() in codify.")
	}

	// Check Uninstall
	if !strings.Contains(codeStr, "Uninstall") {
		t.Errorf("Missing Uninstall() in codify.")
	}

	// Check Objects
	if !strings.Contains(codeStr, "Objects") {
		t.Errorf("Missing Objects() in codify.")
	}

	// Compile the source code
	program, err := Compile(code)
	if err != nil {
		t.Errorf("Unable to compile source code: %v", err)
		t.FailNow()
	}

	// Execute the new program
	stdout, stderr, err := program.Execute([]string{""})
	if err != nil {
		t.Errorf("Unable to execute newly compiled program: %v", err)
	}
	t.Logf("Stdout Logs length: %d", len(stdout.Bytes()))
	t.Logf("Stderr Logs length: %d", len(stderr.Bytes()))

	// Remove the new program from the filesystem
	err = program.Remove()
	if err != nil {
		t.Errorf("Unable to remove newly compiled program: %v", err)
	}

}

const (
	TestNginxDeploymentYAML = `apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "3"
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"apps/appsv1","kind":"Deployment","metadata":{"annotations":{},"name":"nginx-deployment","namespace":"default"},"spec":{"replicas":4,"selector":{"matchLabels":{"app":"nginx"}},"template":{"metadata":{"labels":{"app":"nginx"}},"spec":{"containers":[{"image":"nginx:1.14.2","name":"nginx","ports":[{"containerPort":80}]}]}}}}
  creationTimestamp: null
  generation: 3
  name: nginx-deployment
  namespace: default
  resourceVersion: "254882"
  uid: 445061d9-5000-471b-8e06-45f5240dedb6
spec:
  replicas: 4
  selector:
    matchLabels:
      app: nginx
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx:1.14.2
        imagePullPolicy: IfNotPresent
        name: nginx
        ports:
        - containerPort: 80
          protocol: TCP
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
`
)
