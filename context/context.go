package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func (c *Context) ReadJson(req interface{}) error {
	// 帮我读 body
	// 帮我反序列化
	body, err := io.ReadAll(c.R.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, req); err != nil {
		return err
	}
	return nil
}

func (c *Context) WriteJson(code int, resp interface{}) error {
	// 帮我序列化json
	c.W.WriteHeader(code)
	respJson, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = c.W.Write(respJson)
	return err
}
