package check

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"update/utils"
)

var cjxlCurrVerRe = regexp.MustCompile(`cjxl v(\d+\.\d+\.\d+)`)

func LibJXL(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		cmdOutput, err := exec.Command(execPath, "--version").Output()
		if err != nil {
			return "", err
		}
		match := cjxlCurrVerRe.FindStringSubmatch(string(cmdOutput))
		if match == nil {
			return "", fmt.Errorf("unexpected command output")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestTag, latestVer, err := utils.GetGitHubLatestTag("libjxl/libjxl")
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := utils.ToGitHubDirectURL(
		"libjxl/libjxl",
		latestTag,
		"jxl-x64-windows-static.zip",
	)

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: fmt.Sprintf("https://github.com/libjxl/libjxl/releases/tag/%s", latestTag),
		Error:           nil,
	}
}
