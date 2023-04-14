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

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-ini/ini"
	"github.com/google/go-github/v50/github"
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

func (d defaultBranch) Test() *QualityResult {
	repo := getRepoMetadata(getRepoBaseInfo())

	if *repo.Fork {
		return &QualityResult{Passed: false, ErrorDescription: "Check determined running on a fork."}
	}

	if *repo.DefaultBranch != "main" {
		return &QualityResult{
			Passed:           false,
			ErrorDescription: "Default branch not set to 'main'.",
		}
	}

	return &QualityResult{Passed: true}
}

func (d defaultBranch) IsOptional() bool {
	return false
}

func NewDefaultBranch() QualityGuideline {
	return &defaultBranch{}
}

// repoInfo struct provides required information to query GitHub API.
// repoInfo.local shall be true, when running check locally.
type repoInfo struct {
	owner    string
	reponame string
}

// getRepoBaseInfo returns repoInfo as pointer. Func determines according to
// available environment variables if running locally or as part of a GitHub
// Action/Workflow/Check.
func getRepoBaseInfo() *repoInfo {
	const (
		BASEURL = "https://github.com/"
		SSHBASE = "git@github.com:"
		SECTION = `remote "origin"`
		SUFFIX  = ".git"
	)
	if os.Getenv("GITHUB_REPOSITORY") != "" && os.Getenv("GITHUB_REPOSITORY_OWNER") != "" {
		// env variable is available when executed as GH Action/Workflow/Check
		result := repoInfo{
			owner:    os.Getenv("GITHUB_REPOSITORY_OWNER"),
			reponame: strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")[1],
		}

		return &result
	}

	// Parse local git configuration when executing locally
	cfg, err := ini.Load(".git/config")
	if err != nil {
		fmt.Printf("Failed to read file: %v", err)
	}

	url := cfg.Section(SECTION).Key("url").String()

	var repoSplitInfo []string
	if strings.Contains(url, BASEURL) {
		repoSplitInfo = strings.Split(strings.TrimSuffix(strings.TrimPrefix(url, BASEURL), SUFFIX), "/")
	} else if strings.Contains(url, SSHBASE) {
		repoSplitInfo = strings.Split(strings.TrimSuffix(strings.TrimPrefix(url, SSHBASE), SUFFIX), "/")
	}

	result := repoInfo{
		owner:    repoSplitInfo[0],
		reponame: repoSplitInfo[1],
	}

	return &result
}

// getRepoInfo returns *github.Repository object and error. If GitHub API call
// failed, an error is returned.
func getRepoMetadata(repo *repoInfo) *github.Repository {
	ctx := context.Background()
	client := *github.NewClient(nil)

	repoInfo, _, err := client.Repositories.Get(ctx, repo.owner, repo.reponame)
	if err != nil {
		fmt.Printf("Error querying GitHub API:\n%v\n", err)
	}

	return repoInfo
}
