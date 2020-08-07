package utils

import (
	"bytes"
	"crypto/tls"
	"net/http"

	"github.com/sirupsen/logrus"
)

func HttpDo(method string, url string, body []byte) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	return client.Do(req)
}
