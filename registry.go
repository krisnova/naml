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
	"os"

	"github.com/kris-nova/logger"
)

var registry = make(map[string]Deployable)

// Register an application with naml
func Register(app Deployable) {

	// Validate the application
	if app == nil {
		logger.Critical("Unable to register NIL application.")
		os.Exit(1)
	}

	if app.Meta() == nil {
		logger.Critical("Unable to register NIL ObjectMeta for application")
		os.Exit(1)
	}

	if app.Meta().Name == "" {
		logger.Critical("Unable to register NIL ObjectMeta.Name for application")
		os.Exit(1)
	}

	registry[app.Meta().Name] = app
}

// Registry will return the registry
func Registry() map[string]Deployable {
	return registry
}

// Find an application by name
func Find(name string) Deployable {
	if app, ok := registry[name]; ok {
		return app
	}
	return nil
}
