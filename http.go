package library

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewNetClient() *http.Client {

	once.Do(func() {

		var netTransport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 60 * time.Second,
			}).DialContext,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout: 60 * time.Second,
		}

		netClient = &http.Client{
			Timeout:   time.Second * 60,
			Transport: otelhttp.NewTransport(netTransport),
		}
	})

	return netClient
}

func HTTPPost(url string, headers map[string]string, payload interface{}) (httpStatus int, response string) {

	if payload == nil {

		payload = "{}"
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, err.Error()
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
		return 0, err.Error()
	}

	st := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, err.Error()
	}

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
		return 0, err.Error()
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
		return 0, err.Error()
	}

	st := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, err.Error()
	}

	return st, string(body)
}

func HTTPRequest(ctx context.Context, method, url string, headers map[string]string, payload interface{}) (httpStatus int, response string) {

	var req *http.Request

	method = strings.ToUpper(strings.TrimSpace(method))

	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" && method != "PATCH" {

		return 0, "invalid method"
	}

	if payload == nil {

		jsonData, _ := json.Marshal(payload)
		reqL, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonData))
		if err != nil {

			log.Printf("got error making http request %s", err.Error())
			return 0, err.Error()
		}

		req = reqL

	} else {

		reqL, err := http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {

			log.Printf("got error making http request %s", err.Error())
			return 0, err.Error()
		}

		req = reqL
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if headers != nil {

		for k, v := range headers {

			req.Header.Set(k, v)
		}
	}

	resp, err := NewNetClient().Do(req)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, err.Error()
	}

	st := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, err.Error()
	}

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
		return 0, err.Error()
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
		return 0, err.Error()
	}

	st := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, err.Error()
	}

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
		return 0, err.Error()
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
		return 0, err.Error()
	}

	st := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, err.Error()
	}

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
		return 0, err.Error()
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
		return 0, err.Error()
	}

	defer resp.Body.Close()
	st := resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, err.Error()
	}

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
		return 0, err.Error()
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
		return 0, err.Error()
	}

	defer resp.Body.Close()
	st := resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return st, err.Error()
	}

	return st, string(body)
}

func ToMapStringInterface(d map[string]string) map[string]interface{} {

	e := make(map[string]interface{})

	for k, v := range d {
		e[k] = v
	}

	return e
}
