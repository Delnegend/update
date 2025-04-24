package check

import (
	"fmt"
	"regexp"
	"update/utils"
)

var nvidiaICATLatestVerRe = regexp.MustCompile(`https://icat-public-releases.s3.amazonaws.com/ICAT-([\d\.]+).exe`)
var nvidiaICATCurrVerRe = regexp.MustCompile(`(\d+\.\d+\.\d+)`)

func NvidiaICAT(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVer, err := utils.GetExeVersion(execPath)
		if err != nil {
			return "", fmt.Errorf("can't get executable version: %v", err)
		}
		match := nvidiaICATCurrVerRe.FindStringSubmatch(currVer)
		if match == nil {
			return "", fmt.Errorf("unexpected command output")
		}
		return match[1], nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		content, err := utils.Fetch("https://www.nvidia.com/en-us/geforce/technologies/icat/#")
		if err != nil {
			return "", err
		}
		match := nvidiaICATLatestVerRe.FindStringSubmatch(string(content))
		if match == nil {
			return "", err
		}
		return match[1], nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := fmt.Sprintf("https://icat-public-releases.s3.amazonaws.com/ICAT-%s.exe", latestVer)

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: "https://www.nvidia.com/en-us/geforce/technologies/icat/",
	}
}
