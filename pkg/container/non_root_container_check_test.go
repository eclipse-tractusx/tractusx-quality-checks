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
	"fmt"
	"testing"
)

func TestShouldPassIfNoDockerfileFound(t *testing.T) {
	result := NewNonRootContainer().Test()

	if result == nil || result.Passed == false {
		t.Errorf("Non-Root Container Check should pass, if there is no Dockerfile found")
	}
}

// func TestValidateUser() validates the USER instruction in a Dockerfile to match lowercase or
func TestValidateUser(t *testing.T) {
	testCases := []struct {
		username         string
		groupname        string
		errorDescription string
		want             bool
	}{
		// TODO: carslen, 12.07.23: add missing error description

		{"1000", "1010", "lala", true},
		{"65537", "123", "uid gt 65536 is not possible", false},
		{"65536", "65537", "gid gt 65536 is not possible", false},
		{"hiddenrootgrp", "0", "gid = 0 not allowed", false},
		{"user123", "123", "", false},
		{"USERNAME", "", "", false},
		{"UserName", "", "", false},
		{"userName", "", "", false},
		//{NoUserFound, "", "empty user is invalid", false},
		{},
	}
	for _, tc := range testCases {
		result := validateUser(&user{tc.username, tc.groupname})
		if result != tc.want {
			t.Errorf("USER '%s':'%s' is not allowed: %s", tc.username, tc.groupname, tc.errorDescription)
		}
	}
}

func TestUserMethod(t *testing.T) {
	testCases := []struct {
		file *dockerfile
		want user
	}{
		{correttoDockerfile(), user{"corretto", ""}},                                                // USER corretto
		{temurinDockerfile(), user{"temurin", ""}},                                                  // USER temurin
		{newDockerfile().appendCommand("USER 0"), user{"0", ""}},                                    // uid as int without gid defined
		{newDockerfile().appendCommand("USER root"), user{"root", ""}},                              // uid as string without gid defined
		{newDockerfile().appendCommand("USER 0:1"), user{"0", "1"}},                                 // uid:gid as int
		{newDockerfile().appendCommand("USER 1000:0"), user{"1000", "0"}},                           // uid:gid as int
		{newDockerfile().appendCommand("USER testuser:testgroup"), user{"testuser", "testgroup"}},   // uid:gid as string
		{newDockerfile().appendCommand("FROM distroless").appendEmptyLine(), user{NoUserFound, ""}}, // No USER instruction in Dockerfile
	}

	for _, tc := range testCases {
		result := tc.file.user()

		// TODO: carslen, 12.07.23: datt tut nicht
		if result != &tc.want {
			t.Errorf("got user: '%s' group: '%s', expected user: '%s' group: '%s'", tc.file.user().user, tc.file.user().user, tc.want.user, tc.want.group)
		}
		fmt.Println(tc.file.filename)
	}
}

func TestQualityCheckPass2(t *testing.T) {
	testCases := []struct {
		file *dockerfile
		want bool
	}{
		{correttoDockerfile(), true},
		{temurinDockerfile(), true},
		{newDockerfile().appendCommand("FROM nginx").appendCommand("USER 102"), true},
		{newDockerfile().appendCommand("FROM nginx").appendCommand("USER 671234"), false},
		{newDockerfile().appendCommand("FROM nginx").appendCommand("USER 0"), false},
		{newDockerfile().appendCommand("FROM nginx").appendEmptyLine(), false},
	}

	for _, tc := range testCases {
		_ = tc.file.writeTo(".")
		result := NewNonRootContainer().Test()

		if result.Passed != tc.want {
			t.Errorf("got '%t', expected '%t' as result", result.Passed, tc.want)
		}
	}
}
