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
	"os"

	"github.com/kris-nova/logger"
)

var registry = make(map[string]Deployable)

// RegisterAndExit will register the app or exit with an error message in stdout
func RegisterAndExit(app Deployable) {
	err := RegisterAndError(app)
	if err != nil {
		logger.Critical("%v", err)
		os.Exit(1)
	}

}

// Register an application with naml
func Register(app Deployable) {
	RegisterAndExit(app)
}

func RegisterAndError(app Deployable) error {

	// Validate the application
	if app == nil {
		return fmt.Errorf("Unable to register NIL application.")

	}

	if app.Meta() == nil {
		return fmt.Errorf("Unable to register NIL ObjectMeta for application")

	}

	if app.Meta().Name == "" {
		return fmt.Errorf("Unable to register NIL ObjectMeta.Name for application")

	}

	registry[app.Meta().Name] = app
	return nil
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
