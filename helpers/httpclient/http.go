package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/rayyone/go-core/errors"
	"github.com/rayyone/go-core/helpers/httpclient/contenttype"
	loghelper "github.com/rayyone/go-core/helpers/log"
	"github.com/rayyone/go-core/helpers/retry"
)

var (
	once       sync.Once
	httpClient *http.Client
)

// RequestOption Function to change request options
type RequestOption func(*requestOptions)

type requestOptions struct {
	Headers              map[string]string
	ErrorResult          interface{}
	ResponseType         interface{}
	ReportOnRequestError bool
	HTTPClient           *http.Client
	HTTPClientTimeout    time.Duration
	Debug                bool
	RetryOptions         retry.Options
}

func getDefaultRequestOptions() requestOptions {
	return requestOptions{
		Headers:              map[string]string{"Content-Type": string(contenttype.ApplicationJSON)},
		ResponseType:         nil,
		ReportOnRequestError: true,
		HTTPClientTimeout:    10 * time.Second,
	}
}

// SetDebug Set HTTP Client debug
func SetDebug(debug bool) RequestOption {
	return func(o *requestOptions) {
		o.Debug = debug
	}
}

// SetHttpClientTimeout Set HTTP Client timeout
func SetHttpClientTimeout(timeout time.Duration) RequestOption {
	return func(o *requestOptions) {
		o.HTTPClientTimeout = timeout
	}
}

// SetHttpClient Set a custom HTTP Client
func SetHttpClient(httpClient *http.Client) RequestOption {
	return func(o *requestOptions) {
		o.HTTPClient = httpClient
	}
}

// SetErrorResult Set error retriever obj
func SetErrorResult(errorResult interface{}) RequestOption {
	return func(o *requestOptions) {
		o.ErrorResult = &errorResult
	}
}

// AddHeader Add new header
func AddHeader(headerKey string, headerValue string) RequestOption {
	return func(o *requestOptions) {
		o.Headers[headerKey] = headerValue
	}
}

// RemoveHeader remove header by key
func RemoveHeader(headerKey string) RequestOption {
	return func(o *requestOptions) {
		delete(o.Headers, headerKey)
	}
}

// ContentType set content type header
func ContentType(contentType string) RequestOption {
	return func(o *requestOptions) {
		o.Headers["Content-Type"] = contentType
	}
}

// BearerTokenAuthorization add bearer authorization to headers
func BearerTokenAuthorization(token string) RequestOption {
	return func(o *requestOptions) {
		o.Headers["Authorization"] = "Bearer " + token
	}
}

// ReportOnRequestError set report on request error
func ReportOnRequestError(flag bool) RequestOption {
	return func(o *requestOptions) {
		o.ReportOnRequestError = flag
	}
}

// DontReportOnRequestError dont report on request error
func DontReportOnRequestError() RequestOption {
	return func(o *requestOptions) {
		o.ReportOnRequestError = false
	}
}

// WithRetry with retry
func WithRetry(retryOptions retry.Options) RequestOption {
	return func(o *requestOptions) {
		o.RetryOptions = retryOptions
	}
}

// Get GET Request
func Get(endpoint string, queryParams url.Values, result interface{}, opts ...RequestOption) error {
	endpointWithQuery, err := url.Parse(endpoint)
	if err != nil {
		errMsg := fmt.Sprintf("Error: API Call - Cannot parse URL. Error: %v", err)
		return errors.BadRequest.New(errMsg)
	}

	if queryParams != nil && len(queryParams) > 0 {
		endpointWithQuery.RawQuery = queryParams.Encode()
	}

	req, err := buildRequest(http.MethodGet, endpointWithQuery.String(), nil, opts...)
	if err != nil {
		return err
	}

	return SendRequest(req, result, opts...)
}

// Post POST Request
func Post(url string, payload BodyParams, result interface{}, opts ...RequestOption) error {
	return requestWithBodyParams(http.MethodPost, url, payload, result, opts...)
}

// Put PUT Request
func Put(url string, payload BodyParams, result interface{}, opts ...RequestOption) error {
	return requestWithBodyParams(http.MethodPut, url, payload, result, opts...)
}

// Delete DELETE Request
func Delete(url string, payload BodyParams, result interface{}, opts ...RequestOption) error {
	return requestWithBodyParams(http.MethodDelete, url, payload, result, opts...)
}

func requestWithBodyParams(method string, url string, payload BodyParams, result interface{}, opts ...RequestOption) error {
	options := getDefaultRequestOptions()
	for _, o := range opts {
		o(&options)
	}

	contentType := options.Headers["Content-Type"]
	isFormData := contentType == string(contenttype.FormData)
	var bodyParams io.Reader
	var err error
	if isFormData {
		bodyParams, contentType, err = getFormDataBodyParams(payload)
		if err != nil {
			errMsg := fmt.Sprintf("Error: API Call - Cannot build form data payload. Error: %v", err)
			return errors.BadRequest.New(errMsg)
		}
		opts = append(opts, ContentType(contentType))
	} else {
		payloadBs, err := json.Marshal(payload)
		if err != nil {
			errMsg := fmt.Sprintf("Error: API Call - Cannot encode payload. Error: %v", err)
			return errors.BadRequest.New(errMsg)
		}
		bodyParams = bytes.NewBuffer(payloadBs)
	}

	req, err := buildRequest(method, url, bodyParams, opts...)
	if err != nil {
		return err
	}

	return SendRequest(req, result, opts...)
}

