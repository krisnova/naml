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

package main

import (
	"fmt"
	yamyams "github.com/kris-nova/yamyams/pkg"
)

func List() {
	fmt.Println("$ yamyams install    [app]")
	fmt.Println("$ yamyams uninstall  [app]")
	fmt.Println("")
	for _, app := range yamyams.Registry() {
		fmt.Printf("[%s]\n", app.Meta().Name)
		fmt.Printf("\tnamespace  : %s\n", app.Meta().Namespace)
		fmt.Printf("\tversion    : %s\n", app.Meta().ResourceVersion)
		if description, ok := app.Meta().Labels["description"]; ok {
			fmt.Printf("\tdescription : %s\n", description)
		}
		fmt.Printf("\n")
	}
	fmt.Println("")
}
