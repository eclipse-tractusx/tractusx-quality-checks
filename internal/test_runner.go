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
	"errors"
	"fmt"
)

type GuidelineTestRunner struct {
	guidelines []QualityGuideline
	printer    Printer
}

func NewTestRunner(tests []QualityGuideline) *GuidelineTestRunner {
	return &GuidelineTestRunner{guidelines: tests, printer: &stdoutPrinter{}}
}

func (runner *GuidelineTestRunner) Run() error {
	allPassed := true
	for _, guideline := range runner.guidelines {
		runner.printer.Print(fmt.Sprintf("Testing Quality Guideline: %s", guideline.Name()))

		result := guideline.Test()
		if !result.Passed {
			runner.printer.Print(fmt.Sprintf("Failed! Guideline description: %s; More infos: %s", guideline.Description(), guideline.ExternalDescription()))
		}

		allPassed = allPassed && result.Passed
	}

	if !allPassed {
		return errors.New("not all tests have passed")
	}
	return nil
}

type stdoutPrinter struct {
}

func (p *stdoutPrinter) Print(message string) {
	fmt.Println(message)
}
