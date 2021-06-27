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
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"sigs.k8s.io/kind/pkg/exec"

	"github.com/kris-nova/logger"
)

const (
	// childListenAddr will be the address we use to listen on the child
	// process. By default :0 will tell the kernel to pick a port.
	childListenAddr = ":0"
)

// RuntimeChild will run naml in a very unusual mode.
//
// We define our own RPC for the program.
//
// [ naml rpc spec]
//
// The spec is that the child will attempt to listen on an ephemeral port,
// and then the first line of the RPC is the TCP port that the child
// is listening on which is passed over stdout to a parent.
//
// Next the RPC states that all methods are exposed via JSON
//
// TODO: Nóva to come add security auth and validation to our RPC
//
func RuntimeChild() error {
	logger.Debug("Running in child runtime mode")

	// Listen on ephemeral port
	listen, err := net.Listen("tcp", childListenAddr)
	if err != nil {
		return fmt.Errorf("unable to listen on ephemeral port: %v", err)
	}

	// Print ephemeral port for server ipc as per the "spec"
	tcpPort := listen.Addr().(*net.TCPAddr).Port

	// The child will print the TCP port it is listening on
	// for the first part of the handshake with a parent.
	fmt.Println(tcpPort)

	// Register the list() procedure
	http.HandleFunc("/list", list)
	http.HandleFunc("/install", install)
	http.HandleFunc("/uninstall", uninstall)
	http.HandleFunc("/meta", meta)

	// TODO consider using custom protocol in favor of HTTP
	return http.Serve(listen, nil)
}

// list will return the list of applications with the child binary
// list returns []string
func list(w http.ResponseWriter, r *http.Request) {
	var childApps []string
	for appName := range Registry() {
		childApps = append(childApps, appName)
	}

	rjson, err := json.Marshal(&childApps)
	if err != nil {
		w.Write([]byte("unable to json marshal"))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(rjson)
	w.WriteHeader(http.StatusOK)
}

// meta will return meta information about the child
func meta(w http.ResponseWriter, r *http.Request) {
	// Todo &ChildMeta can hold our auth material
	rjson, err := json.Marshal(&ChildMeta{})
	if err != nil {
		w.Write([]byte("unable to json marshal"))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(rjson)
	w.WriteHeader(http.StatusOK)
}

// install will install an application with the child binary
func install(w http.ResponseWriter, r *http.Request) {

}

// uninstall will install an application with the child binary
func uninstall(w http.ResponseWriter, r *http.Request) {

}

// AddChild will attempt to validate and add
// a runtime child NAML to the current runtime.
func AddChild(path string) error {
	// Execute the child and find the TCP port
	file, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("unable to execut file: %v", err)
	}

	// Todo use the file.Sys() interface to validate the file type

	logger.Info("Starting IPC with child %s", file.Name())

	// Execute the child binary and pass in "c" for
	// shorthand to tell the child to run in child runtime mode.
	// logger.Info("Command %s c", path)
	cmd := exec.Command(path, "c")
	childOut := &bytes.Buffer{}
	childErr := &bytes.Buffer{}
	cmd.SetStdout(childOut)
	cmd.SetStderr(childErr)

	errCh := make(chan error)
	go func() {
		err = cmd.Run()
		if err != nil {
			errCh <- fmt.Errorf("unable to execute child process: %v", err)
		}
		errCh <- nil
	}()

	// We are in RPC land, so we are handling cases
	// as we see them.
	for {
		select {
		// Our one and only case is if we have an error over the channel
		case runErr := <-errCh:
			// The cmd.Run() has returned an error, something is very wrong.
			if runErr != nil {
				return fmt.Errorf("unable to execute child process with err: %v", runErr)
			}
			logger.Info("Child has exited.")
			break
		default:
			switch {
			case childOut.Len() > 0:
				message := childOut.Bytes()
				// According to the naml RPC
				// this should be the TCP port
				rpcHello := string(message)
				rpcSplit := strings.Split(rpcHello, "\n")
				if len(rpcSplit) != 2 {
					return fmt.Errorf("child failed validation, initial rpc message invalid use of newline")
				}
				// According to the spec the first line is our TCP port
				tcpPort := rpcSplit[0]

				// Check for > 6 bytes
				// [*]38179
				if childOut.Len() > 6 {
					return fmt.Errorf("child failed validation, initial rpc message greater than 5 bytes")
				}

				// Say hello
				metaAddr := fmt.Sprintf("http://127.0.0.1:%s", tcpPort)
				u, err := url.Parse(fmt.Sprintf("%s/meta", metaAddr))
				if err != nil {
					return fmt.Errorf("unable parse url %s", metaAddr)
				}
				logger.Info("Child meta Get(): %s", u.String())
				response, err := http.Get(u.String())
				if err != nil {
					return fmt.Errorf("unable to authenticate with child [1]: %v", err)
				}
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					return fmt.Errorf("unable to authenticate with child [2]: %v", err)
				}

				childMeta := ChildMeta{}
				err = json.Unmarshal(body, &childMeta)
				if err != nil {
					return fmt.Errorf("unable to authenticate with child [3]: %v", err)
				}

				// We know we have a healthy child, register the healthy child connection string in memory
				childMeta.Addr = metaAddr
				logger.Info("Child Meta Addr: %s", childMeta.Addr)

				// Cache the child meta pointer
				children[path] = &childMeta

				return nil
			case childErr.Len() > 0:
				return fmt.Errorf("unable to execute child process with err: %s", string(childErr.Bytes()))
				break
			}
		}
	}
	return nil
}
