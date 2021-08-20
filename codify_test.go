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

// TestGoToYAML will generate source code from YAML, and then compile the source, and then run the new program
func TestGoToYAML(t *testing.T) {
	buf := bytes.Buffer{}
	buf.Write([]byte(TestYAML))
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

	// Compile the source code
	program, err := Compile(code)
	if err != nil {
		t.Errorf("Unable to compile source code: %v", err)
		t.FailNow()
	}

	// Output as YAML
	_, stderr, err := program.Execute([]string{"output"})
	if err != nil {
		t.Errorf("Unable to execute newly compiled program: %v", err)
	}
	if len(stderr.Bytes()) > 0 {
		t.Errorf("Unable to parse YAML: %v", err)
	}

	//// Pass the previous output to Codify()
	//newCode, err := Codify(stdout, v)
	//if err != nil {
	//	t.Errorf("Invalid YAML from output: %v", err)
	//}
	//
	//// Compile the new source code
	//_, err = Compile(newCode)
	//if err != nil {
	//	t.Errorf("Unable to compile new source code: %v", err)
	//	t.FailNow()
	//}

	// Output as YAML
	//stdout, stderr, err = newProgram.Execute([]string{"output"})
	//if err != nil {
	//	t.Errorf("Unable to execute newly compiled program: %v", err)
	//}
	//if len(stderr.Bytes()) > 0 {
	//	t.Errorf("Unable to parse YAML: %v", err)
	//}

	//	if len(newCode) != len(code) {
	//		t.Errorf("Code Length Original %d", len(code))
	//		t.Errorf("Code Length New      %d", len(newCode))
	//		t.Errorf("YAML conversion resulted in different programs. Non deterministic YAML!")
	//		t.FailNow()
	//		newSplit := strings.Split(string(newCode), `
	//`)
	//		oldSplit := strings.Split(string(code), `
	//`)
	//
	//		var highSplit, lowSplit []string
	//		if len(newSplit) > len(oldSplit) {
	//			highSplit = newSplit
	//			lowSplit = oldSplit
	//		} else {
	//			highSplit = oldSplit
	//			lowSplit = newSplit
	//		}
	//		for i := 0; i < len(highSplit); i++ {
	//			if len(lowSplit) <= i {
	//				t.Logf("+ ---> %s\n", highSplit[i])
	//				continue
	//			}
	//			if highSplit[i] != lowSplit[i] {
	//				t.Logf("\n")
	//				t.Logf("! [ORIGINAL]  ---> %s\n", highSplit[i])
	//				t.Logf("! [GENERATED] ---> %s\n", lowSplit[i])
	//				t.Logf("\n")
	//			} else {
	//				t.Log(highSplit[i])
	//			}
	//		}
	//	}

	// Remove the new program from the filesystem
	err = program.Remove()
	if err != nil {
		t.Errorf("Unable to remove newly compiled program: %v", err)
	}

}

