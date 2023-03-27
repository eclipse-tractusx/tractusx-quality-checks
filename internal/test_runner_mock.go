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

import "fmt"

type FailingQualityGuideline struct {
	name                string
	description         string
	externalDescription string
}

func (f FailingQualityGuideline) Name() string {
	return f.name
}

func (f FailingQualityGuideline) Description() string {
	return f.description
}

func (f FailingQualityGuideline) ExternalDescription() string {
	return f.externalDescription
}

func (f FailingQualityGuideline) Test() *QualityResult {
	return &QualityResult{Passed: false}
}

type PassingQualityGuideline struct {
	name                string
	description         string
	externalDescription string
}

func (p PassingQualityGuideline) Name() string {
	return p.name
}

func (p PassingQualityGuideline) Description() string {
	return p.description
}

func (p PassingQualityGuideline) ExternalDescription() string {
	return p.externalDescription
}

func (p PassingQualityGuideline) Test() *QualityResult {
	return &QualityResult{Passed: true}
}

type PrinterMock struct {
	messages []string
}

func (p *PrinterMock) Print(message string) {
	fmt.Println(message)
	p.messages = append(p.messages, message)
}
