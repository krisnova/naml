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
	"encoding/json"
	"fmt"

	"sigs.k8s.io/yaml"
)

const (
	OutputYAML OutputEncoding = 0
	OutputJSON OutputEncoding = 1
)

type OutputEncoding int

func RunOutput(appName string, o OutputEncoding) error {
	app := Find(appName)
	if app == nil {
		return fmt.Errorf("unable to find app: %s", appName)
	}

	// Install the application "nowhere" to register the components in memory
	app.Install(nil)

	switch o {

	// ---- [ JSON ] ----
	case OutputJSON:
		raw, err := json.MarshalIndent(app.Objects(), " ", "	")
		if err != nil {
			return fmt.Errorf("unable to JSON marshal: %v", err)
		}
		fmt.Println(string(raw))
		return nil

	// ---- [ YAML ] ----
	case OutputYAML:
		raw, err := yaml.Marshal(app.Objects())
		if err != nil {
			return fmt.Errorf("unable to YAML marshal: %v", err)
		}
		fmt.Println(string(raw))
		return nil

	// ---- [ DEFAULT ] ----
	default:
		raw, err := yaml.Marshal(app.Objects())
		if err != nil {
			return fmt.Errorf("unable to YAML marshal: %v", err)
		}
		fmt.Println(string(raw))
		return nil
	}
	return nil
}
