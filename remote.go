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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"

	"github.com/kris-nova/logger"
)

// RPCApplication is the pointer in the parent process
// to the child process internal application.
type RPCApplication struct {
	// Name is the name of the remote application
	AppName string

	// Remote is the associated remote RPC server
	Remote *RPCPointer

	// AppDescription is the application description
	AppDescription string

	// AppVersion is the application version
	AppVersion string
}

// RPCError is when something goes wrong over the RPC
type RPCError struct {
	Message string
}

// Install is the remote application install wrapper.
func (c *RPCApplication) Install(clientset *kubernetes.Clientset) error {
	if clientset != nil {
		return fmt.Errorf("*** security concern: clientset != nil ***")
	}
	ajson, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("unable to json marshal remote application: %v", err)
	}
	appBuffer := bytes.NewBuffer(ajson)
	response, err := http.Post(fmt.Sprintf("%s/install", c.Remote.Addr), contentType, appBuffer)
	if err != nil {
		return fmt.Errorf("unable to remote install(): %v", err)
	}
	if response.StatusCode != http.StatusOK {
		respjson, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("unable to remote install() status code: %d: %v", response.StatusCode, err)
		}
		rpcerr := RPCError{}
		err = json.Unmarshal(respjson, &rpcerr)
		if err != nil {
			return fmt.Errorf("unable to remote install() json error status code: %d: %v", response.StatusCode, err)
		}
		return fmt.Errorf("unable to remote install() status code: %d: %s", response.StatusCode, rpcerr.Message)
	}
	logger.Info("Success!")
	return nil
}

func (c *RPCApplication) Uninstall(clientset *kubernetes.Clientset) error {
	if clientset != nil {
		return fmt.Errorf("*** security concern: clientset != nil ***")
	}
	ajson, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("unable to json marshal remote application: %v", err)
	}
	appBuffer := bytes.NewBuffer(ajson)
	response, err := http.Post(fmt.Sprintf("%s/uninstall", c.Remote.Addr), contentType, appBuffer)
	if err != nil {
		return fmt.Errorf("unable to remote install(): %v", err)
	}
	if response.StatusCode != http.StatusOK {
		respjson, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("unable to remote install() status code: %d: %v", response.StatusCode, err)
		}
		rpcerr := RPCError{}
		err = json.Unmarshal(respjson, &rpcerr)
		if err != nil {
			return fmt.Errorf("unable to remote install() json error status code: %d: %v", response.StatusCode, err)
		}
		return fmt.Errorf("unable to remote install() status code: %d: %s", response.StatusCode, rpcerr.Message)
	}
	logger.Info("Success!")
	return nil
}

func (c *RPCApplication) Description() string {
	return c.AppDescription
}

func (c *RPCApplication) Meta() *v1.ObjectMeta {
	return &v1.ObjectMeta{
		Name:            c.AppName,
		ResourceVersion: c.AppVersion,
	}
}

// RegisterRemoteApplications will call list() on all remote RPC servers
// and register the applications as pointers on the remote.
func RegisterRemoteApplications() error {
	for path, remote := range remotes {
		logger.Info("Calling list() on remote RPC %s from %s", remote.Addr, path)
		var remoteApps []*RPCApplication
		resp, err := http.Get(fmt.Sprintf("%s/list", remote.Addr))
		if err != nil {
			return fmt.Errorf("unable to list() remote [1] %s: %v", path, err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("unable to list() remote [2] %s: %v", path, err)
		}
		err = json.Unmarshal(body, &remoteApps)
		if err != nil {
			return fmt.Errorf("unable to list() remote [3] %s: %v", path, err)
		}
		for _, app := range remoteApps {
			app.Remote = remote // Don't forget to set the remote!
			logger.Info("Registering remote app %s", app.AppName)

			//
			// ***************
			Register(app) // *
			// ***************
			//

		}
	}
	return nil
}
