package utils

import (
	"strings"
)

var (
	ArgoCD_URL            string
	DefaultGogsNs         string
	PFApps                []Apps
	ArgoCDAdminSecretName string
	AppVersion            string
)

func init() {
	setDefaults()
}

func setDefaults() {
	AppVersion = "v0.1.0"
	ArgoCD_URL = "https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml"
	DefaultGogsNs = "gogs"
	ArgoCDAdminSecretName = "argocd-initial-admin-secret"
	PFApps = []Apps{
		{
			Name:        "ArgoCD",
			DefaultNs:   "argocd",
			DefaultPort: ArgoPort,
			Label:       ArgoCDServerLabel,
			PodPort:     "8080",
		},
		{
			Name:        "Git",
			DefaultNs:   "gogs",
			DefaultPort: GogsPort,
			Label:       GogsLabel,
			PodPort:     "3000",
		},
	}
}

type Apps struct {
	Name        string
	DefaultNs   string
	DefaultPort string
	Label       string
	PodPort     string
}

// GetAppSpec returns true if the app present in the list and
// the App struct for the specified app from the list
func GetAppSpec(appName string) (bool, Apps) {
	for _, app := range PFApps {
		if strings.ToLower(app.Name) == strings.ToLower(appName) {
			return true, app
		}
	}
	return false, Apps{}
}
