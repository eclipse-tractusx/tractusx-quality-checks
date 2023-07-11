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

func TestShouldFailIfChangelogFileIsMissing(t *testing.T) {
	changelogTest := NewChangelogExists()

	result := changelogTest.Test()

	if result.Passed {
		t.Errorf("ChangelogExist should fail if no Changelog file present")
	}
}

func TestShouldPassIfChangelogExists(t *testing.T) {
	_, _ = os.Create("CHANGELOG.md")
	defer func() {
		_ = os.Remove("CHANGELOG.md")
	}()
	changelogTest := NewChangelogExists()

	result := changelogTest.Test()

	if !result.Passed {
		t.Errorf("ChangelogExist should pass, if a CHANGELOG.md exists")
	}
}

func TestShouldProvideErrorDescriptionIfFailing(t *testing.T) {
	changelogTest := NewChangelogExists()

	result := changelogTest.Test()

	if result.ErrorDescription == "" {
		t.Errorf("Failing tests should provide an error description")
	}
}
