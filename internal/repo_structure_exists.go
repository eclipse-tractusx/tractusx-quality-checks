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

package txqualitychecks

import (
	"os"
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

func (c RepoStructureExists) Test() *QualityResult {

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
		"docs",
		"charts",
	}

	missingMandatoryFiles := checkMissingFiles(listOfMandatoryFilesToBeChecked)
	missingOptionalFiles := checkMissingFiles(listOfOptionalFilesToBeChecked)

	optionalMessage := "Warning! The check detected following optional files missing: "
	mandatoryMessage := "The check detected following mandatory files missing: "

	if len(missingOptionalFiles) > 0 {
		optionalMessage = fmtMessage(optionalMessage, missingOptionalFiles) + "\n\t"
	}

	if len(missingMandatoryFiles) > 0 {
		mandatoryMessage = fmtMessage(mandatoryMessage, missingMandatoryFiles)
	}

	if len(missingMandatoryFiles) == 0 && len(missingOptionalFiles) > 0 {
		return &QualityResult{ErrorDescription: optionalMessage, Passed: true}
	} else if len(missingMandatoryFiles) > 0 || len(missingOptionalFiles) > 0 {
		return &QualityResult{ErrorDescription: optionalMessage + mandatoryMessage}
	}

	return &QualityResult{Passed: true}
}

// Function to verify files existance.
// Return missing ones.
func checkMissingFiles(listOfFiles []string) []string {
	var missingFiles []string

	for _, file := range listOfFiles {

		if _, err := os.Stat(file); os.IsNotExist(err) {
			missingFiles = append(missingFiles, file)
		}
	}
	return missingFiles
}

// Function to format output message containing missing files.
func fmtMessage(startingMessage string, listOfFiles []string) string {
	message := startingMessage
	for _, missingFile := range listOfFiles {
		message += missingFile + " "
	}
	return message
}
