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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
)

type Source struct {
	File *os.File
}

type Program struct {
	Source *Source
	File   *os.File
}

var mtx sync.Mutex

// Execute will execute a compiled NAML program
// stdout, stderr, err := program.Execute([]string{""})
func (p *Program) Execute(flags []string) (*bytes.Buffer, *bytes.Buffer, error) {
	mtx.Lock()
	defer mtx.Unlock()
	cmd := exec.Command(p.File.Name(), flags...)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	return stdout, stderr, err
}

// Remove will remove the program from the filesystem
func (p *Program) Remove() error {
	return os.Remove(p.File.Name())
}

// Compile will use the Go compiler to compile source code for NAML
func Compile(src []byte) (*Program, error) {

	dir := "/tmp"

	// Write a temporary source file
	f, err := ioutil.TempFile(dir, "*.go")
	if err != nil {
		return nil, fmt.Errorf("unable to write temporary file: %v", err)
	}

	_, err = f.Write(src)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to write source code to temporary file: %v", err)
	}

	source := &Source{
		File: f,
	}

	f, err = ioutil.TempFile(dir, "*.program")
	if err != nil {
		return nil, fmt.Errorf("unable to write program file: %v", err)
	}
	program := &Program{
		File: f,
	}

	cmd := exec.Command("go", "build", "-ldflags", "-X 'github.com/kris-nova/naml.Version=tests'", "-o", program.File.Name(), source.File.Name())
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		e := "\n\n"
		e = fmt.Sprintf("%s+-------------------------+---------------------------------\n", e)
		e = fmt.Sprintf("%s| Codify Compile Failure  |\n", e)
		e = fmt.Sprintf("%s+-------------------------+\n", e)
		e = fmt.Sprintf("%s| \n", e)
		e = fmt.Sprintf("%s| %s\n", e, stdout.String())
		e = fmt.Sprintf("%s| %s", e, stderr.String())
		e = fmt.Sprintf("%s+----------------------------------------------------------\n", e)
		return nil, fmt.Errorf("%s", e)
	}

	return program, nil
}

func Src(path string) ([]byte, error) {
	if path == "." {
		path = "main.go"
	}
	return ioutil.ReadFile(path)
}
