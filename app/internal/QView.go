package check

import (
	"fmt"
	"strings"
	"update/utils"
)

func QView(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVer, err := utils.GetExeVersion(execPath)
		if err != nil {
			return "", err
		}
		return strings.Join(strings.Split(currVer, ".")[:2], "."), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestTag, latestVer, err := utils.GetGitHubLatestTag("jurplel/qView")
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := utils.ToGitHubDirectURL(
		"qView/qView",
		latestTag,
		fmt.Sprintf("qView-%s-win64.exe", latestVer),
	)

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: fmt.Sprintf("https://github.com/qView/qView/releases/tag/%s", latestTag),
	}
}
