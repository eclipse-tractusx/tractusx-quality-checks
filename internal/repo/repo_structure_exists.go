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
	"fmt"
	"os"
	"strings"

	"github.com/eclipse-tractusx/tractusx-quality-checks/internal"
	"github.com/eclipse-tractusx/tractusx-quality-checks/pkg/product"
)

type RepoStructureExists struct {
}

func (c RepoStructureExists) IsOptional() bool {
	return false
}

func NewRepoStructureExists() *RepoStructureExists {
	return &RepoStructureExists{}
}

func (c RepoStructureExists) Name() string {
	return "TRG 2.03 - Repo structure"
}

func (c RepoStructureExists) Description() string {
	return "All repositories must follow specified files and folders structure."
}

func (c RepoStructureExists) ExternalDescription() string {
	return "https://eclipse-tractusx.github.io/docs/release/trg-2/trg-2-3"
}

func (c RepoStructureExists) Test() *txqualitychecks.QualityResult {

	// Slice containing required files and folders in the repo structure.
	// Before modification make sure you align to TRG 2.03 guideline.

	listOfOptionalFilesToBeChecked := []string{
		"AUTHORS.md",
		"INSTALL.md",
	}

	listOfMandatoryFilesToBeChecked := []string{
		"CODE_OF_CONDUCT.md",
		"CONTRIBUTING.md",
		"DEPENDENCIES",
		"LICENSE",
		"NOTICE.md",
		"README.md",
		"SECURITY.md",
	}

	mandatoryForLeadingRepo := []string{"docs", "charts"}

	if product.IsLeadingRepo() {
		listOfMandatoryFilesToBeChecked = append(listOfMandatoryFilesToBeChecked, mandatoryForLeadingRepo...)
	}

	missingMandatoryFiles := CheckMissingFiles(listOfMandatoryFilesToBeChecked)
	missingOptionalFiles := CheckMissingFiles(listOfOptionalFilesToBeChecked)

	if len(missingOptionalFiles) > 0 {
		fmt.Printf("Warning! Guideline description: %s\n\t%s\n\tMore infos: %s\n",
			c.Description(), "The check detected following optional files missing: "+strings.Join(missingOptionalFiles, " "),
			c.ExternalDescription())
	}

	if len(missingMandatoryFiles) > 0 {
		return &txqualitychecks.QualityResult{ErrorDescription: "The check detected following mandatory files missing: " + strings.Join(missingMandatoryFiles, " ")}
	}

	return &txqualitychecks.QualityResult{Passed: true}
}

// Function to verify files existance.
// Return missing ones.
func CheckMissingFiles(listOfFiles []string) []string {
	var missingFiles []string

	for _, file := range listOfFiles {

		if _, err := os.Stat(file); os.IsNotExist(err) {
			missingFiles = append(missingFiles, file)
		}
	}
	return missingFiles
}

func CreateFiles(files []string) {
	for _, file := range files {
		os.Create(file)
	}
}

func CreateDirs(dirs []string) {
	for _, dir := range dirs {
		os.Mkdir(dir, 0750)
	}
}

func CleanFiles(files []string) {

	for _, file := range files {
		os.Remove(file)
	}
}
