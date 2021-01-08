/*
* BSD 3-Clause License
* Copyright © 2020. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package properties

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// ConfigProperties type
type ConfigProperties map[string]string

// ReadPropertiesFile function for read a property file
func ReadPropertiesFile(fileLocation string) (ConfigProperties, error) {

	config := ConfigProperties{}

	if len(fileLocation) == 0 {
		return config, nil
	}

	file, err := os.Open(fileLocation)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}

	if err := scan.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return config, nil
}
