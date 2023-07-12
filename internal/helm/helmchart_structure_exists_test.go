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
	"github.com/eclipse-tractusx/tractusx-quality-checks/internal/repo"
	"os"
	"testing"
)

func TestShouldPassIfHelmDirIsMissing(t *testing.T) {
	helmStructureTest := NewHelmStructureExists()

	result := helmStructureTest.Test()

	if !result.Passed {
		t.Errorf("Helm directory doesn't exist hence test should pass.")
	}
}

func TestShouldFailIfNoHelmChartsFound(t *testing.T) {
	os.Mkdir("charts", 0750)
	defer os.Remove("charts")

	helmStructureTest := NewHelmStructureExists()
	result := helmStructureTest.Test()
	if result.Passed {
		t.Errorf("Helm directory doesn't contain any charts hence test should fail.")
	}

}

func TestShouldFailIfHelmStructureIsMissing(t *testing.T) {

	os.Mkdir("charts", 0750)
	os.Mkdir("charts/exampleChart", 0750)
	defer os.Remove("charts")

	helmStructureTest := NewHelmStructureExists()

	result := helmStructureTest.Test()

	if result.Passed {
		t.Errorf("Helm structure is missing hence test should fail.")
	}
}

func TestShouldPassIfHelmStructureExist(t *testing.T) {

	helmStructureDirsExample := []string{
		"charts",
		"charts/exampleChart",
		"charts/exampleChart/templates",
	}
	helmStructureFilesExample := []string{
		"charts/exampleChart/.helmignore",
		"charts/exampleChart/Chart.yaml",
		"charts/exampleChart/LICENSE",
		"charts/exampleChart/README.md",
		"charts/exampleChart/values.yaml",
		"charts/exampleChart/templates/NOTES.txt",
	}

	repo.CreateDirs(helmStructureDirsExample)
	repo.CreateFiles(helmStructureFilesExample)
	defer os.RemoveAll("charts")

	helmStructureTest := NewHelmStructureExists()

	result := helmStructureTest.Test()

	if !result.Passed {
		t.Errorf("Helm structure exists hence test should pass.")
	}

}
