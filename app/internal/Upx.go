package check

import (
	"fmt"
	"os/exec"
	"strings"
	"update/utils"
)

func Upx(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVerBytes, err := exec.Command(execPath, "--version").Output()
		if err != nil {
			return "", err
		}
		currVerString := string(currVerBytes)
		for line := range strings.SplitSeq(currVerString, "\n") {
			if strings.HasPrefix(line, "upx ") {
				return strings.TrimSpace(strings.TrimPrefix(line, "upx ")), nil
			}
		}
		return "", fmt.Errorf("unexpected command output")
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestTag, latestVer, err := utils.GetGitHubLatestTag("upx/upx")
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := utils.ToGitHubDirectURL(
		"upx/upx",
		latestTag,
		"upx-5.0.0-win64.zip",
	)

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: fmt.Sprintf("https://github.com/upx/upx/releases/tag/%s", latestTag),
		Error:           nil,
	}
}
