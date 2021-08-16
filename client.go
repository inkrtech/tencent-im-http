/**
 * @Author: fuxiao
 * @Email: 576101059@qq.com
 * @Date: 2021/8/14 4:11 下午
 * @Desc: TODO
 */

package http

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Client struct {
	http.Client
	headers       map[string]string
	cookies       map[string]string
	ctx           context.Context
	baseUrl       string
	retryCount    int
	retryInterval time.Duration
	middlewares   []MiddlewareFunc
}

const (
	defaultUserAgent = "DobyteHttpClient"
	
	HeaderUserAgent     = "User-Agent"
	HeaderContentType   = "Content-Type"
	HeaderAuthorization = "Authorization"
	HeaderCookie        = "Cookie"
	HeaderHost          = "Host"
	
	ContentTypeJson           = "application/json"
	ContentTypeXml            = "application/xml"
	ContentTypeFormData       = "form-data"
	ContentTypeFormUrlEncoded = "application/x-www-form-urlencoded"
)

func NewClient() *Client {
	client := &Client{
		Client: http.Client{
			Transport: &http.Transport{
				DisableKeepAlives: true,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		headers:     make(map[string]string),
		cookies:     make(map[string]string),
		middlewares: make([]MiddlewareFunc, 0),
	}
	client.headers[HeaderUserAgent] = defaultUserAgent
	
	return client
}

// Set a header for the request.
func (c *Client) SetHeader(key, value string) *Client {
	c.headers[key] = value
	return c
}

// Set multiple headers for the request.
func (c *Client) SetHeaders(headers map[string]string) *Client {
	for key, value := range headers {
		c.headers[key] = value
	}
	return c
}

// Set a cookie for the request.
func (c *Client) SetCookie(key, value string) *Client {
	c.cookies[key] = value
	return c
}

// Set multiple cookies for the request.
func (c *Client) SetCookies(cookies map[string]string) *Client {
	for key, value := range cookies {
		c.cookies[key] = value
	}
	return c
}

// Set User-Agent for the request.
func (c *Client) SetUserAgent(agent string) *Client {
	c.headers[HeaderUserAgent] = agent
	return c
}

// Set Content-Type for the request.
func (c *Client) SetContentType(contentType string) *Client {
	c.headers[HeaderContentType] = contentType
	return c
}

// Enable browser mode for the request.
func (c *Client) SetBrowserMode() *Client {
	jar, _ := cookiejar.New(nil)
	c.Jar = jar
	return c
}

//
func (c *Client) SetBaseUrl(baseUrl string) *Client {
	c.baseUrl = baseUrl
	return c
}

// SetBasicAuth set HTTP basic authentication information for the request.
func (c *Client) SetBasicAuth(username, password string) *Client {
	c.headers[HeaderAuthorization] = "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	return c
}

// SetBearerToken set HTTP Bearer-Token authentication information for the request.
func (c *Client) SetBearerToken(token string) *Client {
	c.headers[HeaderAuthorization] = "Bearer " + token
	return c
}

// SetContext set context for the request.
func (c *Client) SetContext(ctx context.Context) *Client {
	c.ctx = ctx
	return c
}

// SetTimeOut sets the request timeout for the client.
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.Client.Timeout = timeout
	return c
}

// SetRetry sets count and interval of retry for the request.
func (c *Client) SetRetry(retryCount int, retryInterval time.Duration) *Client {
	c.retryCount = retryCount
	c.retryInterval = retryInterval
	return c
}

func (c *Client) SetKeepAlive(enable bool) {
	//c.Transport.
}

// Use sets middleware for the request.
func (c *Client) Use(middlewares ...MiddlewareFunc) *Client {
	c.middlewares = append(c.middlewares, middlewares...)
	return c
}

func (c *Client) Request(method, url string, data ...interface{}) (*Response, error) {
	return NewRequest(c).request(method, url, data...)
}

func (c *Client) Get(url string, data ...interface{}) (*Response, error) {
	return c.Request(MethodGet, url, data...)
}

func (c *Client) Post(url string, data ...interface{}) (*Response, error) {
	return c.Request(MethodPost, url, data...)
}

func (c *Client) Put(url string, data ...interface{}) (*Response, error) {
	return c.Request(MethodPut, url, data...)
}

func (c *Client) Patch(url string, data ...interface{}) (*Response, error) {
	return c.Request(MethodPatch, url, data...)
}

func (c *Client) Delete(url string, data ...interface{}) (*Response, error) {
	return c.Request(MethodDelete, url, data...)
}

func (c *Client) Head(url string, data ...interface{}) (*Response, error) {
	return c.Request(MethodHead, url, data...)
}

func (c *Client) Options(url string, data ...interface{}) (*Response, error) {
	return c.Request(MethodOptions, url, data...)
}

func (c *Client) Connect(url string, data ...interface{}) (*Response, error) {
	return c.Request(MethodConnect, url, data...)
}

func (c *Client) Trace(url string, data ...interface{}) (*Response, error) {
	return c.Request(MethodTrace, url, data...)
}
