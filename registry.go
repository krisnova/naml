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
	myapplication "github.com/kris-nova/naml/apps/_example"
	mydeployment "github.com/kris-nova/naml/apps/sampleapp"
	naml "github.com/kris-nova/naml/pkg"
)

// Version is set at compile time and used for this specific version of naml
var Version string

// Load is where we can set up applications.
//
// This is called whenever the naml program starts.
func Load() {

	// We can keep them very simple, and hard code all the logic like this one.
	naml.Register(myapplication.New())

	// We can also have several instances of the same application like this.
	naml.Register(mydeployment.New("default", "example-1", "beeps", 3))
	naml.Register(mydeployment.New("default", "example-2", "boops", 1))
	naml.Register(mydeployment.New("default", "example-3", "cyber boops", 7))

}
