package check

import (
	"fmt"
	"os/exec"
	"strings"
	"update/utils"
)

func Onefetch(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVerBytes, err := exec.Command(execPath, "--version").Output()
		if err != nil {
			return "", err
		}
		currVerSlice := strings.SplitN(string(currVerBytes), " ", 2)
		if len(currVerSlice) == 2 {
			return strings.TrimSpace(currVerSlice[1]), nil
		}
		return "", fmt.Errorf("unexpected command output")
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestTag, latestVer, err := utils.GetGitHubLatestTag("o2sh/onefetch")
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := utils.ToGitHubDirectURL("o2sh/onefetch", latestTag, "onefetch-win.tar.gz")
	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: fmt.Sprintf("https://github.com/o2sh/onefetch/releases/tag/%s", latestTag),
		Error:           nil,
	}
}
