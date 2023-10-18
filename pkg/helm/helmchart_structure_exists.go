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
	"os"
	"path"
	"strings"

	"github.com/eclipse-tractusx/tractusx-quality-checks/internal/filesystem"
	"github.com/eclipse-tractusx/tractusx-quality-checks/internal/helm"
	"github.com/eclipse-tractusx/tractusx-quality-checks/pkg/tractusx"
)

var helmStructureFiles []string = []string{
	".helmignore",
	"LICENSE",
	"README.md",
	"values.yaml",
}

type HelmStructureExists struct {
	baseDir string
}

func NewHelmStructureExists(baseDir string) tractusx.QualityGuideline {
	return &HelmStructureExists{baseDir}
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

func (r *HelmStructureExists) IsOptional() bool {
	return false
}

func (r *HelmStructureExists) Test() *tractusx.QualityResult {
	var missingFiles []string
	var chartYamlFiles []string
	errorDescriptionCharts := ""
	chartsValid := true

	chartDir := path.Join(r.baseDir, "charts")
	helmCharts, err := os.ReadDir(chartDir)
	if err != nil || len(helmCharts) == 0 {
		return &tractusx.QualityResult{Passed: true}
	}

	for _, hc := range helmCharts {
		if !isChartDirectory(path.Join(chartDir, hc.Name())) {
			continue
		}
		chartYamlFiles = append(chartYamlFiles, path.Join(chartDir, hc.Name(), "Chart.yaml"))
		getMissingChartFiles(path.Join(chartDir, hc.Name()), &missingFiles)
	}

	if len(chartYamlFiles) > 0 {
		for _, fpath := range chartYamlFiles {
			isValid, msg := validateChart(fpath)
			chartsValid = chartsValid && isValid
			errorDescriptionCharts += msg
		}
	}

	if len(missingFiles) > 0 || !chartsValid {
		errMsg := "+ Following Helm Chart structure files are missing: " + strings.Join(missingFiles, ", ") + errorDescriptionCharts
		return &tractusx.QualityResult{ErrorDescription: errMsg, ErrorDescriptionWeb: strings.ReplaceAll(errMsg, "\n", "<br>")}
	}
	return &tractusx.QualityResult{Passed: true}
}

func isChartDirectory(dir string) bool {
	chartYamlPath := path.Join(dir, "Chart.yaml")
	_, err := os.Stat(chartYamlPath)
	return err == nil
}

func getMissingChartFiles(chartPath string, missingFiles *[]string) {
	for _, fileToCheck := range helmStructureFiles {
		missingFile := filesystem.CheckMissingFiles([]string{path.Join(chartPath, fileToCheck)})
		if missingFile != nil {
			*missingFiles = append(*missingFiles, missingFile...)
		}
	}
}

func validateChart(chartyamlfile string) (bool, string) {
	isValid := true
	returnMessage := "\n\t+ Analysis for " + chartyamlfile + ": "
	cyf := helm.ChartYamlFromFile(chartyamlfile)
	missingFields := cyf.GetMissingMandatoryFields()

	if len(missingFields) > 0 {
		isValid = false
		returnMessage += "\n\t\t - Missing mandatory fields: " + strings.Join(missingFields, ", ")
	}
	if !cyf.IsVersionValid() {
		isValid = false
		returnMessage += "\n\t\t - " + cyf.Version + " Version of the Helm Chart is incorrect. It needs to follow Semantic Version."
	}

	return isValid, returnMessage
}
