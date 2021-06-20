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

package apps

import (
	"github.com/kris-nova/yamyams/apps/sampleapp"
	yamyams "github.com/kris-nova/yamyams/pkg"
	"testing"
)

// TestSampleApp is an example integration test that can be used to
// install and uninstall a sample application in Kubernetes.
func TestSampleApp(t *testing.T) {
	client, err := yamyams.ClientFromPath(yamyams.TestClusterKubeConfigPath())
	if err != nil {
		t.Errorf("unable to create client: %v", err)
	}
	app := sampleapp.New("default", "sample-app", "beeps-boops", 2)
	err = app.Install(client)
	if err != nil {
		t.Errorf("unable to install sample-app: %v", err)
	}
	err = app.Uninstall(client)
	if err != nil {
		t.Errorf("unable to uninstall sample-app: %v", err)
	}
}

// TestSampleAppName shows how you can test arbitrary parts of your application.
func TestSampleAppName(t *testing.T) {
	app := sampleapp.New("default", "sample-app", "beeps-boops", 2)
	if app.Name != "sample-app" {
		t.Errorf(".Name is not plumbed through from New()")
	}
}
