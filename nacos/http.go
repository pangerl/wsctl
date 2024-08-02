package nacos

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

func (d *Nacos) get(apiurl string) []byte {
	u, err := url.Parse(apiurl)
	if err != nil {
		panic(err)
	}
	if len(conf.Username) != 0 && len(conf.Password) != 0 {
		if len(u.RawQuery) == 0 {
			apiurl += "?accessToken=" + url.QueryEscape(d.Token)
		} else {
			apiurl += "&accessToken=" + url.QueryEscape(d.Token)
		}
	}
	req, _ := http.NewRequest("GET", apiurl, nil)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		if res.StatusCode == 403 {
			panic(fmt.Sprintf("%s请求状态码异常:%d 请使用--username --password参数进行鉴权", apiurl, res.StatusCode))
		}
		panic(fmt.Sprintf("%s请求状态码异常:%d", apiurl, res.StatusCode))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)
	resp, _ := ioutil.ReadAll(res.Body)
	return resp

}

func (d *Nacos) post(apiurl string, formData map[string]string) []byte {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	for key, val := range formData {
		_ = bodyWriter.WriteField(key, val)
	}
	contentType := bodyWriter.FormDataContentType()
	err := bodyWriter.Close()
	if err != nil {
		return nil
	}
	var req *http.Request
	u, err := url.Parse(apiurl)
	if u.Path == "/nacos/v1/auth/login" {
		req, _ = http.NewRequest("POST", apiurl, bodyBuf)
		req.Header.Set("Content-Type", contentType)
	}
	res, err := d.Client.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		if u.Path == "/nacos/v1/auth/login" && res.StatusCode == 403 {
			panic(fmt.Sprintf("%s请求状态码异常,认证失败!:%d", apiurl, res.StatusCode))
		}
		panic(fmt.Sprintf("%s请求状态码异常:%d", apiurl, res.StatusCode))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)
	resp, _ := ioutil.ReadAll(res.Body)
	return resp
}
