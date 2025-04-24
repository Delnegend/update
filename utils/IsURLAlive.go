package utils

import (
	"io"
	"net/http"
)

func IsURLAlive(url string, skipCheck bool) bool {
	if skipCheck {
		return false
	}
	if url == "" {
		return false
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Range", "bytes=0-0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode <= 399
}
