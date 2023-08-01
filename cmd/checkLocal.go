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

package cmd

import (
	"fmt"
	"os"

	txqualitychecks "github.com/eclipse-tractusx/tractusx-quality-checks/internal"
	"github.com/eclipse-tractusx/tractusx-quality-checks/internal/docs"
	"github.com/eclipse-tractusx/tractusx-quality-checks/internal/helm"
	"github.com/eclipse-tractusx/tractusx-quality-checks/internal/repo"
	"github.com/eclipse-tractusx/tractusx-quality-checks/pkg/container"
	"github.com/spf13/cobra"
)

var releaseGuidelines = []txqualitychecks.QualityGuideline{
	docs.NewReadmeExists(),
	docs.NewInstallExists(),
	docs.NewChangelogExists(),
	repo.NewLeadingRepositoryDefined(),
	repo.NewDefaultBranch(),
	container.NewAllowedBaseImage("./"),
	repo.NewRepoStructureExists(),
	helm.NewHelmStructureExists(),
}

// checkLocalCmd represents the checkLocal command
var checkLocalCmd = &cobra.Command{
	Use:   "checkLocal",
	Short: "Does run a quality check on local files",
	Long:  `Execute the checkLocal command in any directory you want to check for quality compliance with eclipse-tractusx rules`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running local checks of eclipse-tractusx release guidelines")
		runner := txqualitychecks.NewTestRunner(releaseGuidelines)
		err := runner.Run()

		if err != nil {
			fmt.Println("Error occurred! Check command output for details on failed checks")
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(checkLocalCmd)
}
