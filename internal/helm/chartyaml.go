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
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"reflect"
	"regexp"
)

type chartyaml struct {
	ApiVersion  string `yaml:"apiVersion"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	AppVersion  string `yaml:"appVersion"`
	Version     string `yaml:"version"`
	// Below Fields considered not mandatory
	// Commented out for future usage.
	//
	// Sources      []string `yaml:"sources"`
	// Home         string   `yaml:"home"`
	// Dependencies []struct {
	// 	Name       string `yaml:"name"`
	// 	Repository string `yaml:"repository"`
	// 	Version    string `yaml:"version"`
	// 	Condition  string `yaml:"condition"`
	// } `yaml:"dependencies"`
	// Maintainers []struct {
	// 	Name  string `yaml:"name"`
	// 	Email string `yaml:"email"`
	// 	Url   string `yaml:"url"`
	// } `yaml:"maintainers"`
}

func newChartYaml() *chartyaml {
	return &chartyaml{
		ApiVersion:  "",
		Name:        "",
		Description: "",
		AppVersion:  "",
		Version:     "",
	}
}

func chartYamlFromFile(ymlfile string) *chartyaml {
	data, err := ioutil.ReadFile(ymlfile)
	if err != nil {
		fmt.Printf("Unable to read %v.\n", ymlfile)
		return nil
	}

	var c chartyaml
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		fmt.Printf("Unable to parse YAML file: %v.\n", ymlfile)
		return nil
	}

	return &c
}

func (c *chartyaml) isVersionValid() bool {
	regexPattern := `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	match, err := regexp.MatchString(regexPattern, c.Version)
	if err != nil {
		fmt.Println("Error occured when validating semantic version.")
		return false
	}
	if !match {
		return false
	}
	return true
}

func (c *chartyaml) checkMandatoryFields() []string {
	chartValues := reflect.ValueOf(*c)
	numFields := chartValues.NumField()
	chartType := chartValues.Type()

	missingFields := []string{}
	for i := 0; i < numFields; i++ {
		field := chartType.Field(i)
		fieldValue := chartValues.Field(i)

		if fieldValue.Len() == 0 {
			missingFields = append(missingFields, field.Name)
		}
	}
	return missingFields
}
