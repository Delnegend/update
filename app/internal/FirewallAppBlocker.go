package check

import (
	"fmt"
	"regexp"
	"update/utils"
)

var fabLatestVerRe = regexp.MustCompile(`title="Firewall App Blocker \(Fab\) v(\d+\.\d+)"`)

func FirewallAppBlocker(getExec func() (string, error)) CheckResult {
	webpageURL := "https://www.sordum.org/?s=firewall+app+blocker"

	latestVer, err := func() (string, error) {
		content, err := utils.Fetch(webpageURL)
		if err != nil {
			return "", fmt.Errorf("can't fetch webpage: %v", err)
		}
		if !fabLatestVerRe.MatchString(string(content)) {
			return "", fmt.Errorf("unexpected website content")
		}
		return fabLatestVerRe.FindStringSubmatch(string(content))[1] + ".0.0", nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	currVer, err := utils.GetExeVersion(`C:\Portable\Firewall App Blocker\Fab_x64.exe`)
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: can't get executable version: %v", err)}
	}

	directURL := "https://www.sordum.org/files/download/firewall-app-blocker/fab.zip"

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: webpageURL,
	}
}
