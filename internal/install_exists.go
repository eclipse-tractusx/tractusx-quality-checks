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

import "os"

type InstallExists struct {
}

// NewInstallExists returns a new check based on QualityGuideline interface.
func NewInstallExists() QualityGuideline {
	return &InstallExists{}
}

func (r *InstallExists) Name() string {
	return "TRG 1.02 - INSTALL.md"
}

func (r *InstallExists) Description() string {
	return "File INSTALL.md contains comprehensive instructions for installation."
}

func (r *InstallExists) ExternalDescription() string {
	return "https://eclipse-tractusx.github.io/docs/release/trg-1/trg-1-2"
}

func (r *InstallExists) Test() *QualityResult {
	_, err := os.Stat("INSTALL.md")

	if err != nil {
		return &QualityResult{
			ErrorDescription: "Optional file INSTALL.md not found in current directory.",
		}
	}
	return &QualityResult{Passed: true}
}

func (r *InstallExists) IsOptional() bool {
	return true
}
