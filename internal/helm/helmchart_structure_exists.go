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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/eclipse-tractusx/tractusx-quality-checks/internal"
	"github.com/eclipse-tractusx/tractusx-quality-checks/pkg/filesystem"
	"gopkg.in/yaml.v3"
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
	missingHCFields := []string{}
	for _, hc := range helmCharts {
		if hc.IsDir() {
			for _, fname := range helmStructureFiles {
				fpath := filepath.Join(mainDir, hc.Name(), fname)
				isMissing := filesystem.CheckMissingFiles([]string{fpath})
				if fname == "Chart.yaml" &&  isMissing == nil  {
					missingHCFields = verifyChartYaml(fpath)
				} else if isMissing != nil {
					missingFiles = append(missingFiles, isMissing...)
				}
			}
		}
	}

	if len(missingFiles) > 0 || len(missingHCFields) > 0 {
		return &txqualitychecks.QualityResult{ErrorDescription: buildErrorDescription(missingFiles, missingHCFields)}
	}

	return &txqualitychecks.QualityResult{Passed: true}
}

func (r *HelmStructureExists) IsOptional() bool {
	return false
}

func verifyChartYaml(ymlfile string) []string {
	type Chart struct {
		ApiVersion   string   `yaml:"apiVersion"`
		Name         string   `yaml:"name"`
		Description  string   `yaml:"description"`
		AppVersion   string   `yaml:"appVersion"`
		Version      string   `yaml:"version"`
		Home         string   `yaml:"home"`
		Sources      []string `yaml:"sources"`
		Dependencies []struct {
			Name       string `yaml:"name"`
			Repository string `yaml:"repository"`
			Version    string `yaml:"version"`
			Condition  string `yaml:"condition"`
		} `yaml:"dependencies"`
		Maintainers []struct {
			Name  string `yaml:"name"`
			Email string `yaml:"email"`
			Url   string `yaml:"url"`
		} `yaml:"maintainers"`
	}

	data, err := ioutil.ReadFile(ymlfile)
	if err != nil {
		fmt.Printf("Unable to read %v.\n", ymlfile)
		return nil
	}

	var chart Chart
	err = yaml.Unmarshal(data, &chart)
	if err != nil {
		fmt.Printf("Unable to parse YAML file: %v.\n", ymlfile)
		return nil
	}

	chartValues := reflect.ValueOf(chart)
	numFields := chartValues.NumField()
	chartType := chartValues.Type()

	missingFields := []string{} 

	for i := 0; i < numFields; i++ {
		field := chartType.Field(i)
		fieldValue := chartValues.Field(i)
	
		if fieldValue.Len() == 0 {
			missingFields = append(missingFields, field.Name)
		} 
	}
	if len(missingFields) > 0 {
		fmt.Println(missingFields)
		return missingFields
	}
	
	return nil
}

func buildErrorDescription(missingFiles []string, missingHCFields []string) string {
	errorDescription := ""

	if len(missingFiles) != 0 {
		errorDescription += "\n\tFollowing Helm Chart structure files are missing: "+strings.Join(missingFiles, ", ")
	}
	if len(missingHCFields) != 0 {
		errorDescription += "\n\tChart.yaml doesn't contain required fields: "+strings.Join(missingHCFields, ", ")
	}
	return errorDescription
}