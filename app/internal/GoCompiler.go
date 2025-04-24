package check

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"update/utils"
)

var (
	goCompilerCurrVerRe   = regexp.MustCompile(`go version go(\d+\.\d+\.\d+)`)
	goCompilerLatestVerRe = regexp.MustCompile(`go(\d+\.\d+\.\d+)\.windows-amd64\.msi`)
)

func GoCompiler(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVerBytes, err := exec.Command(execPath, "version").Output()
		if err != nil {
			return "", err
		}
		match := goCompilerCurrVerRe.FindStringSubmatch(string(currVerBytes))
		if match == nil {
			return "", fmt.Errorf("unexpected command output")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		content, err := utils.Fetch("https://go.dev/dl/")
		if err != nil {
			return "", fmt.Errorf("can't fetch the homepage: %v", err)
		}
		match := goCompilerLatestVerRe.FindStringSubmatch(string(content))
		if match == nil {
			return "", fmt.Errorf("unexpected website content")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := fmt.Sprintf("https://go.dev/dl/go%s.windows-amd64.msi", latestVer)

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: "https://go.dev/dl/",
	}

}
