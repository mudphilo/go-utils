package library

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func HTTPPost(url string, headers map[string]string, payload interface{}) (httpStatus int, response string) {

	if payload == nil {

		payload = "{}"
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	logHeaders := make(map[string]string)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	logHeaders["Content-Type"] = "application/json"
	logHeaders["Accept"] = "application/json"

	if headers != nil {

		for k, v := range headers {

			req.Header.Set(k, v)
			logHeaders[k] = v
		}
	}

	resp, err := NewNetClient().Do(req)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	st := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, ""
	}

	logRequest("POST", url, logHeaders, payload, st, req.Header, string(body))

	return st, string(body)
}

func HTTPPostWithContext(ctx context.Context, url string, headers map[string]string, payload interface{}) (httpStatus int, response string) {

	if payload == nil {

		payload = "{}"
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	logHeaders := make(map[string]string)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	logHeaders["Content-Type"] = "application/json"
	logHeaders["Accept"] = "application/json"

	if headers != nil {

		for k, v := range headers {

			req.Header.Set(k, v)
			logHeaders[k] = v
		}
	}

	resp, err := NewNetClient().Do(req)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	st := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, ""
	}

	logRequest("POST", url, logHeaders, payload, st, req.Header, string(body))

	return st, string(body)
}

func HTTPGet(remoteURL string, headers map[string]string, payload map[string]string) (httpStatus int, response string) {

	var fields []string

	endpoint := remoteURL
	if payload != nil {

		for key, value := range payload {

			val := fmt.Sprintf("%s=%v", key, url.QueryEscape(value))

			fields = append(fields, val)
		}

		params := strings.Join(fields, "&")

		endpoint = fmt.Sprintf("%s?%s", remoteURL, params)

	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	logHeaders := make(map[string]string)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	logHeaders["Content-Type"] = "application/json"
	logHeaders["Accept"] = "application/json"

	if headers != nil {

		for k, v := range headers {

			req.Header.Set(k, v)
			logHeaders[k] = v
		}
	}

	resp, err := NewNetClient().Do(req)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	st := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, ""
	}

	logRequest("GET", endpoint, logHeaders, nil, st, req.Header, string(body))

	return st, string(body)
}

func HTTPGetWithContext(ctx context.Context, remoteURL string, headers map[string]string, payload map[string]string) (httpStatus int, response string) {

	var fields []string

	endpoint := remoteURL

	if payload != nil {

		for key, value := range payload {

			val := fmt.Sprintf("%s=%v", key, url.QueryEscape(value))

			fields = append(fields, val)
		}

		params := strings.Join(fields, "&")

		endpoint = fmt.Sprintf("%s?%s", remoteURL, params)

	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	logHeaders := make(map[string]string)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	logHeaders["Content-Type"] = "application/json"
	logHeaders["Accept"] = "application/json"

	if headers != nil {

		for k, v := range headers {

			req.Header.Set(k, v)
			logHeaders[k] = v
		}
	}

	resp, err := NewNetClient().Do(req)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	st := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, ""
	}

	logRequest("GET", endpoint, logHeaders, nil, st, req.Header, string(body))

	return st, string(body)
}

func HTTPFormPost(endpoint string, headers map[string]string, payload map[string]string) (httpStatus int, response string) {

	method := "POST"

	var stringPayload []string

	if payload != nil {

		for key, value := range payload {

			stringPayload = append(stringPayload, fmt.Sprintf("%s=%v", key, value))

		}

	}

	requestPayload := strings.NewReader(strings.Join(stringPayload, "&"))

	req, err := http.NewRequest(method, endpoint, requestPayload)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	logHeaders := make(map[string]string)
	logHeaders["Content-Type"] = "application/x-www-form-urlencoded"

	if headers != nil {

		for k, v := range headers {

			req.Header.Set(k, v)
			logHeaders[k] = v
		}
	}

	resp, err := NewNetClient().Do(req)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	defer resp.Body.Close()
	st := resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, ""
	}

	logRequest("POST", endpoint, logHeaders, payload, st, req.Header, string(body))

	return st, string(body)
}

func HTTPFormPostWithContext(ctx context.Context, endpoint string, headers map[string]string, payload map[string]string) (httpStatus int, response string) {

	method := "POST"

	var stringPayload []string

	if payload != nil {

		for key, value := range payload {

			stringPayload = append(stringPayload, fmt.Sprintf("%s=%v", key, value))

		}

	}

	requestPayload := strings.NewReader(strings.Join(stringPayload, "&"))

	req, err := http.NewRequestWithContext(ctx, method, endpoint, requestPayload)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	logHeaders := make(map[string]string)
	logHeaders["Content-Type"] = "application/x-www-form-urlencoded"

	if headers != nil {

		for k, v := range headers {

			req.Header.Set(k, v)
			logHeaders[k] = v
		}
	}

	resp, err := NewNetClient().Do(req)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	defer resp.Body.Close()
	st := resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, ""
	}

	logRequest("FORM", endpoint, logHeaders, payload, st, req.Header, string(body))

	return st, string(body)
}

func logRequest(method, endpoint string, requestHeaders map[string]string, requestBody interface{},
	responseStatus int, responseHeader http.Header, responseBody string) {

	if os.Getenv("debug") == "1" || os.Getenv("DEBUG") == "1" {

		// sensitive headers (lowercase)
		sensitiveHeaders := map[string]struct{}{
			"authorization":       {},
			"x-api-key":           {},
			"x-auth-token":        {},
			"cookie":              {},
			"set-cookie":          {},
			"client-secret":       {},
			"proxy-authorization": {},
		}

		maskIfSensitive := func(key, value string) string {
			normalized := strings.ToLower(http.CanonicalHeaderKey(key))
			if _, exists := sensitiveHeaders[normalized]; exists {
				return "***"
			}
			return value
		}

		responseHeaders := make(map[string]string)
		for k, v := range responseHeader {
			joined := strings.Join(v, ",")
			responseHeaders[k] = maskIfSensitive(k, joined)
		}

		var heads, rheads []string

		for k, v := range requestHeaders {
			heads = append(heads, fmt.Sprintf("\t%s : %s", k, maskIfSensitive(k, v)))
		}

		for k, v := range responseHeaders {
			rheads = append(rheads, fmt.Sprintf("\t%s : %s", k, v))
		}

		body := "none"
		if requestBody != nil {
			jsonData, _ := json.Marshal(requestBody)
			body = string(jsonData)
		}

		log.Printf("**** BEGIN HTTP %s REQUEST ****\n"+
			"Remote Url : %s\n"+
			"Request Headers:\n"+
			"%s\n"+
			"Request Payload\n"+
			"\t%s\n"+
			"Response Status: %d\n"+
			"Response Headers\n"+
			"%s\n"+
			"Response Body\n"+
			"**** END HTTP %s REQUEST ****\n"+
			"\t%s",
			strings.ToUpper(method),
			endpoint,
			strings.Join(heads, "\n"),
			body,
			responseStatus,
			strings.Join(rheads, "\n"),
			strings.ToUpper(method),
			responseBody,
		)
	}
}
func ToMapStringInterface(d map[string]string) map[string]interface{} {

	e := make(map[string]interface{})

	for k, v := range d {
		e[k] = v
	}

	return e
}
