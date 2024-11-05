package request

import (
	"bytes"
	"github.com/ysfgrl/gerror"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpService struct {
	BaseUrl string
	Token   string
}

func (base *HttpService) request(method string, path string, body []byte) ([]byte, *gerror.Error) {
	req, err := http.NewRequest(method, base.BaseUrl+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, gerror.GetError(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+base.Token)

	client := &http.Client{
		Timeout: time.Second * 5,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, gerror.GetError(err)
	}
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, gerror.GetError(err)
	}
	return resBytes, nil
}

func (base *HttpService) Get(path string) ([]byte, *gerror.Error) {
	return base.request("GET", path, nil)
}
func (base *HttpService) Put(path string, body []byte) ([]byte, *gerror.Error) {
	return base.request("PUT", path, body)
}

func (base *HttpService) Post(path string, body []byte) ([]byte, *gerror.Error) {
	return base.request("POST", path, body)
}
