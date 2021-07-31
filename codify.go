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
	"io"
	"k8s.io/client-go/kubernetes/scheme"
	"text/template"
)

type MainGoValues struct {
	AuthorName string
	AuthorEmail string
	CopyrightYear string
	AppNameTitle string
	AppNameLower string
	Description string
	Version string
	Install string
	Uninstall string
}

func Codify(input io.Reader, v *MainGoValues) ([]byte, error) {

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

	// Setup template
	tpl := template.New("main")

	// Create the base file
	tpl, err = tpl.Parse(FormatMainGo)

	// Read the YAML
	serializer := scheme.Codecs.UniversalDeserializer()
	obj, _, err := serializer.Decode(ibytes, nil, nil)
	if err != nil {
		return code, fmt.Errorf("unable to deserialize: %v", err)
	}

	// Generate install source code
	iBytes, err := InstallBytes(obj)
	if err != nil {
		return code, fmt.Errorf("unable to generate Install(): %v", err)
	}
	v.Install = string(iBytes)

	// Generate uninstall source code
	uBytes, err := UninstallBytes(obj)
	if err != nil {
		return code, fmt.Errorf("unable to generate Uninstall(): %v", err)
	}
	v.Uninstall = string(uBytes)

	// Parse the system
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, v)
	if err != nil {
		return code, fmt.Errorf("unable to generate source code: %v", err)
	}

	// Return the buffer bytes :)
	return buf.Bytes(), nil
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

func UninstallBytes(object interface{}) ([]byte, error) {
	var code []byte
	return code, nil
}