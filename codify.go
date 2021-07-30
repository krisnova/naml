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
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
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
	return CodifyKubeObject(obj)
}

func CodifyKubeObject(object interface{}) ([]byte, error) {

	// ---
	// Left off here, we need to come up with a way to start
	// generating go code for us :)
	//
	// Have fun :)
	// ---

	var code []byte
	switch x := object.(type){
	case *v1.Pod:
		logger.Debug("pod %s", x.Name)
	default:

	}
	return code, nil
}
