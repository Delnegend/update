package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func GetGitHubLatestTag(repo string) (string, string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url := fmt.Sprintf("https://github.com/%s/releases/latest", repo)
	resp, err := client.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	loc := resp.Header.Get("Location")
	if loc == "" {
		return "", "", errors.New("no redirect location found")
	}
	parts := strings.Split(loc, "/")

	tag := parts[len(parts)-1]

	return tag, strings.TrimPrefix(tag, "v"), nil
}
