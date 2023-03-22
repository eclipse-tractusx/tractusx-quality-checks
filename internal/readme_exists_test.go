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
	"os"
	"testing"
)

func TestShouldFailIfReadmeDoesNotExist(t *testing.T) {
	readmeTest := NewReadmeExists()

	result := readmeTest.Test()

	if result.Passed {
		t.Errorf("Readme check should fail, if no README file present")
	}
}

func TestProvideErrorDescriptionOnFailingTest(t *testing.T) {
	readmeTest := NewReadmeExists()
	expectedError := "Did not find a README.md file in current directory!"

	result := readmeTest.Test()

	if result.ErrorDescription != expectedError {
		t.Errorf("Readme check does not provide correct error description on failing check! \nexprected: %s, \ngot: %s", expectedError, result.ErrorDescription)
	}
}

func TestShouldPassIfReadmeExists(t *testing.T) {
	readmeTest := NewReadmeExists()
	os.Create("README.md")

	result := readmeTest.Test()

	if !result.Passed {
		t.Errorf("README exists, but test still fails")
	}

	os.Remove("README.md")
}

func TestShouldNotProvideErrorDescriptionForPassingTest(t *testing.T) {
	readmeTest := NewReadmeExists()
	os.Create("README.md")

	result := readmeTest.Test()

	if result.ErrorDescription != "" {
		t.Errorf("Passing tests should not contain any error description")
	}

	os.Remove("README.md")
}
