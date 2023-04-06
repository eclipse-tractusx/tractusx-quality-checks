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

// InstallExists defines a struct for
type InstallExists struct {
}

// NewInstallExists func return a new check based on QualityGuideline interface.
func NewInstallExists() QualityGuideline {
	return &InstallExists{}
}

// Name implements ExternalDescription from interface
// QualityGuideline and returns the name of the test.
func (r *InstallExists) Name() string {
	return "TRG 1.02 - INSTALL.md"
}

// Description implements ExternalDescription from interface
// QualityGuideline and returns a brief description of this test.
func (r *InstallExists) Description() string {
	return "File INSTALL.md contains comprehensive instructions for installation."
}

// ExternalDescription implements ExternalDescription from interface
// QualityGuideline and returns an external description like a URL with further
// details for this test.
func (r *InstallExists) ExternalDescription() string {
	return "https://eclipse-tractusx.github.io/docs/release/trg-1/trg-1-2"
}

// IsOptional implements IsOptional from interface QualityGuideline to control if
// the test is optional.
func (r *InstallExists) IsOptional() bool {
	return true
}

// Test implements the test if optional file INSTALL.md exists. If file is
// missing an error description is added to QualityResult and Passed set to
// false, other it returns QualityResult Passed true.
func (r *InstallExists) Test() *QualityResult {
	_, err := os.Stat("INSTALL.md")

	if err != nil {
		return &QualityResult{
			ErrorDescription: "Optional file INSTALL.md not found in current directory.",
		}
	}
	return &QualityResult{Passed: true}
}
