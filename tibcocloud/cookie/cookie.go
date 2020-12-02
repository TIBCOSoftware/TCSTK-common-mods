/*
* BSD 3-Clause License
* Copyright © 2020. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
*/
package cookie

import (
	"strings"
)

// Details is the stored structure for a TIBCO Cloud cookie
type Details struct {
	Region   string
	Tenant   string
	username string
	clientID string
	version  string
	Tsc      string
	Domain   string
}

type cookies map[string]Details

var currentCookies cookies = make(map[string]Details)

// Get returns the cookie for the connectionName
func Get(name string) (cookieDetails Details, ok bool) {
	cookieDetails, ok = currentCookies[name]
	return
}

// Set add/update a cookie for the connectionName
func Set(name string, cookie Details) {
	currentCookies[name] = cookie
	return
}

// New creates a new cookie
func New(region string, tenant string, username string, clientID string, version string, tsc string, domain string) (cookie Details) {
	cookie = Details{region, tenant, username, clientID, version, tsc, domain}
	return cookie
}

// Remove removes name cookie
func Remove(name string) {
	delete(currentCookies, name)
	return
}

// ExtractValue returns just the cookie value without Path value, Max-Age, HttpOnly or secure
func ExtractValue(cookie string, prefix string) (cookieValue string) {

	parsedCookie := strings.Split(cookie, ";")
	cookieValue = strings.TrimPrefix(parsedCookie[0], prefix)

	return
}
