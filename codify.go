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
	"fmt"
	"github.com/kris-nova/logger"
	"github.com/kris-nova/naml/codify"
	"go/format"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"strings"
	"text/template"
)

// YAMLDelimiter is the official delimiter used to append multiple
// YAML files together into the same file.
//
//	Reference: https://yaml.org/spec/1.2/spec.html
//
// Furthermore let it be documented that at the 2018 KubeCon pub trivia
// Bryan Liles (https://twitter.com/bryanl) correctly had answered the
// trivia question with the correct delimiter of 3 characters "---" and
// was awarded no points for his correct answer, while an opposing team
// was awarded a single point for their incorrect answer of 2 characters "--".
//
// If the correct delimiter points would have been awarded to Brian's team
// they would technically should have been crowned KubeCon pub champions of 2018.
const YAMLDelimiter string = "---"

// We ARE in fact doing a lot of string handling here
// So we use strings as often as possible.

// MainGoValues are ultimately what is rendered
// into the .naml files in /src. These values
// are what will be created in the output.
type MainGoValues struct {
	AuthorName    string
	AuthorEmail   string
	CopyrightYear string
	AppNameTitle  string
	AppNameLower  string
	Description   string
	Version       string
	Install       string
	Uninstall     string
}

type CodifyObject interface {
	// Install returns the snippet of code that would
	// traditionall live INSIDE of a function. This
	// will define literally (what it can) a struct
	// for the object, and pass it to the corresponding
	// kubernetes library.
	Install() string

	// Uninstall is the reverse library call of install
	Uninstall() string
}

// Codify will take any valid Kubernetes YAML as an io.Reader
// and do it's best to return a syntactically correct Go program
// that is NAML compliant.
//
// The NAML codebase is Apache 2.0 licensed, so we assume that
// any calling code will adopt the same Apache license.
func Codify(input io.Reader, v *MainGoValues) ([]byte, error) {
	// Right away we assume this is a finite stream that can
	// fit into memory. Therefore there is no need for a scanner.
	var code []byte
	ibytes, err := io.ReadAll(input)
	if err != nil {
		return code, fmt.Errorf("unable to read all of stdin: %v", err)
	}
	logger.Debug("Read %d bytes from stdin", len(ibytes))

	// Setup template
	tpl := template.New("main")

	// Create the base file
	tpl, err = tpl.Parse(FormatMainGo)

	// Find the objects
	objs, err := YAMLToCodifyObjects(ibytes)
	if err != nil {
		return code, fmt.Errorf("unable to parse objects: %v", err)
	}
	logger.Debug("Found %d objects to parse", len(objs))

	// Append both install and uninstall for every object
	for _, obj := range objs {
		if v.Install == "" {
			v.Install = obj.Install()
		} else {
			v.Install = fmt.Sprintf("%s\n%s", v.Install, obj.Install())
		}
		if v.Uninstall == "" {
			v.Uninstall = obj.Uninstall()
		} else {
			v.Uninstall = fmt.Sprintf("%s\n%s", v.Uninstall, obj.Uninstall())
		}
	}

	// Parse the system
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, v)
	if err != nil {
		return code, fmt.Errorf("unable to generate source code: %v", err)
	}

	// Grab the source code in []byte form
	src := buf.Bytes()

	// Go fmt!
	fmtBytes, err := format.Source(src)
	if err != nil {
		return code, fmt.Errorf("unable to auto format code: %v", err)
	}
	return fmtBytes, nil
}

func YAMLToCodifyObjects(raw []byte) ([]CodifyObject, error) {
	rawStr := string(raw)
	var objects []CodifyObject
	yamls := strings.Split(rawStr, YAMLDelimiter)
	// We support more than one "YAML" per the delimiter
	// So we need to deal in sets.
	for _, yaml := range yamls {
		raw := []byte(yaml)
		cObjects, err := toCodify(raw)
		if err != nil {
			return objects, fmt.Errorf("unable to codify: %v", err)
		}
		// Merge the items
		for _, c := range cObjects {
			objects = append(objects, c)
		}
	}
	return objects, nil
}

func toCodify(raw []byte) ([]CodifyObject, error) {
	var objects []CodifyObject

	serializer := scheme.Codecs.UniversalDeserializer()
	decoded, _, err := serializer.Decode([]byte(raw), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to deserialize: %v", err)
	}

	// -------------------------------------------------------------------
	// [NAML Type Switch]
	//
	// Because of the lack of generics in Go, we are having a lot of fun
	// doing things like this.
	//
	// Anyway if you are interested in adding a NAML type it MUST be switched
	// on here.
	//
	switch x := decoded.(type) {
	case *corev1.List:
		// Lists are recursive items
		// But we error each time and just
		// base the error from the inner system.
		for _, item := range x.Items {
			cObjects, err := toCodify(item.Raw)
			if err != nil {
				return objects, err
			}
			// Merge the items
			for _, c := range cObjects {
				objects = append(objects, c)
			}
		}
	case *corev1.Pod:
		objects = append(objects, codify.NewPod(x))
	case *appsv1.Deployment:
		objects = append(objects, codify.NewDeployment(x))
	case *appsv1.StatefulSet:
		objects = append(objects, codify.NewStatefulSet(x))
	case *appsv1.DaemonSet:
		objects = append(objects, codify.NewDaemonSet(x))
	case *corev1.ConfigMap:
		objects = append(objects, codify.NewConfigMap(x))
	case *corev1.Service:
		objects = append(objects, codify.NewService(x))
	case *corev1.PersistentVolume:
		objects = append(objects, codify.NewPersistentVolume(x))
	case *corev1.PersistentVolumeClaim:
		objects = append(objects, codify.NewPersistentVolumeClaim(x))
	case *batchv1.Job:
		objects = append(objects, codify.NewJob(x))
	case *appsv1.ReplicaSet:
	case *corev1.Endpoints:
		// Ignore ReplicaSet, Endpoints
		break
	default:
		return nil, fmt.Errorf("missing NAML support for type: %s", x.GetObjectKind().GroupVersionKind().Kind)
	}
	// -------------------------------------------------------------------
	return objects, nil
}
