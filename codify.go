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
	"fmt"
	"github.com/kris-nova/logger"
	"io"
	"k8s.io/client-go/kubernetes/scheme"
	"strings"
	"time"
)

func Codify(input io.Reader) ([]byte, error) {

	// For the first few versions let's just
	// read "all" of stdin
	//
	// Todo: use scanner to read by \n
	var code []byte
	ibytes, err := io.ReadAll(input)
	if err != nil {
		return code, fmt.Errorf("unable to read all of stdin: %v", err)
	}
	logger.Debug("Read %d bytes from stdin", len(ibytes))

	serializer := scheme.Codecs.UniversalDeserializer()
	obj, _, err := serializer.Decode(ibytes, nil, nil)
	if err != nil {
		return code, fmt.Errorf("unable to deserialize: %v", err)
	}
	iBytes, err := InstallBytes(obj)
	if err != nil {
		return code, fmt.Errorf("unable to generate Install(): %v", err)
	}
	fileString := MainGo(time.Now().Year(), "NamlApp", "Kris Nóva", "<kris@nivenly.com>", iBytes)
	return []byte(fileString), nil
}

func InstallBytes(object interface{}) ([]byte, error) {
	var code []byte
	// ---
	// Left off here, we need to come up with a way to start
	// generating go code for us :)
	//
	// Have fun :)
	// ---

	//switch x := object.(type){
	//case *v1.Pod:
	//	logger.Debug("pod %s", x.Name)
	//default:
	//
	//}
	return code, nil
}

func DefaultMainGo() string {
	year := time.Now().Year()
	return MainGo(year, "NamlApp","Kris Nóva", "<kris@nivenly.com>", []byte("// app goes here"))
}

func MainGo(copyrightYear int, appName, authorName, authorEmail string, install []byte) string {
	lowerName := strings.ToLower(appName)
	titleName := strings.ToTitle(lowerName)
	return fmt.Sprintf(FormatMain, copyrightYear, authorName, authorEmail,
		titleName, lowerName,
		titleName, titleName, string(install),
		titleName, titleName, titleName)
}


const FormatMain = `
// Copyright © %d %s %s
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

package main

import (
	"context"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kris-nova/naml"
	"k8s.io/client-go/kubernetes"
)

func main() {

	naml.Register(&%s{"%s"})

	err := naml.RunCommandLine()
	if err != nil {
		os.Exit(1)
	}
}

type %s struct {
	metav1.ObjectMeta
	Name string
}

func (n *%s) Install(client *kubernetes.Clientset) error {
	%s 
	return err
}

func (n *%s) Uninstall(client *kubernetes.Clientset) error {
	return nil
}

func (n *%s) Description() string {
	return n.Name
}

func (n *%s) Meta() *metav1.ObjectMeta {
	return &n.ObjectMeta
}
`