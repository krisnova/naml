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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/kind/pkg/exec"

	"github.com/kris-nova/logger"
)

// [ NAML RPC Spec ]
//
// Handshake tcp advertisement over stdout.
// ----------------------------------------------------------------------------
// <tcpPort> # No greater than 6 bytes TCP port.
// # Empty newline to denote the end of the advertised port.
// ----------------------------------------------------------------------------
//
// RPC Servers will attempt to bind to an ephemeral
// TCP port at runtime. If a port is available from the
// the kernel the RPC will use it and advertise it as
// as new line in stdout which will be delimited by
// another newline.
//
const (
	// childListenAddr will be the address we use to listen on the child
	// process. By default :0 will tell the kernel to pick an ephemeral port.
	childListenAddr = ":0"

	// serverAuthToken is unimplemented for now. however we can use fields
	// like this later as needed.
	serverAuthToken = "valid"

	// localServerAddr will be the name of the local server we try to connect to.
	// In this example we use 127.0.0.1 in favor of localhost as we aren't sure DNS
	// is set to localhost on this machine.
	localServerAddr = "127.0.0.1"

	// contentType is the naml rpc contentType according to the HTTP spec.
	contentType = "application/json"
)

// remotes is a package level cache of known healthy remote naml RPC servers.
var remotes = make(map[string]*RPCPointer)

// RunRPC will run in RPC mode.
//
// TODO: Nóva to come add security auth and validation to our RPC
//
func RunRPC() error {
	logger.Debug("Running in rpc mode")

	// Listen on ephemeral port
	listen, err := net.Listen("tcp", childListenAddr)
	if err != nil {
		return fmt.Errorf("unable to listen on ephemeral port: %v", err)
	}

	// Find the port we are listening on.
	tcpPort := listen.Addr().(*net.TCPAddr).Port

	// The child will print the TCP port it is listening on
	// for the first part of the handshake with a parent.
	fmt.Println(tcpPort)

	// Register the remote procedures.
	http.HandleFunc("/list", list)
	http.HandleFunc("/install", install)
	http.HandleFunc("/uninstall", uninstall)
	http.HandleFunc("/handshake", handshake)

	// Listen for remote execution (as long as the parent is running)
	return http.Serve(listen, nil)
}

// list will return the list of applications with the child binary
// list returns []string
func list(w http.ResponseWriter, r *http.Request) {
	var remoteApps []*RPCApplication
	for appName, app := range Registry() {
		// --------------------------------------------------------------------
		// Here is where we list our applications.
		//
		// There is great concern here, as whatever fields are passed between
		// client and remote naml processes can be exposed.
		//
		// Do not pass secret material here.
		//
		// If in doubt. Don't.
		remote := &RPCApplication{
			AppName:        appName,
			AppDescription: app.Description(),
			AppVersion:     app.Meta().ResourceVersion,
		}
		remoteApps = append(remoteApps, remote)
		// --------------------------------------------------------------------

	}
	rjson, err := json.Marshal(&remoteApps)
	if err != nil {
		w.Write([]byte("unable to json marshal"))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(rjson)
	w.WriteHeader(http.StatusOK)
}

// RPCPointer is a pointer that allows us
// to connect to a remote RPC server
type RPCPointer struct {
	Message string
	Addr    string
}

// handshake will return meta information about the child
func handshake(w http.ResponseWriter, r *http.Request) {
	rjson, err := json.Marshal(&RPCPointer{
		Message: serverAuthToken},
	)
	if err != nil {
		w.Write([]byte("unable to json marshal"))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(rjson)
	w.WriteHeader(http.StatusOK)
}

// install will install an application with the child binary
func install(w http.ResponseWriter, r *http.Request) {
	remotejson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, fmt.Sprintf("unable to read body: %v", err))
		return
	}
	remoteApp := RPCApplication{}
	err = json.Unmarshal(remotejson, &remoteApp)
	if err != nil {
		writeError(w, fmt.Sprintf("unable to json unmarshal body: %v", err))
		return
	}

	app := Find(remoteApp.AppName)
	if app == nil {
		writeError(w, fmt.Sprintf("unable to find app: %s", remoteApp.AppName))
		return
	}

	// Now we handle auth
	client, err := Client()
	if err != nil {
		writeError(w, fmt.Sprintf("unable to generate client: %v", err))
		return
	}

	// Try to install the remote application
	err = app.Install(client)
	if err != nil {
		writeError(w, fmt.Sprintf("unable to install on remote: %v", err))
		return
	}

	// Success!
	w.WriteHeader(http.StatusOK)
}

