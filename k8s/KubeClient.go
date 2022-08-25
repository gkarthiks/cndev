package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gkarthiks/cndev/utils"
	discoveryK8s "github.com/gkarthiks/k8s-discovery"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

var (
	Client     *discoveryK8s.K8s
	K8sVersion string
)

func init() {
	Client, _ = discoveryK8s.NewK8s()
	K8sVersion, _ = Client.GetVersion()
	//fmt.Printf("Version of running Kubernetes: %s\n", version)
}

// GetRestMapperAndDynamicClient creates the Rest Mapper and Dynamic Client Interface for the
// given cluster defined in the kubeconfig or from within the cluster
// and while deleting it just deletes resource-by-resource
func GetRestMapperAndDynamicClient() (*restmapper.DeferredDiscoveryRESTMapper, dynamic.Interface) {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(Client.RestConfig)
	if err != nil {
		logrus.Panicf("error while instantiating the discovery client: %v", err)
	}
	restMapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))

	dynamicClient, err := dynamic.NewForConfig(Client.RestConfig)
	if err != nil {
		logrus.Panicf("error while creating dynamic client: %v", err)
	}
	return restMapper, dynamicClient
}

// serverSidePatch creates or updates the object via SSA mechanism using the ApplyPatchType and
// dynamic resource interface. This in-turn makes the code to work with any version of the kubernetes.
func serverSidePatch(gvk *schema.GroupVersionKind, unstructuredObj *unstructured.Unstructured, userDefinedNamespace string, isInstall bool) {
	restMapper, dynamicClient := GetRestMapperAndDynamicClient()
	mapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		logrus.Fatalf("error while finding GVR: %v", err)
	}
	var dynamicResourceInterface dynamic.ResourceInterface

	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources
		var namespace string
		if unstructuredObj.GetNamespace() == "" {
			namespace = userDefinedNamespace
		} else {
			namespace = unstructuredObj.GetNamespace()
		}
		dynamicResourceInterface = dynamicClient.Resource(mapping.Resource).Namespace(namespace)
	} else {
		// cluster-scoped resources
		dynamicResourceInterface = dynamicClient.Resource(mapping.Resource)
	}

	data, err := json.Marshal(unstructuredObj)
	if err != nil {
		logrus.Fatalf(utils.Red("error while marshaling object into JSON: %v"), err)
	}

	if isInstall {
		patchedObj, err := dynamicResourceInterface.Patch(context.Background(), unstructuredObj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
			FieldManager: utils.CNDev,
		})
		if err != nil {
			logrus.Fatalf(utils.Red("error while applying the %s object: %v"), unstructuredObj.GetName(), err)
		} else {
			logrus.Infof(utils.Magenta("applied object: %v \n"), patchedObj.GetName())
		}
	} else {
		gracePeriodSeconds := int64(0)
		deleteBackground := metav1.DeletePropagationBackground
		err := dynamicResourceInterface.Delete(context.Background(), unstructuredObj.GetName(), metav1.DeleteOptions{
			GracePeriodSeconds: &gracePeriodSeconds,
			PropagationPolicy:  &deleteBackground,
		})
		if err != nil {
			logrus.Fatalf("error while deleting the %s object: %v", unstructuredObj.GetName(), err)
		} else {
			logrus.Infof(utils.Magenta("deleted object: API=%v, Name: %v \n"), unstructuredObj.GetObjectKind().GroupVersionKind(), unstructuredObj.GetName())
		}
	}
}

// InstallOrDelete installs or deletes the unstructured data according
// to the isInstall boolean
func InstallOrDelete(sepYamlFiles []string, userDefinedNamespace string, isInstall bool) {
	logrus.Debugf("starting to execute install or delete based on isInstall: %v", isInstall)
	var decodedUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	for _, file := range sepYamlFiles {
		unstructuredObj := &unstructured.Unstructured{}
		_, gvk, err := decodedUnstructured.Decode([]byte(file), nil, unstructuredObj)
		if err != nil {
			logrus.Fatalf(utils.Red("error while decoding YAML manifest into unstructured.Unstructured: %v"), err)
		} else {
			logrus.Debugf("GVK: %v\n", gvk)
		}
		serverSidePatch(gvk, unstructuredObj, userDefinedNamespace, isInstall)
	}
}

// FindNamespace returns the namespace for matching labels
func FindNamespace(labelSelector map[string]string) (namespace string, err error) {
	logrus.Debugf("finding the namespace with matching lables %v\n", labelSelector)
	ns, err := Client.Clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelSelector).String(),
	})
	if err != nil {
		return "", fmt.Errorf("error while listing the namespaces matching the labels: %v; err: %v", labelSelector, err)
	} else if len(ns.Items) > 1 {
		return "", fmt.Errorf("more than 1 namespace found for the matching labels: %v \n List of namespaces found: %v", labelSelector, ns.Items)
	} else {
		return ns.Items[0].Name, nil
	}
}

// DeleteNamespace deletes all the namespaces in the v1 namespacelist
func DeleteNamespace(namespace string) {
	logrus.Debugf("deleting the namespace: %s", namespace)
	gracePeriodSeconds := int64(0)
	deleteBackground := metav1.DeletePropagationBackground
	err := Client.Clientset.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
		PropagationPolicy:  &deleteBackground,
	})
	if err != nil {
		logrus.Errorf("error while deleting ns: %v\n", err)
	} else {
		logrus.Infof("Namespace %s deleted\n", namespace)
	}
}
