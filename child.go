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
//    ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗
//    ████╗  ██║██╔═████╗██║   ██║██╔══██╗
//    ██╔██╗ ██║██║██╔██║██║   ██║███████║
//    ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║
//    ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║
//    ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝

package naml

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"

	"github.com/kris-nova/logger"
)

// ChildMeta is the meta data for all requests
type ChildMeta struct {
	Addr string
}

type ChildApp struct {
	Name string
}

func (c *ChildApp) Install(clientset *kubernetes.Clientset) error {
	return nil
}

func (c *ChildApp) Uninstall(clientset *kubernetes.Clientset) error {
	return nil
}

func (c *ChildApp) Meta() *v1.ObjectMeta {
	return &v1.ObjectMeta{
		Name: c.Name,
	}
}

// children are our child RPC namls
var children = make(map[string]*ChildMeta)

func RegisterChildren() error {

	for path, child := range children {
		logger.Info("Calling list() on child %s from child %s", child.Addr, path)
		childRegistry := []string{""}
		resp, err := http.Get(fmt.Sprintf("%s/list", child.Addr))
		if err != nil {
			return fmt.Errorf("unable to list() child [1] %s: %v", path, err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("unable to list() child [2] %s: %v", path, err)
		}
		err = json.Unmarshal(body, &childRegistry)
		if err != nil {
			return fmt.Errorf("unable to list() child [3] %s: %v", path, err)
		}
		logger.Info("Registering children...")
		for _, appName := range childRegistry {
			logger.Info("Registering child app %s", appName)

			// Here is where we *FINALLY* register the child app
			childApp := &ChildApp{
				Name: appName,
			}

			//
			// ***************
			Register(childApp) // *
			// ***************
			//

		}
	}

	return nil
}
