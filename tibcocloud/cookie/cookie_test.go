/*
* BSD 3-Clause License
* Copyright © 2020. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
*/
package cookie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCookie(t *testing.T) {
	cookie1 := New("region", "tenant", "username", "clientID", "version", "tsc", "domain")
	Set("test", cookie1)
	cookie2, _ := Get("test")
	assert.Equal(t, cookie1, cookie2)
}

func TestEmptyCurrentCookies(t *testing.T) {
	_, ok := Get("test")
	assert.Equal(t, ok, false)
}

func TestNonExistingCookie(t *testing.T) {
	initCookies()
	_, ok := Get("test")
	assert.Equal(t, ok, false)
}

func TestExistingCookie(t *testing.T) {
	initCookies()
	_, ok := Get("test1")
	assert.Equal(t, ok, true)
}

func TestDeleteCookie(t *testing.T) {
	initCookies()
	Remove("test2")
	_, ok := Get("test2")
	assert.Equal(t, ok, false)
}

func initCookies() {
	Set("test1", New("region1", "tenant1", "username1", "clientID1", "version1", "tsc1", "domain1"))
	Set("test2", New("region2", "tenant2", "username2", "clientID2", "version2", "tsc2", "domain2"))
	Set("test3", New("region3", "tenant3", "username3", "clientID3", "version3", "tsc3", "domain3"))
}
