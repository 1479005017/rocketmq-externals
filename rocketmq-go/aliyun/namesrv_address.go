package aliyun

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func FetchNamesrvAddress(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	resp.Body.Close()

	return strings.TrimSpace(string(body))
}
