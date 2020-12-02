/*
* BSD 3-Clause License
* Copyright © 2020. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package rest

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/TCSTK-common-mods/tibcocloud/cookie"
)

// Input for the rest call
type Input struct {
	URL         string
	HTTPMethod  string
	ContentType string // Could be replace for instrospecting Data

	Headers        map[string]string
	Authentication interface{}

	Return string // Option to specify the expected output

	// Path Params
	PathParams map[string]string

	// Query Params
	QueryParams url.Values

	// JSON elements for the REST calls
	Data interface{} `json:"data"`

	// Responses containers
	Response interface{} `json:"Response"` // The message when HTTP Code is 200
	Error    interface{} // The error message when HTTP Code is not 200
}

// Output response
type Output struct {
	Status  int         // The HTTP status code
	Headers http.Header // The HTTP response headers
	Data    interface{} // The HTTP response data
}

// MakeCall to debug retry
func MakeCall(request Input) (response Output, err error) {
	for i := 0; i < 3; i++ {
		response, err = MakeCall2(request)
		if response.Status <= 500 {
			// fmt.Printf("***** Status code: [%d]\n", response.Status)
			break
		}
	}
	return response, err
}

// MakeCall2 to the rest service
func MakeCall2(request Input) (response Output, err error) {

	var reqBody io.Reader
	if request.HTTPMethod == http.MethodPut ||
		request.HTTPMethod == http.MethodPost ||
		request.HTTPMethod == http.MethodPatch ||
		request.HTTPMethod == http.MethodDelete {
		if request.ContentType != "application/multipart" {
			if request.Data != nil {
				if str, ok := request.Data.(string); ok {
					reqBody = bytes.NewBuffer([]byte(str))
				} else if data, ok := request.Data.(url.Values); ok {
					reqBody = strings.NewReader(data.Encode())
				} else {
					b, _ := json.Marshal(request.Data) //todo handle error
					reqBody = bytes.NewBuffer([]byte(b))
				}
			}
		} else {
			dataByte, _ := base64.StdEncoding.DecodeString(request.Data.(string))
			data := bytes.NewReader(dataByte)
			// data := strings.NewReader(fileContent)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("artifactContents", filepath.Base(request.URL))
			if err != nil {
				// return "nil", err
			}
			_, err = io.Copy(part, data)

			err = writer.Close()
			if err != nil {
				// return "nil", err
			}
			request.ContentType = writer.FormDataContentType()
			reqBody = body
		}
	} else {
		reqBody = nil
	}

	if len(request.PathParams) > 0 {
		request.URL = generateURL(request.URL, request.PathParams)
	}

	if len(request.QueryParams) > 0 {
		request.URL = request.URL + "?" + request.QueryParams.Encode()
	}

	req, err := http.NewRequest(request.HTTPMethod, request.URL, reqBody)
	if err != nil {

	}

	// Add headers
	if len(request.Headers) > 0 {
		for k, v := range request.Headers {
			req.Header.Set(k, v)
		}
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", request.ContentType)
	}

	// Add security header
	if request.Authentication != nil {
		cookie, ok := request.Authentication.(cookie.Details)
		if ok {
			req.AddCookie(&http.Cookie{Name: "tsc", Value: cookie.Tsc})
			req.AddCookie(&http.Cookie{Name: "domain", Value: cookie.Domain})
		} else {
			req.Header.Set("Authorization", request.Authentication.(string))
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}

	if resp == nil {
		return response, nil
	}

	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	// Check the HTTP Header Content-Type
	respContentType := resp.Header.Get("Content-Type")

	if request.Return != "" {
		respContentType = strings.ToLower(request.Return)
	}

	var responseData interface{}
	if resp.StatusCode == 200 {
		responseData = request.Response
	} else {
		responseData = request.Error
	}

	switch request.Response.(type) {
	case string:
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return response, err
		}
		responseData = string(b)
	default:
		switch respContentType {
		case "application/json",
			"application/json;charset=utf-8",
			"application/json; charset=utf-8":
			d := json.NewDecoder(resp.Body)
			d.UseNumber()
			err = d.Decode(responseData)
			if err != nil {
				switch {
				case err == io.EOF:
					// empty body
				default:
					return response, err
				}
			}
		case "int":
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return response, err
			}
			responseData, _ = strconv.Atoi(string(b))
		default:
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return response, err
			}
			responseData = string(b)
		}
	}

	response = Output{Status: resp.StatusCode, Data: responseData, Headers: resp.Header}

	return response, nil
}

func generateURL(uri string, values map[string]string) string {

	var buffer bytes.Buffer
	buffer.Grow(len(uri))

	addrStart := strings.Index(uri, "://")

	i := addrStart + 3

	for i < len(uri) {
		if uri[i] == '/' {
			break
		}
		i++
	}

	buffer.WriteString(uri[0:i])

	for i < len(uri) {
		if uri[i] == ':' {
			j := i + 1
			for j < len(uri) && uri[j] != '/' {
				j++
			}

			if i+1 == j {

				buffer.WriteByte(uri[i])
				i++
			} else {

				param := uri[i+1 : j]
				value := values[param]
				buffer.WriteString(value)
				if j < len(uri) {
					buffer.WriteString("/")
				}
				i = j + 1
			}

		} else {
			buffer.WriteByte(uri[i])
			i++
		}
	}

	return buffer.String()
}

// GenerateBasicToken generate a Basic Authentication
func GenerateBasicToken(username string, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}
