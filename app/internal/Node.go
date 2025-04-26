package check

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"update/utils"
)

var nodeLTSLatestVerRe = regexp.MustCompile(`nodejs\.org\/dist\/v([\d\.]+)`)

func Node(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		cmdOutputBytes, err := exec.Command(execPath, "--version").Output()
		if err != nil {
			return "", fmt.Errorf("can't run node --version: %v", err)
		}
		return strings.TrimPrefix(strings.TrimSpace(string(cmdOutputBytes)), "v"), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		content, err := utils.Fetch("https://nodejs.org/en")
		if err != nil {
			return "", fmt.Errorf("can't fetch the homepage: %v", err)
		}
		match := nodeLTSLatestVerRe.FindStringSubmatch(string(content))
		if match == nil {
			return "", fmt.Errorf("unexpected website content")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DownloadPageURL: "https://nodejs.org/",
	}
}
