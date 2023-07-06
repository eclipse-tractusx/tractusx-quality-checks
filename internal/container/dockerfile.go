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
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

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
