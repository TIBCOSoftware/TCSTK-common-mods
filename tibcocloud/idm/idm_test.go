/*
* BSD 3-Clause License
* Copyright © 2020. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
*/
package idm

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ConfigProperties map[string]string

var fileLocation = "../../credentials.properties"

func ReadPropertiesFile() (ConfigProperties, error) {
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

func TestCreateCookie(t *testing.T) {

	props, err := ReadPropertiesFile()
	if err != nil {
		panic("Error while reading properties file")
	}

	cookie, _ := LoginOauth("kkk", props["region"], props["tenant"], props["username"], props["password"], props["clientid"], props["version"])

	assert.NotEmpty(t, cookie.Domain)
}
