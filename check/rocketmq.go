// Package check @Author lanpang
// @Date 2024/9/10 下午6:13:00
// @Desc
package check

import (
	"io"
	"log"
	"net/http"
)

func doRequest(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("E! fail to close the res", err)
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Println("E! fail to read request data", err)
		return nil
	} else {
		return body
	}
}
