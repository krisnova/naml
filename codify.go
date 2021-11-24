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
	"go/format"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/template"

	policyv1 "k8s.io/api/policy/v1beta1"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/fatih/color"

	"github.com/kris-nova/logger"
	"github.com/kris-nova/naml/codify"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

const (

	// codifyGoFormat will toggle if Codify() will also "go fmt" the generated code.
	codifyGoFormat bool = true

	// YAMLDelimiter is the official delimiter used to append multiple
	// YAML files together into the same file.
	//
	//	Reference: https://yaml.org/spec/1.2/spec.html
	//
	YAMLDelimiter string = "\n---\n"
)

// We ARE in fact doing a lot of string handling here
// So we use strings as often as possible.

// CodifyValues are ultimately what is rendered
// into the .naml files in /src. These values
// are what will be created in the output.
type CodifyValues struct {
	LibraryMode   bool
	AuthorName    string
	AuthorEmail   string
	CopyrightYear string
	AppNameTitle  string
	AppNameLower  string
	Description   string
	Version       string
	Install       string
	Uninstall     string
	Packages      string
	PackageName   string
}

type CodifyObject interface {

	// Install returns the snippet of code that would
	// traditionally live inside a function. This
	// will define literally (what it can) a struct
	// for the object, and pass it to the corresponding
	// kubernetes library.
	Install() (string, []string)

	// Uninstall is the reverse library call of install.
	Uninstall() string
}

// Codify will take any valid Kubernetes YAML as an io.Reader
// and do it's best to return a syntactically correct Go program
// that is NAML compliant.
//
// The NAML codebase is Apache 2.0 licensed, so we assume that
// any calling code will adopt the same Apache license.
func Codify(input io.Reader, v *CodifyValues) ([]byte, error) {

	if v.PackageName == "" {
		return []byte(""), fmt.Errorf("missing packageName")
	}

	var code []byte

	// Setup template with a unique name based on the input
	tpl := template.New(fmt.Sprintf("%+v%+v", input, v.AppNameLower))

	// Create the base file
	templateString := FormatMainGo
	if v.LibraryMode {
		templateString = FormatLibraryGo
	}
	tpl, err := tpl.Parse(templateString)
	if err != nil {
		return code, fmt.Errorf("unable to create main go tempalte: %v", err)
	}

	// Find the objects
	objs, delta, err := ReaderToCodifyObjects(input)
	if err != nil {
		return code, fmt.Errorf("unable to parse objects: %v", err)
	}

	// Create map of used packages
	packages := make(map[string]bool)

	// Append both install and uninstall for every object
	for _, obj := range objs {
		// get the install code and packages it depends on
		install, localPackages := obj.Install()

		// add all packages to the package map
		for _, pkg := range localPackages {
			packages[pkg] = true
		}

		// add to install
		v.Install = fmt.Sprintf("%s\n%s", v.Install, install)

		if v.Uninstall == "" {
			v.Uninstall = obj.Uninstall()
		} else {
			v.Uninstall = fmt.Sprintf("%s\n%s", v.Uninstall, obj.Uninstall())
		}
	}

	// Gather list of packages and sort them
	packagesSlice := make([]string, 0)
	for k, _ := range packages {
		packagesSlice = append(packagesSlice, k)
	}
	sort.Strings(packagesSlice)

	// define list of import aliases
	packageAliases := codify.KubernetesImportPackageMap

	packagesCode := ""
	for _, pkg := range packagesSlice {
		packageAlias := ""
		if name, ok := packageAliases[pkg]; ok == true {
			packageAlias = name + " "
		}
		packagesCode += fmt.Sprintf("\t%s\"%s\"\n", packageAlias, pkg)
	}

	v.Packages = packagesCode

	// Parse the system
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, v)
	if err != nil {
		return code, fmt.Errorf("unable to generate source code: %v", err)
	}

	// Grab the source code in []byte form
	src := buf.Bytes()

	// Go fmt!
	var fmtBytes []byte
	if codifyGoFormat {
		fmtBytes, err = format.Source(src)
	} else {
		fmtBytes = src
	}
	if err != nil {
		// Unable to auto format the code so let's debug!
		lines := strings.Split(string(src), "\n")
		lns := strings.Split(err.Error(), ":")
		line, err := strconv.Atoi(lns[0])
		if err != nil || len(lines) < line {
			return code, fmt.Errorf("unable to auto format code: %v", err)
		}
		fmt.Printf("------------------------------------------\n")
		fmt.Printf("Text template compile error in main.go.tpl\n")
		fmt.Printf("------------------------------------------\n")
		fmt.Printf("\n")
		fmt.Printf("%s\n", color.WhiteString(lines[line-1]))
		fmt.Printf("%s   %s\n", color.RedString(lines[line]), color.YellowString("<---"))
		fmt.Printf("%s\n", color.WhiteString(lines[line+1]))
		fmt.Printf("\n")
		fmt.Printf("------------------------------------------\n")
		fmt.Printf("\n")
		return code, fmt.Errorf("unable to auto format code: %v", err)
	}

	if delta > 0 {
		return fmtBytes, fmt.Errorf("incomplete: unable to Codify (%d) YAML system(s)", delta)
	}
	return fmtBytes, nil
}

