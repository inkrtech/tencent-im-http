/**
 * @Author: fuxiao
 * @Email: 576101059@qq.com
 * @Date: 2021/8/16 9:40 上午
 * @Desc: TODO
 */

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/inkrtech/tencent-im-http/internal"
)

const (
	MethodGet     = http.MethodGet
	MethodHead    = http.MethodHead
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodPatch   = http.MethodPatch
	MethodDelete  = http.MethodDelete
	MethodConnect = http.MethodConnect
	MethodOptions = http.MethodOptions
	MethodTrace   = http.MethodTrace
)

const fileUploadingKey = "@file:"

type Request struct {
	client        *Client
	retryCount    int
	retryInterval time.Duration
	Request       *http.Request
}

func NewRequest(c *Client) *Request {
	return &Request{
		client:        c,
		retryCount:    c.retryCount,
		retryInterval: c.retryInterval,
	}
}

func (r *Request) Next() (*Response, error) {
	if v := r.Request.Context().Value(middlewareKey); v != nil {
		if m, ok := v.(*middleware); ok {
			return m.Next()
		}
	}
	return r.call()
}

func (r *Request) request(method, url string, data ...interface{}) (resp *Response, err error) {
	r.Request, err = r.prepare(method, url, data...)
	if err != nil {
		return nil, err
	}

	if count := len(r.client.middlewares); count > 0 {
		handlers := make([]MiddlewareFunc, 0, count+1)
		handlers = append(handlers, r.client.middlewares...)
		handlers = append(handlers, func(r *Request) (*Response, error) {
			return r.call()
		})
		r.Request = r.Request.WithContext(context.WithValue(r.Request.Context(), middlewareKey, &middleware{
			req:      r,
			handlers: handlers,
			index:    -1,
		}))
		resp, err = r.Next()
	} else {
		resp, err = r.call()
	}

	return resp, err
}

// prepare build a http request.
func (r *Request) prepare(method, url string, data ...interface{}) (req *http.Request, err error) {
	method = strings.ToUpper(method)
	url = r.client.baseUrl + url

	var params string
	if len(data) > 0 {
		switch data[0].(type) {
		case string:
			params = data[0].(string)
		case []byte:
			params = string(data[0].([]byte))
		default:
			switch r.client.headers[HeaderContentType] {
			case ContentTypeJson:
				if b, err := json.Marshal(data[0]); err != nil {
					return nil, err
				} else {
					params = string(b)
				}
			case ContentTypeXml:
				if b, err := xml.Marshal(data[0]); err != nil {
					return nil, err
				} else {
					params = string(b)
				}
			default:
				params = internal.BuildParams(data[0])
			}
		}
	}

	if method == MethodGet {
		buffer := bytes.NewBuffer(nil)

		if params != "" {
			switch r.client.headers[HeaderContentType] {
			case ContentTypeJson, ContentTypeXml:
				buffer = bytes.NewBuffer([]byte(params))
			default:
				if strings.Contains(url, "?") {
					url = url + "&" + params
				} else {
					url = url + "?" + params
				}
			}
		}

		if req, err = http.NewRequest(method, url, buffer); err != nil {
			return nil, err
		}
	} else {
		if strings.Contains(params, fileUploadingKey) {
			var (
				buffer = bytes.NewBuffer(nil)
				writer = multipart.NewWriter(buffer)
			)

			for _, item := range strings.Split(params, "&") {
				array := strings.Split(item, "=")
				if len(array[1]) > 6 && strings.Compare(array[1][0:6], fileUploadingKey) == 0 {
					path := array[1][6:]
					if !internal.Exists(path) {
						return nil, errors.New(fmt.Sprintf(`"%s" does not exist`, path))
					}
					if file, err := writer.CreateFormFile(array[0], filepath.Base(path)); err == nil {
						if f, err := os.Open(path); err == nil {
							if _, err = io.Copy(file, f); err != nil {
								if err := f.Close(); err != nil {
									log.Printf(`%+v`, err)
								}
								return nil, err
							}
							if err := f.Close(); err != nil {
								log.Printf(`%+v`, err)
							}
						} else {
							return nil, err
						}
					} else {
						return nil, err
					}
				} else {
					if err = writer.WriteField(array[0], array[1]); err != nil {
						return nil, err
					}
				}
			}

			if err = writer.Close(); err != nil {
				return nil, err
			}

			if req, err = http.NewRequest(method, url, buffer); err != nil {
				return nil, err
			} else {
				req.Header.Set(HeaderContentType, writer.FormDataContentType())
			}
		} else {
			paramBytes := []byte(params)
			if req, err = http.NewRequest(method, url, bytes.NewReader(paramBytes)); err != nil {
				return nil, err
			} else {
				if v, ok := r.client.headers[HeaderContentType]; ok {
					req.Header.Set(HeaderContentType, v)
				} else if len(paramBytes) > 0 {
					if (paramBytes[0] == '[' || paramBytes[0] == '{') && json.Valid(paramBytes) {
						req.Header.Set(HeaderContentType, ContentTypeJson)
					} else if matched, _ := regexp.Match(`^[\w\[\]]+=.+`, paramBytes); matched {
						req.Header.Set(HeaderContentType, ContentTypeFormUrlEncoded)
					}
				}
			}
		}
	}

	if r.client.ctx != nil {
		req = req.WithContext(r.client.ctx)
	} else {
		req = req.WithContext(context.Background())
	}

	if len(r.client.headers) > 0 {
		for key, value := range r.client.headers {
			if key != "" {
				req.Header.Set(key, value)
			}
		}
	}

	if len(r.client.cookies) > 0 {
		var cookies = make([]string, 0)
		for key, value := range r.client.cookies {
			if key != "" {
				cookies = append(cookies, key+"="+value)
			}
		}
		req.Header.Set(HeaderCookie, strings.Join(cookies, ";"))
	}

	if host := req.Header.Get(HeaderHost); host != "" {
		req.Host = host
	}

	return req, nil
}

// call nitiate an HTTP request and return the response data.
func (r *Request) call() (resp *Response, err error) {
	resp = &Response{Request: r.Request}

	for {
		if resp.Response, err = r.client.Do(r.Request); err != nil {
			if resp.Response != nil {
				if err := resp.Response.Body.Close(); err != nil {
					log.Printf(`%+v`, err)
				}
			}

			if r.retryCount > 0 {
				r.retryCount--
				time.Sleep(r.retryInterval)
			} else {
				break
			}
		} else {
			break
		}
	}

	return resp, err
}
