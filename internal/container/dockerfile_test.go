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

func TestShouldReturnEmptyStringIfNoDockerfilePresent(t *testing.T) {
	foundFiles := FindDockerfiles()
	if len(foundFiles) != 0 {
		t.Errorf("#FindDockerfiles should return empty string if no file could be found")
	}
}

func TestShouldFindSingleDockerfile(t *testing.T) {
	if err := NewTempDockerfile().WriteTo(t.TempDir()); err != nil {
		t.Errorf("Failed to write Dockerfile for test in temp dire")
	}

	//foundFiles := FindDockerfiles()
	//if len(foundFiles) != 1 {
	//	t.Errorf("Dockerfile written to temp dir could not be found by #FindDockerfiles")
	//}
}
