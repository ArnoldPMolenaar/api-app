package requests

type AreAppsAvailable struct {
	AppNames []string `query:"appName"`
}
