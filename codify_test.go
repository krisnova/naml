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
	"testing"
)

func TestYAMLDelimiterBottom(t *testing.T) {

	testString := `apiVersion: v1
kind: Service
metadata:
  labels:
    app: example
    bogus: ---
  name: example
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: example
  type: LoadBalancer
---

apiVersion: v1
kind: Service
metadata:
  labels:
    app: example
    bogus: ---
  name: example
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: example
  type: LoadBalancer
---
`

	buf := bytes.Buffer{}
	buf.Write([]byte(testString))
	objects, delta, err := ReaderToCodifyObjects(&buf)
	if delta != 0 {
		t.Errorf("Failure parsing YAML systems")
	}
	if err != nil {
		t.Errorf("YAML delimiter check: %v", err)
	}
	if len(objects) != 2 {
		t.Errorf("YAML delimiter split: %d", len(objects))
	}
}

func TestYAMLDelimiterNoSpace(t *testing.T) {

	testString := `
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: example
    bogus: ---
  name: example
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: example
  type: LoadBalancer
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: example
    bogus: ---
  name: example
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: example
  type: LoadBalancer
`

	buf := bytes.Buffer{}
	buf.Write([]byte(testString))
	objects, delta, err := ReaderToCodifyObjects(&buf)
	if delta != 0 {
		t.Errorf("Failure parsing YAML systems")
	}
	if err != nil {
		t.Errorf("YAML delimiter check: %v", err)
	}
	if len(objects) != 2 {
		t.Errorf("YAML delimiter split: %d", len(objects))
	}
}

func TestYAMLDelimiterTopLoad(t *testing.T) {

	testString := `
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: example
    bogus: ---
  name: example
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: example
  type: LoadBalancer

---

apiVersion: v1
kind: Service
metadata:
  labels:
    app: example
    bogus: ---
  name: example
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: example
  type: LoadBalancer
`

	buf := bytes.Buffer{}
	buf.Write([]byte(testString))
	objects, delta, err := ReaderToCodifyObjects(&buf)
	if delta != 0 {
		t.Errorf("Failure parsing YAML systems")
	}
	if err != nil {
		t.Errorf("YAML delimiter check: %v", err)
	}
	if len(objects) != 2 {
		t.Errorf("YAML delimiter split: %d", len(objects))
	}
}

func TestYAMLDelimiterHappy(t *testing.T) {

	testString := `apiVersion: v1
kind: Service
metadata:
  labels:
    app: example
    bogus: ---
  name: example
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: example
  type: LoadBalancer

---

apiVersion: v1
kind: Service
metadata:
  labels:
    app: example
    bogus: ---
  name: example
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: example
  type: LoadBalancer
`

	buf := bytes.Buffer{}
	buf.Write([]byte(testString))
	objects, delta, err := ReaderToCodifyObjects(&buf)
	if delta != 0 {
		t.Errorf("Failure parsing YAML systems")
	}
	if err != nil {
		t.Errorf("YAML delimiter check: %v", err)
	}
	if len(objects) != 2 {
		t.Errorf("YAML delimiter split: %d", len(objects))
	}
}

func TestYAMLDelimiterInline(t *testing.T) {

	testString := `apiVersion: v1
kind: Service
metadata:
  labels:
    app: example
    bogus: ---
  name: example
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: example
  type: LoadBalancer`

	buf := bytes.Buffer{}
	buf.Write([]byte(testString))
	objects, delta, err := ReaderToCodifyObjects(&buf)
	if delta != 0 {
		t.Errorf("Failure parsing YAML systems")
	}
	if err != nil {
		t.Errorf("inline YAML delimiter check: %v", err)
	}
	if len(objects) != 1 {
		t.Errorf("inline YAML delimiter split: %d", len(objects))
	}
}

func TestYAMLDelimiterDeltaMissing(t *testing.T) {

	testString := `apiVersion: v1
kind: Service
metadata:
  labels:
    app: example
    bogus: ---
  name: example
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: example
  type: LoadBalancer

---

apiVersion: unknownbad.io/v1alpha1
kind: UnknownBad
metadata:
    name: example-obc
spec:
    unknownField: example-unknown
    unknownFieldName: unknown.naml.io

`

	buf := bytes.Buffer{}
	buf.Write([]byte(testString))
	objects, delta, err := ReaderToCodifyObjects(&buf)
	if delta != 1 {
		t.Errorf("Expecting 1 failed object")
	}
	if err != nil {
		t.Errorf("inline YAML delimiter check: %v", err)
	}
	if len(objects) != 1 {
		t.Errorf("inline YAML delimiter split: %d", len(objects))
	}
}