// ReaderToBytes is basically a wrapper for ReadAll, however
// we add in some specific error language for stdin.
func ReaderToBytes(input io.Reader) ([]byte, error) {
	var code []byte
	ibytes, err := io.ReadAll(input)
	if err != nil {
		return code, fmt.Errorf("unable to read all of stdin: %v", err)
	}
	logger.Debug("Read %d bytes from stdin", len(ibytes))
	return ibytes, nil
}

// ReaderToCodifyObjects will convert an io.Reader to naml compatible Go objects.
//
// This function works by doing the best it can, and will return as many CodifyObjects
// as possible.
// The function will return a positive integer for every YAML object it detects, that
// it is unable to Codify.
// If the delta is greater than 0, that means we have encountered a loss.
func ReaderToCodifyObjects(input io.Reader) ([]CodifyObject, int, error) {
	var objects []CodifyObject
	ibytes, err := ReaderToBytes(input)
	if err != nil {
		return objects, -1, err
	}
	clean := string(cleanRaw(ibytes))
	yamls := strings.Split(clean, YAMLDelimiter)
	// We support more than one "YAML" per the delimiter
	// So we need to deal in sets.
	d := len(yamls)
	for _, yaml := range yamls {
		raw := []byte(yaml)
		cObjects, err := toCodify(raw)
		if err != nil {
			return objects, -1, fmt.Errorf("unable to codify: %v", err)
		}
		// Merge the items
		for _, c := range cObjects {
			objects = append(objects, c)
		}
	}
	c := len(objects)
	return objects, d - c, nil
}

// cleanRaw will clean raw yaml
// it will remove empty lines, empty lines with whitespace, and comments
func cleanRaw(raw []byte) []byte {
	rawString := string(raw)
	lines := strings.Split(rawString, "\n")
	var cleanLines []string
	for _, line := range lines {
		// Check for comments
		if strings.HasPrefix(line, "#") {
			continue
		}
		// Ignore empty lines
		if line == "" {
			continue
		}

		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		cleanLines = append(cleanLines, line)
	}
	cleanedRawString := strings.Join(cleanLines, "\n")
	return []byte(cleanedRawString)
}

func toCodify(raw []byte) ([]CodifyObject, error) {
	var objects []CodifyObject
	if len(raw) <= 1 {
		return objects, nil
	}

	serializer := scheme.Codecs.UniversalDeserializer()
	var decoded runtime.Object
	decoded, _, err := serializer.Decode([]byte(raw), nil, nil)
	if err != nil {
		// Here we try CRDs
		decoded, _, err = serializer.Decode([]byte(raw), nil, &apiextensionsv1.CustomResourceDefinition{})
		if err != nil {
			return nil, fmt.Errorf("trying CRD: unable to deserialize in codify: %v\n\nraw:\n\n%s", err, string(raw))
		}
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
	case *batchv1.CronJob:
		objects = append(objects, codify.NewCronJob(x))
	case *rbacv1.Role:
		objects = append(objects, codify.NewRole(x))
	case *rbacv1.ClusterRole:
		objects = append(objects, codify.NewClusterRole(x))
	case *rbacv1.RoleBinding:
		objects = append(objects, codify.NewRoleBinding(x))
	case *rbacv1.ClusterRoleBinding:
		objects = append(objects, codify.NewClusterRoleBinding(x))
	case *corev1.ServiceAccount:
		objects = append(objects, codify.NewServiceAccount(x))
	case *corev1.Secret:
		objects = append(objects, codify.NewSecret(x))
	case *networkingv1.IngressClass:
		objects = append(objects, codify.NewIngressClass(x))
	case *networkingv1.Ingress:
		objects = append(objects, codify.NewIngress(x))
	case *policyv1.PodSecurityPolicy:
		objects = append(objects, codify.NewPodSecurityPolicy(x))
	case *admissionregistrationv1.ValidatingWebhookConfiguration:
		objects = append(objects, codify.NewValidatingwebhookConfiguration(x))
	case *policyv1.PodDisruptionBudget:
		objects = append(objects, codify.NewPodDisruptionBudget(x))
	// CRDs is going to take some special care...
	case *apiextensionsv1.CustomResourceDefinition:
	// ignore CRDs for now!
	//objects = append(objects, codify.NewCustomResourceDefinition(x))
	case *corev1.Namespace:
		objects = append(objects, codify.NewNamespace(x))
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
