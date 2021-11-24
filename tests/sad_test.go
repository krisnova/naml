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

package tests

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestBeforeSadManifests(t *testing.T) {
	files, err := ioutil.ReadDir("sad")
	if err != nil {
		t.Errorf("unable to list test_nivenly.yaml directory: %v", err)
	}
	for _, file := range files {
		t.Logf("testing SAD [%s]", file.Name())
		err := generateCompileRunYAML(filepath.Join("sad", file.Name()))
		if err == nil {
			t.Errorf("expecting err with bad test")
			t.FailNow()
		}
	}
	t.Logf("Sad manifest tests complete")
}
