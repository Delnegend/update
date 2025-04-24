package check

import (
	"fmt"
	"os/exec"
	"strings"
	"update/utils"
)

func RsRPC(getExec func() (string, error)) CheckResult {

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
		currVerSlice := strings.SplitN(currVerString, " ", 2)
		if len(currVerSlice) == 2 {
			return strings.TrimSpace(currVerSlice[1]), nil
		}
		return "", fmt.Errorf("unexpected command output")
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	_, latestVer, err := utils.GetGitHubLatestTag("SpikeHD/rsRPC")
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DownloadPageURL: "https://github.com/SpikeHD/rsRPC/actions",
	}
}
