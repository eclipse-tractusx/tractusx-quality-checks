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

package repo

import (
	"github.com/eclipse-tractusx/tractusx-quality-checks/internal"
	"github.com/eclipse-tractusx/tractusx-quality-checks/pkg/repo"
)

type defaultBranch struct {
}

func (d defaultBranch) Name() string {
	return "TRG 2.01 - Default Branch"
}

func (d defaultBranch) Description() string {
	return "Default branch must be main."
}

func (d defaultBranch) ExternalDescription() string {
	return "https://eclipse-tractusx.github.io/docs/release/trg-2/trg-2-1"
}

func (d defaultBranch) Test() *txqualitychecks.QualityResult {
	repoinfo := repo.GetRepoMetadata(repo.GetRepoBaseInfo())

	if *repoinfo.Fork {
		return &txqualitychecks.QualityResult{Passed: false, ErrorDescription: "Check determined running on a fork."}
	}

	if *repoinfo.DefaultBranch != "main" {
		return &txqualitychecks.QualityResult{
			Passed:           false,
			ErrorDescription: "Default branch not set to 'main'.",
		}
	}

	return &txqualitychecks.QualityResult{Passed: true}
}

func (d defaultBranch) IsOptional() bool {
	return false
}

func NewDefaultBranch() txqualitychecks.QualityGuideline {
	return &defaultBranch{}
}
