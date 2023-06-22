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
	"testing"

	"github.com/eclipse-tractusx/tractusx-quality-checks/internal/mocks"
)

func TestShouldPassIfNoDockerfilePresent(t *testing.T) {
	result := AllowedBaseImage{}.Test()

	if result == nil || result.Passed == false {
		t.Errorf("Allowed base image check should pass, if there is no Dockerfile found")
	}
}

func TestShouldFailIfDockerfileWithUnapprovedBaseImagePresent(t *testing.T) {
	dockerfile := createCorrettoDockerfile()
	defer dockerfile.Delete()

	result := AllowedBaseImage{}.Test()

	if result.Passed {
		t.Errorf("Allowed based image check should fail, if Dockerfile with unapproved base image found")
	}
}

func TestShouldPassIfTemurinIsUsedAsBaseImage(t *testing.T) {
	dockerfile := mocks.NewTempDockerfile().AppendCommand("FROM eclipse/temurin:17")
	_ = dockerfile.Create()
	defer dockerfile.Delete()

	if !(AllowedBaseImage{}.Test().Passed) {
		t.Errorf("eclipse/temurin should be recognized as approved base image")
	}
}

func TestShouldNotFailIfOnlyBuildLayerDeviatesFromTemurin(t *testing.T) {
	dockerfile := mocks.NewTempDockerfile().
		AppendCommand("FROM amazoncorretto:8 as builder").
		AppendCommand("COPY . .").
		AppendCommand("FROM eclipse/temurin:17")
	_ = dockerfile.Create()
	defer dockerfile.Delete()

	if !(AllowedBaseImage{}.Test().Passed) {
		t.Errorf("Unapproved images in build layers should not let the check fail")
	}
}

func createCorrettoDockerfile() *mocks.TempDockerfile {
	dockerfile := mocks.NewTempDockerfile()

	dockerfile.AppendCommand("FROM amazoncorreto:8").AppendEmptyLine().AppendCommand("COPY . .")

	_ = dockerfile.Create()
	return dockerfile
}
