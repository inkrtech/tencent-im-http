/**
 * @Author: fuxiao
 * @Email: 576101059@qq.com
 * @Date: 2021/8/16 2:54 下午
 * @Desc: TODO
 */

package http_test

import (
	"errors"
	"testing"
	
	"github.com/webzh/http"
)

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlc2MiOjE2MjgwNDAzMjYxNTQ2MzIwMDAsImV4cCI6MTYyODIyMDMyNiwiaWF0IjoxNjI4MDQwMzI2LCJpZCI6MX0.KM19c6URIih-5SyycYIjNAdSiPKxMQEz3DoROm0N3nw"

func TestClient_Request(t *testing.T) {
	client := http.NewClient()
	client.SetBaseUrl("http://127.0.0.1:8199").Use(func(r *http.Request) (*http.Response, error) {
		return r.Next()
	}).Use(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("Invalid params.")
	})
	
	resp, err := client.Request(http.MethodGet, "/common/regions")
	if err != nil {
		t.Error(err)
		return
	}
	
	t.Log(resp.Response.Status)
}

func TestClient_Post(t *testing.T) {
	client := http.NewClient()
	client.SetBaseUrl("http://127.0.0.1:8199")
	client.SetBearerToken(token)
	client.SetContentType(http.ContentTypeJson)
	client.Use(func(r *http.Request) (*http.Response, error) {
		r.Request.Header.Set("Client-Type", "2")
		return r.Next()
	})
	
	type updateRegionArg struct {
		Id   int    `json:"id"`
		Pid  int    `json:"pid"`
		Code string `json:"code"`
		Name string `json:"name"`
		Sort int    `json:"sort"`
	}
	
	data := updateRegionArg{
		Id:   1,
		Pid:  0,
		Code: "110000",
		Name: "北京市",
		Sort: 0,
	}
	
	if resp, err := client.Put("/backend/region/update-region", data); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(resp.Response.Status)
		t.Log(resp.Response.Header)
		t.Log(resp.ReadBytes())
		t.Log(resp.ReadString())
		t.Log(resp.GetHeaders())
		t.Log(resp.GetCookies())
	}
}

func TestClient_Download(t *testing.T) {
	url := "https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png"
	
	if path, err := http.NewClient().Download(url, "./"); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(path)
	}
}
