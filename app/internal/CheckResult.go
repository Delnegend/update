package check

type CheckResult struct {
	CurrentVersion  string
	LatestVersion   string
	DirectURL       string
	DirectURLAlive  bool
	DownloadPageURL string
	Error           error
}
