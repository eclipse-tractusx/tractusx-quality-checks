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
	"os"
	pathUtil "path"
)

type TempDockerfile struct {
	FileName string
	Commands []string
}

func NewTempDockerfile() *TempDockerfile {
	return &TempDockerfile{FileName: "Dockerfile", Commands: []string{}}
}

func (f *TempDockerfile) AppendCommand(command string) *TempDockerfile {
	f.Commands = append(f.Commands, command)
	return f
}

func (f *TempDockerfile) AppendEmptyLine() *TempDockerfile {
	f.Commands = append(f.Commands, "")
	return f
}

func (f *TempDockerfile) Create() error {
	return f.WriteTo("./")
}

func (f *TempDockerfile) Delete() error {
	if err := os.Remove(f.FileName); err != nil {
		return err
	}
	return nil
}

func (f *TempDockerfile) WriteTo(path string) error {
	if err := os.MkdirAll(path, 0770); err != nil {
		return err
	}

	file, err := os.Create(pathUtil.Join(path, f.FileName))
	if err != nil {
		return err
	}
	defer file.Close()

	for _, command := range f.Commands {
		file.WriteString(command + "\n")
	}

	return nil
}
