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

package docs

import (
	"os"
	"testing"
)

func TestNewInstallExists(t *testing.T) {
	t.Run("Provide ErrorDescription on Fail", func(t *testing.T) {
		installTest := NewInstallExists()
		expectedError := "Optional file INSTALL.md not found in current directory."

		result := installTest.Test()

		if result.ErrorDescription != expectedError {
			t.Errorf("Install.md check does not provide correct error description fon failing check!\n"+
				"expected: %s\ngot: %s", expectedError, result.ErrorDescription)
		}
	})
	t.Run("Pass if file exists", func(t *testing.T) {
		installTest := NewInstallExists()
		_, err := os.Create("INSTALL.md")
		if err != nil {
			t.Errorf("Error creating test preconditions!")
		}

		result := installTest.Test()

		if !result.Passed {
			t.Errorf("INSTALL exists, but test still fails")
		}

		err = os.Remove("INSTALL.md")
		if err != nil {
			t.Errorf("Error while removing test preconditions!")
		}
	})
	t.Run("No ErrorDescription if test pass", func(t *testing.T) {
		installTest := NewInstallExists()
		_, err := os.Create("INSTALL.md")
		if err != nil {
			t.Errorf("Error creating test preconditions!")
		}

		result := installTest.Test()

		if result.ErrorDescription != "" {
			t.Errorf("Passing test should not contain an error description")
		}

		err = os.Remove("INSTALL.md")
		if err != nil {
			t.Errorf("Error while removing test preconditions!")
		}
	})
}
