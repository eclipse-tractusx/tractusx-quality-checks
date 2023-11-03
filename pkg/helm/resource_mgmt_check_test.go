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

package helm

import (
	"os"
	"path"
	"testing"
)

func TestShouldPassIfDeploymentResourcesRequestsFound(t *testing.T) {
	dir := t.TempDir()
	setupK8SObject(dir, "deployment.yaml", "values.yaml", t)

	result := NewResourceMgmt(dir).Test()
	if !result.Passed {
		t.Errorf("Test should pass, the deployment manifest contains Resources Requests.")
	}
}

func TestShouldFailIfStatefulSetResourcesLimitsNotFound(t *testing.T) {
	dir := t.TempDir()
	setupK8SObject(dir, "statefulset.yaml", "values.yaml", t)

	result := NewResourceMgmt(dir).Test()
	if result.Passed {
		t.Errorf("Test should pass, the statefulset manifest doesn't contain Resources Limits.")
	}
}

func TestShouldFailIfCPURequestIsEqualCPULimitsAtDepl(t *testing.T) {
	dir := t.TempDir()
	setupK8SObject(dir, "deployment.yaml", "values-depequalcpu.yaml", t)

	result := NewResourceMgmt(dir).Test()
	if result.Passed {
		t.Errorf("Test should pass, CPU at Resources Requests is equal CPU at Resources Limits.")
	}
}

func TestShouldFailIfMemRequestIsNotEqualMemLimitsAtSTS(t *testing.T) {
	dir := t.TempDir()
	setupK8SObject(dir, "statefulset.yaml", "values-stsdiffmem.yaml", t)

	result := NewResourceMgmt(dir).Test()
	if result.Passed {
		t.Errorf("Test should pass, Mem at Resources Requests is not equal Mem at Resources Limits.")
	}
}

func copyFile(dest string, source string, t *testing.T) {
	templateFile, err := os.ReadFile(source)
	if err != nil {
		t.Errorf("Could not read source file: %s necessary for this test.", source)
	}
	err = os.WriteFile(dest, templateFile, 0770)
	if err != nil {
		t.Errorf("Could not copy template file to designated path.")
	}
}

func setupChartBasics(dir string, values string, t *testing.T) {
	testchartPath := path.Join("test", "charts", "testchart")
	_ = os.MkdirAll(path.Join(dir, "charts", "testchart", "templates"), 0770)
	copyFile(path.Join(dir, "charts", "testchart", "values.yaml"), path.Join(testchartPath, values), t)
	copyFile(path.Join(dir, "charts", "testchart", "Chart.yaml"), path.Join(testchartPath, "Chart.yaml"), t)
}

func setupK8SObject(dir string, manifest string, values string, t *testing.T) {
	setupChartBasics(dir, values, t)
	copyFile(path.Join(dir, "charts", "testchart", "templates", manifest), path.Join("test", "charts", "testchart", "templates", manifest), t)
}
