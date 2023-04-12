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
	"log"
	"os"
	"strings"

	"github.com/go-ini/ini"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
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
	owner, repoName, isLocal := getRepoName()
	repo, err := getRepoInfo(owner, repoName, isLocal)

	if err != nil {
		log.Fatalf("error %v", err)
	}

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

// getRepoName returns repoName (string) and local (bool). repoName is either
// read from GITHUB_REPOSITORY environment variable when running as GH Action or,
// when executed locally, parsing local Git config.
//
// local it true when repoName was parsed from local git config.
func getRepoName() (owner string, repo string, local bool) {
	const (
		BASEURL = "https://github.com/"
		SSHBASE = "git@github.com:"
		SECTION = `remote "origin"`
		SUFFIX  = ".git"
	)

	if os.Getenv("GITHUB_REPOSITORY") != "" {
		// env variable is available when executed as GH Action/Workflow/Check
		owner = os.Getenv("GITHUB_REPOSITORY_OWNER")
		repo = strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")[1]

		return owner, repo, false
	} else {
		// Parse local git configuration when executing locally
		cfg, err := ini.Load(".git/config")
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		url := cfg.Section(SECTION).Key("url").String()

		var repoInfo []string
		if strings.Contains(url, BASEURL) {
			repoInfo = strings.Split(strings.TrimSuffix(strings.TrimPrefix(url, BASEURL), SUFFIX), "/")
		} else if strings.Contains(url, SSHBASE) {
			repoInfo = strings.Split(strings.TrimSuffix(strings.TrimPrefix(url, SSHBASE), SUFFIX), "/")
		}

		owner = repoInfo[0]
		repo = repoInfo[1]

		return owner, repo, true
	}
}

// getRepoInfo returns *github.Repository object and error. If GitHub API call
// failed, an error is returned.
func getRepoInfo(owner string, repo string, isLocal bool) (*github.Repository, error) {
	var client github.Client
	ctx := context.Background()
	if !isLocal {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = *github.NewClient(tc)
	} else {
		client = *github.NewClient(nil)
	}
	// Get repository information
	repoInfo, _, err := client.Repositories.Get(ctx, owner, repo)

	return repoInfo, err
}
