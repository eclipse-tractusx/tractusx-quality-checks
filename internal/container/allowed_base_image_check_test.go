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

package container

import (
	"testing"
)

func TestShouldPassIfNoDockerfilePresent(t *testing.T) {
	result := AllowedBaseImage{}.Test()

	if result == nil || result.Passed == false {
		t.Errorf("Allowed base image check should pass, if there is no Dockerfile found")
	}
}

func TestShouldFailIfDockerfileWithUnapprovedBaseImagePresent(t *testing.T) {
	dockerfile := createCorrettoDockerfile()
	defer dockerfile.delete()

	result := AllowedBaseImage{}.Test()

	if result.Passed {
		t.Errorf("Allowed based image check should fail, if Dockerfile with unapproved base image found")
	}
}

func TestShouldPassIfTemurinIsUsedAsBaseImage(t *testing.T) {
	dockerfile := newTempDockerfile().appendCommand("FROM eclipse/temurin:17")
	_ = dockerfile.create()
	defer dockerfile.delete()

	if !(AllowedBaseImage{}.Test().Passed) {
		t.Errorf("eclipse/temurin should be recognized as approved base image")
	}
}

func TestShouldNotFailIfOnlyBuildLayerDeviatesFromTemurin(t *testing.T) {
	dockerfile := newTempDockerfile().
		appendCommand("FROM amazoncorretto:8 as builder").
		appendCommand("COPY . .").
		appendCommand("FROM eclipse/temurin:17")
	_ = dockerfile.create()
	defer dockerfile.delete()

	if !(AllowedBaseImage{}.Test().Passed) {
		t.Errorf("Unapproved images in build layers should not let the check fail")
	}
}

func createCorrettoDockerfile() *tempDockerfile {
	dockerfile := newTempDockerfile()

	dockerfile.appendCommand("FROM amazoncorreto:8").appendEmptyLine().appendCommand("COPY . .")

	_ = dockerfile.create()
	return dockerfile
}
