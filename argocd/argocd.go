package argocd

import (
	"context"
	"github.com/gkarthiks/cndev/k8s"
	"github.com/gkarthiks/cndev/utils"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
)

// DeployArgoCD will create the argo namespace and installs the ArgoCD
func DeployArgoCD(argocdNamespace string) {
	logrus.Debugf("Starting the ArgoCD deployment on %s namespace.", argocdNamespace)
	createArgoCDNamespace(argocdNamespace)
	installOrDeleteArgoCD(argocdNamespace, true)
}

// DeleteArgoCD will delete all the argocd components
func DeleteArgoCD() {
	logrus.Debugf("Deleting the ArgoCD resources.")
	argocdNamespace, err := k8s.FindNamespace(map[string]string{utils.ManagedBy: utils.CNDev, utils.PartOf: utils.ArgoCD})
	if err != nil {
		logrus.Errorf(err.Error())
	} else {
		installOrDeleteArgoCD(argocdNamespace, false)
		k8s.DeleteNamespace(argocdNamespace)
	}
	logrus.Debugf("Deleting the Gogs resources.")

}

// createArgoCDNamespace creates the argocd namespace
func createArgoCDNamespace(argocdNamespace string) {
	logrus.Debugf("creating the argocd namespace %s", argocdNamespace)
	namespace := v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   argocdNamespace,
			Labels: map[string]string{utils.ManagedBy: utils.CNDev, utils.PartOf: utils.ArgoNamespace},
		},
	}
	ns, err := k8s.Client.Clientset.CoreV1().Namespaces().Create(context.Background(), &namespace, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("error while creating ns: %v\n", err)
	} else {
		logrus.Infof(utils.Green("created ns object: %v"), ns.Name)
	}
}

// installOrDeleteArgoCD will deploys or delete all argocd components based on the
// isInstall boolean value
func installOrDeleteArgoCD(argocdNamespace string, isInstall bool) {
	logrus.Debugf("executing GET to the %s to fetch argo manifests", utils.ArgoCD_URL)
	resp, err := http.Get(utils.ArgoCD_URL)
	if err != nil {
		logrus.Panicf("error occurred while obtaining the argo manifests: %v\n", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Panicf("couldn't read the data of argocd manifests: %v\n", err)
	}

	sepYamlFiles := strings.Split(string(data), "---")
	logrus.Debugf("total yaml files in argo manifest %v\n", len(sepYamlFiles))
	k8s.InstallOrDelete(sepYamlFiles, argocdNamespace, isInstall)
}

//func PortForward(argoPort string) {
//	logrus.Debugf("Find the argocd-server pod with the labels: %v for port forwarding on %s port.", utils.ArgoCDServerLabel, argoPort)
//	argocdServerPod, err := k8s.GetPod(utils.ArgoCDServerLabel)
//	if err != nil {
//		logrus.Fatalf("error while getting teh pod for matching labels %v to start port-forwarding: %v", utils.ArgoCDServerLabel, err)
//	} else {
//		logrus.Debugf("Starting the port forwarding for ArgoCD on %s port.", argoPort)
//		k8s.PortForward(argocdServerPod, argoPort)
//		logrus.Infof("Please open http://localhost:%s on your computer's browser", argoPort)
//	}
//}
