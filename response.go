/**
 * @Author: fuxiao
 * @Email: 576101059@qq.com
 * @Date: 2021/8/15 4:56 下午
 * @Desc: TODO
 */

package http

import (
	"io/ioutil"
	"net/http"
	"sync"
	
	"github.com/dobyte/http/internal"
)

type Response struct {
	*http.Response
	Request *http.Request
	body    []byte
	cookies map[string]string
	mu      sync.Mutex
}

// ReadBytes retrieves and returns the response content as []byte.
func (r *Response) ReadBytes() []byte {
	if r == nil || r.Response == nil {
		return []byte{}
	}
	
	if r.body == nil {
		var err error
		r.mu.Lock()
		defer r.mu.Unlock()
		if r.body, err = ioutil.ReadAll(r.Response.Body); err != nil {
			return nil
		}
	}
	
	return r.body
}

// ReadString retrieves and returns the response content as string.
func (r *Response) ReadString() string {
	return internal.UnsafeBytesToString(r.ReadBytes())
}

// Scan convert the response into a complex data structure.
func (r *Response) Scan(any interface{}) error {
	return internal.Scan(r.ReadBytes(), any)
}

// Close closes the response when it will never be used.
func (r *Response) Close() error {
	if r == nil || r.Response == nil || r.Response.Close {
		return nil
	}
	r.Response.Close = true
	return r.Response.Body.Close()
}

// HasHeader Determine if a header exists in the cache.
func (r *Response) HasHeader(key string) bool {
	for k, _ := range r.Header {
		if k == key {
			return true
		}
	}
	
	return false
}

// GetHeader Retrieve header's value from the response.
func (r *Response) GetHeader(key string) string {
	return r.Header.Get(key)
}

// GetHeader Retrieve all header's value from the response.
func (r *Response) GetHeaders() map[string]interface{} {
	headers := make(map[string]interface{})
	for k, v := range r.Header {
		if len(v) > 1 {
			headers[k] = v
		} else {
			headers[k] = v[0]
		}
	}
	
	return headers
}

// HasCookie Determine if a cookie exists in the cache.
func (r *Response) HasCookie(key string) bool {
	if r.cookies == nil {
		r.cookies = r.GetCookies()
	}
	_, ok := r.cookies[key]
	
	return ok
}

// GetCookie Retrieve cookie's value from the response.
func (r *Response) GetCookie(key string) string {
	if r.cookies == nil {
		r.cookies = r.GetCookies()
	}
	return r.cookies[key]
}

// GetCookies Retrieve all cookie's value from the response.
func (r *Response) GetCookies() map[string]string {
	cookies := make(map[string]string)
	if r != nil && r.Response != nil {
		for _, cookie := range r.Cookies() {
			cookies[cookie.Name] = cookie.Value
		}
	}
	return cookies
}
