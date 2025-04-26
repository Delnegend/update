package check

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"update/utils"
)

var gitCurrVerRe = regexp.MustCompile(`git version ([\d\.]+)`)
var gitLatestVerRe = regexp.MustCompile(`version">\n.+?([\d\.]+)`)

func Git(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVerBytes, err := exec.Command(execPath, "--version").Output()
		if err != nil {
			return "", err
		}
		match := gitCurrVerRe.FindStringSubmatch(string(currVerBytes))
		if match == nil {
			return "", fmt.Errorf("unexpected command output")
		}
		return strings.TrimSuffix(strings.TrimSpace(match[1]), "."), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		content, err := utils.Fetch("https://git-scm.com/")
		if err != nil {
			return "", fmt.Errorf("can't fetch the homepage: %v", err)
		}
		match := gitLatestVerRe.FindStringSubmatch(string(content))
		if match == nil {
			return "", fmt.Errorf("unexpected website content")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := fmt.Sprintf("https://github.com/git-for-windows/git/releases/download/v%s.windows.1/Git-%s-64-bit.exe", latestVer, latestVer)

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: fmt.Sprintf("https://github.com/git-for-windows/git/releases/tag/v%s.windows.1", latestVer),
	}
}
