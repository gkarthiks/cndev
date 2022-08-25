package gogs

import (
	"github.com/gkarthiks/cndev/assets"
	"github.com/gkarthiks/cndev/k8s"
	"github.com/gkarthiks/cndev/utils"
	"github.com/sirupsen/logrus"
	"strings"
)

// DeployGogs will deploy the gogs got server in gogs namespace
// from the assets/gogs.yaml file
func DeployGogs() {
	logrus.Debugln("Starting the Gogs git-server deployment.")
	installOrDeleteGogs(true)
}

// installOrDeleteGogs reads the manifest file data for gogs and installs
func installOrDeleteGogs(isInstall bool) {
	//path, _ := filepath.Abs(gogsFilePath)
	//data, err := ioutil.ReadFile(path)
	//if err != nil {
	//	logrus.Panicf("couldn't read the data of gogs manifests: %v\n", err)
	//}
	sepYamlFiles := strings.Split(assets.GogsStringFile, "---")
	logrus.Debugf("total yaml files in gogs manifest %v\n", len(sepYamlFiles))

	k8s.InstallOrDelete(sepYamlFiles, utils.DefaultGogsNs, isInstall)
}

// DeleteGogs deletes all the gogs components
func DeleteGogs() {
	logrus.Debugln("Deleting the Gogs git-server deployment.")
	installOrDeleteGogs(false)
}

func PortForwardSvc(gogsPort string) {

}
