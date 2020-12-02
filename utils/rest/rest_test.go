/*
* BSD 3-Clause License
* Copyright © 2020. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
*/
package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLAAuthRestCall(t *testing.T) {

	var request Input
	request.URL = "https://eu.liveapps.cloud.tibco.com/idm/v3/login-oauth"
	request.HTTPMethod = "POST"
	request.ContentType = "application/x-www-form-urlencoded"

	data := "TenantId=tenant&Email=username@domain.com&Password=password&ClientID=clientId"

	request.Data = data

	response, _ := MakeCall(request)
	responseData := response.Data.(map[string]interface{})

	assert.Equal(t, "username@domain.com", responseData["userName"])
}
