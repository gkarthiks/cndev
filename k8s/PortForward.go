package k8s

import (
	"context"
	"errors"
	"fmt"
	"github.com/gkarthiks/cndev/prompts"
	"github.com/gkarthiks/cndev/utils"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type portForwardAPodRequest struct {
	// RestConfig is the kubernetes config
	RestConfig *rest.Config
	// Pod is the selected pod for this port forwarding
	Pod v1.Pod
	// LocalPort is the local port that will be selected to expose the PodPort
	LocalPort int
	// PodPort is the target port for the pod
	PodPort int
	// Steams configures where to write or read input from
	Streams genericclioptions.IOStreams
	// StopCh is the channel used to manage the port forward lifecycle
	StopCh <-chan struct{}
	// ReadyCh communicates when the tunnel is ready to receive traffic
	ReadyCh chan struct{}
}

// portForward will forwards the port from pod to local machine
func portForward(pod v1.Pod, userDefinedPort, podPort string) {
	logrus.Info("starting to execute the port forwarding")
	var wg sync.WaitGroup
	wg.Add(1)
	localPort, err := strconv.Atoi(userDefinedPort)
	containerPort, err := strconv.Atoi(podPort)
	if err != nil {
		logrus.Fatalf("error while converting the local port for port-forwarding: %v\n", err)
	}

	stopCh := make(chan struct{}, 1)
	readyCh := make(chan struct{})
	stream := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		logrus.Info("Closing the connection; shutting the port forwarding")
		close(stopCh)
		wg.Done()
	}()

	go func() {
		err := portForwardAPod(portForwardAPodRequest{
			RestConfig: Client.RestConfig,
			Pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      pod.Name,
					Namespace: pod.Namespace,
				},
			},
			LocalPort: localPort,
			PodPort:   containerPort,
			Streams:   stream,
			StopCh:    stopCh,
			ReadyCh:   readyCh,
		})
		if err != nil {
			logrus.Panicf("error while port-forwarding: %v", err)
		}
	}()

	select {
	case <-readyCh:
		open(fmt.Sprintf("http://localhost:%s/", userDefinedPort))
		break
	}
	logrus.Infof("Port forwarding is ready.")
	wg.Wait()
}

// portForwardAPod forwards the req to specified pod
func portForwardAPod(req portForwardAPodRequest) error {
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", req.Pod.Namespace, req.Pod.Name)
	hostIP := strings.TrimLeft(req.RestConfig.Host, "https://")

	transport, upgrader, err := spdy.RoundTripperFor(req.RestConfig)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", req.LocalPort, req.PodPort)}, req.StopCh, req.ReadyCh, nil, nil)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}

func Forward(appSpec utils.Apps, forwardPort string) {
	if len(forwardPort) < 1 {
		isDefaultPort := prompts.PromptYesNo(fmt.Sprintf("Do you want to continue with the default port: %s for the app %s?", appSpec.DefaultPort, appSpec.Name))
		if isDefaultPort == utils.StringNo {
			validate := func(input string) error {
				_, err := strconv.Atoi(input)
				if err != nil {
					return errors.New("provide a valid local port number")
				}
				return nil
			}
			forwardPort = prompts.PromptUser("What is the local port you want to use?", validate)
		} else {
			forwardPort = appSpec.DefaultPort
		}
	}

	podList, err := Client.Clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{
		LabelSelector: appSpec.Label,
	})
	logrus.Infof("starting to execute the port forwarding request on local port: %s", forwardPort)
	if err != nil {
		logrus.Fatalf("error while listing the pod matching to label: %s, err: %v \n", appSpec.Label, err)
	}
	if len(podList.Items) > 1 {
		logrus.Fatalf("found more than one pod for the matching label: %s\n", appSpec.Label)
	}

	if strings.ToLower(appSpec.Name) == strings.ToLower(utils.ArgoCD) {
		namespace, err := FindNamespace(map[string]string{"part-of": "argocd"})
		if err != nil {
			logrus.Fatalf("error while getting the argocd namespace for getting admin password: %v", err)
		}
		adminPassword, err := Client.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), utils.ArgoCDAdminSecretName, metav1.GetOptions{})
		if err != nil {
			logrus.Errorf("error while getting admin password from secret %s, error: %v\n", utils.ArgoCDAdminSecretName, err)
		} else {
			logrus.Infof(utils.Green("Username/Password for the ArgoCD login: admin/%s"), string(adminPassword.Data["password"]))
		}
	}
	portForward(podList.Items[0], forwardPort, appSpec.PodPort)
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