// TestYAMLToGo will generate source code from YAML, and then compile the source, and then run the new program
func TestYAMLToGo(t *testing.T) {
	buf := bytes.Buffer{}
	buf.Write([]byte(TestYAML))
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
	TestYAML = `# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

apiVersion: v1
kind: Namespace
metadata:
  name: kubernetes-dashboard

---

apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kubernetes-dashboard

---

kind: Service
apiVersion: v1
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kubernetes-dashboard
spec:
  ports:
    - port: 443
      targetPort: 8443
  selector:
    k8s-app: kubernetes-dashboard

---

apiVersion: v1
kind: Secret
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard-certs
  namespace: kubernetes-dashboard
type: Opaque

---

apiVersion: v1
kind: Secret
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard-csrf
  namespace: kubernetes-dashboard
type: Opaque
data:
  csrf: ""

---

apiVersion: v1
kind: Secret
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard-key-holder
  namespace: kubernetes-dashboard
type: Opaque

---

kind: ConfigMap
apiVersion: v1
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard-settings
  namespace: kubernetes-dashboard

---

kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kubernetes-dashboard
rules:
  # Allow Dashboard to get, update and delete Dashboard exclusive secrets.
  - apiGroups: [""]
    resources: ["secrets"]
    resourceNames: ["kubernetes-dashboard-key-holder", "kubernetes-dashboard-certs", "kubernetes-dashboard-csrf"]
    verbs: ["get", "update", "delete"]
    # Allow Dashboard to get and update 'kubernetes-dashboard-settings' config map.
  - apiGroups: [""]
    resources: ["configmaps"]
    resourceNames: ["kubernetes-dashboard-settings"]
    verbs: ["get", "update"]
    # Allow Dashboard to get metrics.
  - apiGroups: [""]
    resources: ["services"]
    resourceNames: ["heapster", "dashboard-metrics-scraper"]
    verbs: ["proxy"]
  - apiGroups: [""]
    resources: ["services/proxy"]
    resourceNames: ["heapster", "http:heapster:", "https:heapster:", "dashboard-metrics-scraper", "http:dashboard-metrics-scraper"]
    verbs: ["get"]

---

kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
rules:
  # Allow Metrics Scraper to get metrics from the Metrics server
  - apiGroups: ["metrics.k8s.io"]
    resources: ["pods", "nodes"]
    verbs: ["get", "list", "watch"]

---

apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kubernetes-dashboard
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: kubernetes-dashboard
subjects:
  - kind: ServiceAccount
    name: kubernetes-dashboard
    namespace: kubernetes-dashboard

---

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kubernetes-dashboard
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kubernetes-dashboard
subjects:
  - kind: ServiceAccount
    name: kubernetes-dashboard
    namespace: kubernetes-dashboard

---

kind: Deployment
apiVersion: apps/v1
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kubernetes-dashboard
spec:
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      k8s-app: kubernetes-dashboard
  template:
    metadata:
      labels:
        k8s-app: kubernetes-dashboard
    spec:
      containers:
        - name: kubernetes-dashboard
          image: kubernetesui/dashboard:v2.3.1
          imagePullPolicy: Always
          ports:
            - containerPort: 8443
              protocol: TCP
          args:
            - --auto-generate-certificates
            - --namespace=kubernetes-dashboard
            # Uncomment the following line to manually specify Kubernetes API server Host
            # If not specified, Dashboard will attempt to auto discover the API server and connect
            # to it. Uncomment only if the default does not work.
            # - --apiserver-host=http://my-address:port
          volumeMounts:
            - name: kubernetes-dashboard-certs
              mountPath: /certs
              # Create on-disk volume to store exec logs
            - mountPath: /tmp
              name: tmp-volume
          livenessProbe:
            httpGet:
              scheme: HTTPS
              path: /
              port: 8443
            initialDelaySeconds: 30
            timeoutSeconds: 30
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            runAsUser: 1001
            runAsGroup: 2001
      volumes:
        - name: kubernetes-dashboard-certs
          secret:
            secretName: kubernetes-dashboard-certs
        - name: tmp-volume
          emptyDir: {}
      serviceAccountName: kubernetes-dashboard
      nodeSelector:
        "kubernetes.io/os": linux
      # Comment the following tolerations if Dashboard must not be deployed on master
      tolerations:
        - key: node-role.kubernetes.io/master
          effect: NoSchedule

---

kind: Service
apiVersion: v1
metadata:
  labels:
    k8s-app: dashboard-metrics-scraper
  name: dashboard-metrics-scraper
  namespace: kubernetes-dashboard
spec:
  ports:
    - port: 8000
      targetPort: 8000
  selector:
    k8s-app: dashboard-metrics-scraper

---

kind: Deployment
apiVersion: apps/v1
metadata:
  labels:
    k8s-app: dashboard-metrics-scraper
  name: dashboard-metrics-scraper
  namespace: kubernetes-dashboard
spec:
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      k8s-app: dashboard-metrics-scraper
  template:
    metadata:
      labels:
        k8s-app: dashboard-metrics-scraper
      annotations:
        seccomp.security.alpha.kubernetes.io/pod: 'runtime/default'
    spec:
      containers:
        - name: dashboard-metrics-scraper
          image: kubernetesui/metrics-scraper:v1.0.6
          ports:
            - containerPort: 8000
              protocol: TCP
          livenessProbe:
            httpGet:
              scheme: HTTP
              path: /
              port: 8000
            initialDelaySeconds: 30
            timeoutSeconds: 30
          volumeMounts:
            - mountPath: /tmp
              name: tmp-volume
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            runAsUser: 1001
            runAsGroup: 2001
      serviceAccountName: kubernetes-dashboard
      nodeSelector:
        "kubernetes.io/os": linux
      # Comment the following tolerations if Dashboard must not be deployed on master
      tolerations:
        - key: node-role.kubernetes.io/master
          effect: NoSchedule
      volumes:
        - name: tmp-volume
          emptyDir: {}
`
)
