/*******************************************************************************
 * Copyright (c) 2023 Contributors to the Eclipse Foundation
 *
 * See the NOTICE file(s) distributed with this work for additional
 * information regarding copyright ownership.
 *
 * This program and the accompanying materials are made available under the
 * terms of the Apache License, Version 2.0 which is available at
 * https://www.apache.org/licenses/LICENSE-2.0.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 ******************************************************************************/

package container

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	pathUtil "path"
	"path/filepath"
	"strings"
)

// dockerfile is a simple utility to create or read dockerfiles.
// the commands are supposed to contain every single instruction of the Dockerfile it represents
type dockerfile struct {
	filename string
	commands []string
}

func newDockerfile() *dockerfile {
	return &dockerfile{filename: "Dockerfile", commands: []string{}}
}

func dockerfileFromPath(path string) (*dockerfile, error) {
	fileInfo, err := os.Stat(path)
	if err != nil || fileInfo.IsDir() {
		return nil, errors.New("Could not find Dockerfile from path: " + path)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("Could not read Dockerfile from path: " + path)
	}

	// TODO: #readLines is a too simple approach. Multiline commands are perfectly valid
	return &dockerfile{filename: filepath.Base(fileInfo.Name()), commands: readLines(file)}, nil
}

func (d *dockerfile) appendCommand(command string) *dockerfile {
	d.commands = append(d.commands, command)
	return d
}

func (d *dockerfile) appendEmptyLine() *dockerfile {
	d.commands = append(d.commands, "")
	return d
}

func (d *dockerfile) writeTo(path string) error {
	if err := os.MkdirAll(path, 0770); err != nil {
		return err
	}

	file, err := os.Create(pathUtil.Join(path, d.filename))
	if err != nil {
		return err
	}
	defer file.Close()

	for _, command := range d.commands {
		file.WriteString(command + "\n")
	}

	return nil
}

// findDockerfilesAt will search the current repository recursively for Dockerfiles.
// If a file is found, the relative path to the file is returned in the result slice.
// If no Dockerfile is found the result will be an empty slice
func findDockerfilesAt(dir string) []string {
	fmt.Println("Start finding Dockerfiles at " + dir)
	var foundFiles []string

	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.Contains(info.Name(), "Dockerfile") {
			foundFiles = append(foundFiles, path)
		}
		return nil
	})

	fmt.Println("Found Dockerfiles: " + strings.Join(foundFiles, ", \n\t"))
	return foundFiles
}

func readLines(file *os.File) []string {
	var lines []string

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