// writeError will attempt to write a JSON error over the RPC
func writeError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	logger.Warning("unable to install remote application: %s", message)
	j := &RPCError{
		Message: message,
	}
	jjson, err := json.Marshal(&j)
	if err != nil {
		return // silently return!
	}
	w.Write(jjson)
}

// uninstall will install an application with the child binary
func uninstall(w http.ResponseWriter, r *http.Request) {
	remotejson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, fmt.Sprintf("unable to read body: %v", err))
		return
	}
	remoteApp := RPCApplication{}
	err = json.Unmarshal(remotejson, &remoteApp)
	if err != nil {
		writeError(w, fmt.Sprintf("unable to json unmarshal body: %v", err))
		return
	}

	app := Find(remoteApp.AppName)
	if app == nil {
		writeError(w, fmt.Sprintf("unable to find app: %s", remoteApp.AppName))
		return
	}

	// Now we handle auth
	client, err := Client()
	if err != nil {
		writeError(w, fmt.Sprintf("unable to generate client: %v", err))
		return
	}

	// Try to install the remote application
	err = app.Uninstall(client)
	if err != nil {
		writeError(w, fmt.Sprintf("unable to uninstall on remote: %v", err))
		return
	}

	// Success!
	w.WriteHeader(http.StatusOK)
}

// AddRPC will attempt to add a remote RPC server to the current runtime.
func AddRPC(path string) error {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("unable to find absolute path of remote: %v", err)
	}

	file, err := os.Stat(absolutePath)
	if err != nil {
		return fmt.Errorf("unable to execut file: %v", err)
	}

	// Todo use the file.Sys() interface to validate the file type

	logger.Info("Starting IPC with remote RPC server %s", file.Name())

	// Execute the child binary and pass in "c" for
	// shorthand to tell the child to run in child runtime mode.
	// logger.Info("Command %s c", path)
	cmd := exec.Command(absolutePath, "rpc")
	childOut := &bytes.Buffer{}
	childErr := &bytes.Buffer{}
	cmd.SetStdout(childOut)
	cmd.SetStderr(childErr)
	errCh := make(chan error)
	go func() {
		err = cmd.Run()
		if err != nil {
			errCh <- fmt.Errorf("unable to execute remote rpc process: %v", err)
		}
		errCh <- nil
	}()

	// ----------------------------------------------------------------------------
	// We are in RPC land, so we are handling cases as we see them.
	//
	// This is concurrent RPC code and should be considered very carefully.
	//
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
				// According to the naml RPC this should be the TCP port
				rpcHello := string(message)
				rpcSplit := strings.Split(rpcHello, "\n")
				if len(rpcSplit) < 1 {
					return fmt.Errorf("child failed validation, initial rpc message invalid use of newline")
				}
				// According to the spec the first line is our TCP port
				tcpPort := rpcSplit[0]

				// Check for > 6 bytes
				// [*]38179
				if childOut.Len() > 6 {
					return fmt.Errorf("child failed validation, initial rpc message greater than 5 bytes")
				}

				// Perform the naml RPC handshake
				remoteAddr := fmt.Sprintf("http://127.0.0.1:%s", tcpPort)
				remoteURL, err := url.Parse(fmt.Sprintf("%s/handshake", remoteAddr))
				if err != nil {
					return fmt.Errorf("unable parse url %s", remoteAddr)
				}
				response, err := http.Get(remoteURL.String())
				if err != nil {
					return fmt.Errorf("unable to authenticate with child [1]: %v", err)
				}
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					return fmt.Errorf("unable to authenticate with child [2]: %v", err)
				}

				rpcPointer := RPCPointer{}
				err = json.Unmarshal(body, &rpcPointer)
				if err != nil {
					return fmt.Errorf("unable to authenticate with child [3]: %v", err)
				}

				// We know we have a healthy rpc remote, register the healthy child connection string in memory
				rpcPointer.Addr = remoteAddr
				logger.Info("RPC Remote Addr: %s", rpcPointer.Addr)

				// Cache the child meta pointer
				remotes[path] = &rpcPointer

				return nil
			case childErr.Len() > 0:
				return fmt.Errorf("unable to execute child process with err: %s", string(childErr.Bytes()))
				break
			}
		}
	}
	return nil
}
