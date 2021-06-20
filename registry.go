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

package yamyams

import (
	mydeployment "github.com/kris-nova/yamyams/apps/_deployment"
	myapplication "github.com/kris-nova/yamyams/apps/_example"
	yamyams "github.com/kris-nova/yamyams/pkg"
)

// Version is set at compile time and used for this specific version of YamYams
var Version string

// Load is where we can set up applications.
//
// This is called whenever the yamyams program starts.
func Load() {

	// We can keep them very simple, and hard code all the logic like this one.
	yamyams.Register(myapplication.New())

	// We can also have several instances of the same application like this.
	yamyams.Register(mydeployment.New("default", "example-1", "beeps", 3))
	yamyams.Register(mydeployment.New("default", "example-2", "boops", 1))
	yamyams.Register(mydeployment.New("default", "example-3", "cyber boops", 7))

}
