package check

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"update/utils"
)

var html2markdownCurrVer = regexp.MustCompile(`GitVersion: +?([\d\.]+)`)

func HTML2Markdown(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVerBytes, err := exec.Command(execPath, "--version").Output()
		if err != nil {
			return "", err
		}
		match := html2markdownCurrVer.FindStringSubmatch(string(currVerBytes))
		if match == nil {
			return "", fmt.Errorf("unexpected command output")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestTag, latestVer, err := utils.GetGitHubLatestTag("JohannesKaufmann/html-to-markdown")
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := utils.ToGitHubDirectURL(
		"JohannesKaufmann/html-to-markdown",
		latestTag,
		"https://github.com/JohannesKaufmann/html-to-markdown/releases/download/v2.3.1/html-to-markdown_Windows_x86_64.zip",
	)

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: fmt.Sprintf("https://github.com/JohannesKaufmann/html-to-markdown/releases/tag/%s", latestTag),
		Error:           nil,
	}
}
