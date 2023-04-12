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

package product

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const MetadataFilename = ".tractusx"

type Metadata struct {
	ProductName       string       `yaml:"product"`
	LeadingRepository string       `yaml:"leadingRepository"`
	Repositories      []Repository `yaml:"repositories"`
}

type Repository struct {
	Name             string `yaml:"name"`
	UsageDescription string `yaml:"usage"`
	Url              string `yaml:"url"`
}

func MetadataFromFile(fileContent []byte) (*Metadata, error) {
	var metadata Metadata

	err := yaml.Unmarshal(fileContent, &metadata)
	if err != nil {
		fmt.Println("Error parsing product metadata file")
		return nil, err
	}

	return &metadata, nil
}

func MetadataFromLocalFile() (*Metadata, error) {
	metadataFileAsBytes, err := os.ReadFile(MetadataFilename)

	if err != nil {
		fmt.Println(fmt.Sprintf("Could not read Tractus-X metadatafile from default location: %s", MetadataFilename))
		return nil, err
	}

	file, err := MetadataFromFile(metadataFileAsBytes)

	if err != nil {
		return nil, err
	}

	return file, nil
}
