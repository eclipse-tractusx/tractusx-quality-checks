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
	"path"
	"testing"
)

func TestShouldReturnEmptyStringIfNoDockerfilePresent(t *testing.T) {
	foundFiles := FindDockerfilesAt("./")
	if len(foundFiles) != 0 {
		t.Errorf("#FindDockerfilesAt should return empty string if no file could be found")
	}
}

func TestShouldFindSingleDockerfile(t *testing.T) {
	tempDir := t.TempDir()
	if err := NewTempDockerfile().WriteTo(tempDir); err != nil {
		t.Errorf("Failed to write Dockerfile for test in temp dire")
	}

	foundFiles := FindDockerfilesAt(tempDir)
	if len(foundFiles) != 1 {
		t.Errorf("Dockerfile written to temp dir could not be found by #FindDockerfilesAt")
	}
}

func TestShouldFindDockerfileInSubdirectory(t *testing.T) {
	tempDir := t.TempDir()
	desiredPath := tempDir + "/abc/def"
	if err := os.MkdirAll(desiredPath, 0777); err != nil {
		t.Errorf("Could not create folder structure for test; err %s", err)
	}

	if err := NewTempDockerfile().WriteTo(desiredPath); err != nil {
		t.Errorf("Could not write test Dockerfile to desired temp subdirectory")
	}

	foundFiles := FindDockerfilesAt(tempDir)

	if len(foundFiles) == 0 {
		t.Errorf("Could not find Dockerfile in subdirectory")
	}
}

func TestShouldFindMultipleDockerfiles(t *testing.T) {
	tempDir := t.TempDir()
	secondTempDir := t.TempDir()

	NewTempDockerfile().WriteTo(tempDir)
	NewTempDockerfile().WriteTo(secondTempDir)

	foundFiles := FindDockerfilesAt(path.Join(tempDir, "../"))

	if len(foundFiles) != 2 {
		t.Errorf("Did not find all Dockerfiles")
	}
}

func TestShouldFindDockerfilesWithUnusualNames(t *testing.T) {
	tempDir := t.TempDir()
	dockerfile := NewTempDockerfile()
	dockerfile.FileName = "Dockerfile.full"
	dockerfile.WriteTo(tempDir)
	dockerfile.FileName = "db.Dockerfile"
	dockerfile.WriteTo(tempDir)
	dockerfile.FileName = "Dockerfile-dev"
	dockerfile.WriteTo(tempDir)

	foundFiles := FindDockerfilesAt(tempDir)

	if len(foundFiles) != 3 {
		t.Errorf("Could not find all uncommonly named Dockerfiles! Found %s", foundFiles)
	}
}
