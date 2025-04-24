package check

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"update/utils"
)

var pythonCurrVerRe = regexp.MustCompile(`Python (\d+\.\d+\.\d+)`)
var pythonLatestVerRe = regexp.MustCompile(`python-(\d+\.\d+\.\d+)-amd64.exe`)

func Python(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVerBytes, err := exec.Command(execPath, "--version").Output()
		if err != nil {
			return "", err
		}
		match := pythonCurrVerRe.FindStringSubmatch(string(currVerBytes))
		if match == nil {
			return "", fmt.Errorf("unexpected command output")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		content, err := utils.Fetch("https://www.python.org/downloads/")
		if err != nil {
			return "", fmt.Errorf("can't fetch the homepage: %v", err)
		}
		match := pythonLatestVerRe.FindStringSubmatch(string(content))
		if match == nil {
			return "", fmt.Errorf("unexpected website content")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := fmt.Sprintf("https://www.python.org/ftp/python/%s/python-%s-amd64.exe", latestVer, latestVer)

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: "https://www.python.org/downloads/",
	}
}
