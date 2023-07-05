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
	"bufio"
	"os"
	"strings"

	txqualitychecks "github.com/eclipse-tractusx/tractusx-quality-checks/internal"
)

type AllowedBaseImage struct {
}

func NewAllowedBaseImage() *AllowedBaseImage {
	return &AllowedBaseImage{}
}

func (a AllowedBaseImage) Name() string {
	return "TRG 4.02 - Base images"
}

func (a AllowedBaseImage) Description() string {
	return "We are aligning all product Docker images to a set of approved ones. This also makes it easier to properly annotate the used images as dependency"
}

func (a AllowedBaseImage) ExternalDescription() string {
	return "https://eclipse-tractusx.github.io/docs/release/trg-4/trg-4-02"
}

func (a AllowedBaseImage) Test() *txqualitychecks.QualityResult {
	_, err := os.Stat("Dockerfile")
	if err != nil {
		return &txqualitychecks.QualityResult{Passed: true}
	}

	file, err := os.Open("Dockerfile")
	defer file.Close()

	if err != nil {
		return &txqualitychecks.QualityResult{ErrorDescription: "Could not read file 'Dockerfile'"}
	}
	baseImage := readBaseImage(file)

	if !strings.Contains(baseImage, "eclipse/temurin") {
		return &txqualitychecks.QualityResult{ErrorDescription: "Docker base images other than eclipse/temurin are not approved. Please switch to Temurin"}
	}

	return &txqualitychecks.QualityResult{Passed: true}
}

func readBaseImage(file *os.File) string {
	lines := readLines(file)
	var baseImageLine string
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "FROM ") {
			baseImageLine = strings.TrimSpace(line)
		}
	}

	baseImage := strings.Trim(baseImageLine, "FROM ")
	return baseImage
}

func readLines(file *os.File) []string {
	var lines []string

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func (a AllowedBaseImage) IsOptional() bool {
	return false
}
