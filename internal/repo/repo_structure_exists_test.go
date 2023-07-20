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

package repo

import (
	"github.com/eclipse-tractusx/tractusx-quality-checks/pkg/filesystem"
	"os"
	"testing"
)

var listOfFilesToBeCreated []string = []string{
	"CODE_OF_CONDUCT.md",
	"CONTRIBUTING.md",
	"DEPENDENCIES",
	"LICENSE",
	"NOTICE.md",
	"README.md",
	"SECURITY.md",
}

var listOfDirsToBeCreated []string = []string{
	"docs",
	"charts",
}

const metadataTestFile = "../../pkg/product/test/metadata_test_template.yaml"

func TestShouldPassIfRepoStructureExistsWithoutOptional(t *testing.T) {
	setEnv(t)
	defer os.Remove(".tractusx")

	filesystem.CreateFiles(listOfFilesToBeCreated)
	filesystem.CreateDirs(listOfDirsToBeCreated)

	repostructureTest := NewRepoStructureExists()
	result := repostructureTest.Test()
	filesystem.CleanFiles(append(listOfFilesToBeCreated, listOfDirsToBeCreated...))

	if !result.Passed {
		t.Errorf("Structure exists with optional files, but test still fails.")
	}
}

func TestShouldPassIfRepoStructureExistsWithOptional(t *testing.T) {
	setEnv(t)
	defer os.Remove(".tractusx")

	listOfFilesToBeCreated = append(listOfFilesToBeCreated, []string{"INSTALL.md", "AUTHORS.md"}...)

	filesystem.CreateFiles(listOfFilesToBeCreated)
	filesystem.CreateDirs(listOfDirsToBeCreated)

	repostructureTest := NewRepoStructureExists()
	result := repostructureTest.Test()
	filesystem.CleanFiles(append(listOfFilesToBeCreated, listOfDirsToBeCreated...))

	if !result.Passed {
		t.Errorf("Structure exists without optional files, but test still fails.")
	}

}

func TestShouldFailIfRepoStructureIsMissing(t *testing.T) {
	setEnv(t)
	defer os.Remove(".tractusx")

	repostructureTest := NewRepoStructureExists()

	result := repostructureTest.Test()

	if result.Passed {
		t.Errorf("RepoStructureExist should fail if repo structure exists.")
	}
}

func setEnv(t *testing.T) {
	copyTemplateFileTo(".tractusx", t)
	os.Setenv("GITHUB_REPOSITORY", "eclipse-tractusx/sig-infra")
	os.Setenv("GITHUB_REPOSITORY_OWNER", "tester")
}

func copyTemplateFileTo(path string, t *testing.T) {
	templateFile, err := os.ReadFile(metadataTestFile)
	if err != nil {
		t.Errorf("Could not read template file necessary for this test")
	}
	err = os.WriteFile(path, templateFile, 0644)
	if err != nil {
		t.Errorf("Could not copy template file to designated path")
	}
}
