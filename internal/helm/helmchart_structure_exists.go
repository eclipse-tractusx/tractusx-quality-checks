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
	"path/filepath"
	"strings"
)

type HelmStructureExists struct {
}

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
	chartYamlFiles := []string{}
	for _, hc := range helmCharts {
		if hc.IsDir() {
			for _, fname := range helmStructureFiles {
				fpath := filepath.Join(mainDir, hc.Name(), fname)
				isMissing := filesystem.CheckMissingFiles([]string{fpath})
				if fname == "Chart.yaml" && isMissing == nil {
					chartYamlFiles = append(chartYamlFiles, fpath)
				} else if isMissing != nil {
					missingFiles = append(missingFiles, isMissing...)
				}
			}
		}
	}

	errorDescriptionCharts := ""
	chartsValid := true
	if len(chartYamlFiles) > 0 {

		for _, fpath := range chartYamlFiles {
			isValid, msg := validateChart(fpath)
			chartsValid = chartsValid && isValid
			errorDescriptionCharts += msg
		}
	}

	if len(missingFiles) > 0 || !chartsValid {
		return &txqualitychecks.QualityResult{ErrorDescription: "+ Following Helm Chart structure files are missing: " + strings.Join(missingFiles, ", ") +
			errorDescriptionCharts}
	}
	return &txqualitychecks.QualityResult{Passed: true}
}

func (r *HelmStructureExists) IsOptional() bool {
	return false
}

func validateChart(chartyamlfile string) (bool, string) {
	isValid := true
	returnMessage := "\n\t+ Analysis for " + chartyamlfile + ": "
	cyf := chartYamlFromFile(chartyamlfile)
	missingFields := cyf.getMissingMandatoryFields()

	if len(missingFields) > 0 {
		isValid = false
		returnMessage += "\n\t\t - Missing mandatory fields: " + strings.Join(missingFields, ", ")
	}
	if !cyf.isVersionValid() {
		isValid = false
		returnMessage += "\n\t\t - " + cyf.Version + " Version of the Helm Chart is incorrect. It needs to follow Semantic Version."
	}

	return isValid, returnMessage
}