func getFormDataBodyParams(payload BodyParams) (io.Reader, string, error) {
	body := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(body)

	for key, val := range payload {
		switch data := val.(type) {
		case []*multipart.FileHeader:
			fileHeaders := data
			for _, fileHeader := range fileHeaders {
				if !strings.Contains(key, "[]") {
					key = key + "[]"
				}
				if err := addFormFile(multipartWriter, key, fileHeader); err != nil {
					return nil, "", err
				}
			}
		case *multipart.FileHeader:
			fileHeader := data
			if err := addFormFile(multipartWriter, key, fileHeader); err != nil {
				return nil, "", err
			}
		default:
			str := fmt.Sprintf("%v", val)
			if err := multipartWriter.WriteField(key, str); err != nil {
				return nil, "", err
			}
		}
	}

	if err := multipartWriter.Close(); err != nil {
		return nil, "", err
	}

	return body, multipartWriter.FormDataContentType(), nil
}

func addFormFile(multipartWriter *multipart.Writer, key string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return errors.Validation.New("Cannot open file")
	}

	part, err := multipartWriter.CreateFormFile(key, fileHeader.Filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	return nil
}

func buildRequest(method string, url string, bodyParams io.Reader, opts ...RequestOption) (*http.Request, error) {
	options := getDefaultRequestOptions()
	for _, o := range opts {
		o(&options)
	}

	if options.Debug {
		loghelper.PrintYellowf("[HTTP Client] Requesting %s - %s", method, url)
		loghelper.PrintMagentaf("[HTTP Client] Body Params: %v", bodyParams)
	}

	req, err := http.NewRequest(method, url, bodyParams)
	if err != nil {
		errMsg := fmt.Sprintf("Error: API Call - Cannot init HTTP Request to '%s'. Error: %v", url, err)
		return nil, errors.BadRequest.New(errMsg)
	}

	return req, nil
}

func SendRequest(req *http.Request, result interface{}, opts ...RequestOption) (err error) {
	options := getDefaultRequestOptions()
	for _, o := range opts {
		o(&options)
	}
	client := options.HTTPClient
	if client == nil {
		client = getHTTPClient(options)
	}
	// Set to random user-agent so this request will not indentified as bot
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")
	for headerKey, headerValue := range options.Headers {
		req.Header.Set(headerKey, headerValue)
	}

	var bodyBs []byte
	var resp *http.Response
	err = retry.WithRetry(func() error {
		start := time.Now()
		resp, err = client.Do(req)
		if options.Debug {
			loghelper.PrintYellowf("[HTTP Client] Request completed in %.2fs", time.Since(start).Seconds())
		}
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			return errors.NewAndDontReport(fmt.Sprintf("Error: API Call - Cannot call API '%s'. Error: %v", req.URL, err))
		}

		bodyBs, _ = ioutil.ReadAll(resp.Body)

		// Read the body even the data is not important
		// This must to do, to avoid memory leak when reusing http
		// Connection. if you don't do this, http connection will be closed
		if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
			return errors.NewAndDontReport(fmt.Sprintf("Error: Cannot discard body response to dev null. Error: %v", err))
		}

		if resp.StatusCode >= 500 {
			return errors.NewAndDontReport(fmt.Sprintf("Error: API Call to '%s' - Returning status code of %d. Body Response: %s", resp.Request.URL, resp.StatusCode, bodyBs))
		}

		return nil
	}, options.RetryOptions)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		body := string(bodyBs)
		errType := errors.BadRequest
		if resp.StatusCode != 400 && (!options.ReportOnRequestError || resp.StatusCode == 422) {
			errType = errors.Validation
		}

		errors.SetExtra("json_response", body)
		if options.ErrorResult != nil {
			if err := json.Unmarshal(bodyBs, &options.ErrorResult); err != nil {
				return errType.Newf("Error: API Call to '%s' - Cannot unmarshal error response. Error: %v. Response: %s", resp.Request.URL, err, body)
			}
		}
		errMsg := "Error: API Call to '%s' - Returning status code of %d. Body Response: %s"
		return errType.Newf(errMsg, resp.Request.URL, resp.StatusCode, body)
	}

	if err := json.Unmarshal(bodyBs, &result); err != nil {
		return errors.BadRequest.Newf("Error: API Call to '%s' - Cannot unmarshal response. Error: %v", resp.Request.URL, err)
	}

	return nil
}

func getHTTPClient(options requestOptions) *http.Client {
	once.Do(func() {
		// Set idle connection to re-use the existing connection before timeout.
		defaultTransport := &http.Transport{
			DialContext: (&net.Dialer{
				KeepAlive: 600 * time.Second,
				Timeout:   10 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   100,
			ResponseHeaderTimeout: options.HTTPClientTimeout,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		httpClient = &http.Client{Transport: defaultTransport, Timeout: options.HTTPClientTimeout}
	})

	return httpClient
}
