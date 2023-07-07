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
	"fmt"
	"strings"

	txqualitychecks "github.com/eclipse-tractusx/tractusx-quality-checks/internal"
)

var baseImageAllowList = []string{
	"eclipse/temurin",
	"nginxinc/nginx-unprivileged",
	"mcr.microsoft.com/dotnet/runtime",
	"mcr.microsoft.com/dotnet/aspnet",
}

type AllowedBaseImage struct {
}

func NewAllowedBaseImage() txqualitychecks.QualityGuideline {
	return &AllowedBaseImage{}
}

func (a *AllowedBaseImage) Name() string {
	return "TRG 4.02 - Base images"
}

func (a *AllowedBaseImage) Description() string {
	return "We are aligning all product Docker images to a set of approved ones. This also makes it easier to properly annotate the used images as dependency"
}

func (a *AllowedBaseImage) ExternalDescription() string {
	return "https://eclipse-tractusx.github.io/docs/release/trg-4/trg-4-02"
}

func (a *AllowedBaseImage) Test() *txqualitychecks.QualityResult {
	foundDockerFiles := findDockerfilesAt("./")

	for _, dockerfilePath := range foundDockerFiles {
		file, err := dockerfileFromPath(dockerfilePath)
		if err != nil {
			fmt.Printf("Could not read dockerfile from Path %s\n", dockerfilePath)
		}

		if !isAllowedBaseImage(file.baseImage()) {
			return &txqualitychecks.QualityResult{ErrorDescription: "We want to align on docker base images. We detected a Dockerfile specifying " +
				file.baseImage() + "\n\tAllowed images are: \n\t - " +
				strings.Join(baseImageAllowList, "\n\t - ")}
		}
	}

	return &txqualitychecks.QualityResult{Passed: true}
}

func (a *AllowedBaseImage) IsOptional() bool {
	return false
}

func isAllowedBaseImage(image string) bool {
	for _, imageFromAllowList := range baseImageAllowList {
		if strings.Contains(image, imageFromAllowList) {
			return true
		}
	}
	return false
}
