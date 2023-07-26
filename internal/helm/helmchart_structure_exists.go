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

package helm

import (
	"github.com/eclipse-tractusx/tractusx-quality-checks/internal"
	"github.com/eclipse-tractusx/tractusx-quality-checks/pkg/filesystem"
	"os"
	"strings"
)

type HelmStructureExists struct {
}

// NewHelmStructureExists returns a new check based on QualityGuideline interface.
func NewHelmStructureExists() txqualitychecks.QualityGuideline {
	return &HelmStructureExists{}
}

func (r *HelmStructureExists) Name() string {
	return "TRG 5.02 - Chart structure"
}

func (r *HelmStructureExists) Description() string {
	return "Helm Chart should follow a specific structure."
}

func (r *HelmStructureExists) ExternalDescription() string {
	return "https://eclipse-tractusx.github.io/docs/release/trg-5/trg-5-02"
}

func (r *HelmStructureExists) Test() *txqualitychecks.QualityResult {
	helmStructureFiles := []string{
		".helmignore",
		"Chart.yaml",
		"LICENSE",
		"README.md",
		"values.yaml",
		"templates",
		"templates/NOTES.txt",
	}

	mainDir := "charts"
	if fi, err := os.Stat(mainDir); err != nil || !fi.IsDir() {
		return &txqualitychecks.QualityResult{Passed: true}
	}

	helmCharts, err := os.ReadDir(mainDir)
	if err != nil || len(helmCharts) == 0 {
		return &txqualitychecks.QualityResult{ErrorDescription: "Can't read Helm Charts at charts/."}
	}

	missingFiles := []string{}
	for _, hc := range helmCharts {
		if hc.IsDir() {
			tmpFilesStructure := []string{}
			for _, fname := range helmStructureFiles {
				tmpFilesStructure = append(tmpFilesStructure, mainDir+"/"+hc.Name()+"/"+fname)
			}
			missingFiles = append(missingFiles, filesystem.CheckMissingFiles(tmpFilesStructure)...)
		}
	}

	if len(missingFiles) > 0 {
		return &txqualitychecks.QualityResult{ErrorDescription: "Following files are missing: " + strings.Join(missingFiles, " ")}
	}

	return &txqualitychecks.QualityResult{Passed: true}
}

func (r *HelmStructureExists) IsOptional() bool {
	return false
}
