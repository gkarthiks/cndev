# CN Dev

```text
                  __         
  _________  ____/ /__ _   __
 / ___/ __ \/ __  / _ \ | / /
/ /__/ / / / /_/ /  __/ |/ / 
\___/_/ /_/\__,_/\___/|___/  
```

*cndev* is a command line utility for the Container Native Developers, who develop applications and microservices and intend to deploy them in the Kubernetes Cluster. A lot of processes are involved between the developer coding the microservice, and it gets deployed onto a cluster in targeted environment.

Tools like ArgoCD is used as GitOps engine for the _continuous deployment_(CD) from source code control system to the environments. But these would require a configuration to make sure the service is deployed correctly. And this configuration is done as a manifest in the git repo again; which would be a centralized repository for the entire developer community. Doing the configuration for the first time or need to redo the configuration for every minimal change is going to need a lot of time-consuming build, approval process for the PR etc.

To reduce these time-consuming process, *cndev* mimics the CD environment into your private environment. You will be bootstrapping the _Git_ server and _ArgoCD_ components in your private cluster.

## Prerequisites
There are certain prerequisites for the _cndev_ to work as expected. Feel free to raise a PR or an issue to make it more usable.
* A locally provisioned Kubernetes Cluster
* A valid *kubeconfig* pointing to a valid Kubernetes cluster
* Admin access to the cluster is required
* `argocd` cli tool and basic knowledge of how to use it

## How does it work?
After having a private cluster or a valid _kubeconfig_ file that points to the Kubernetes cluster that you have admin access to, the next step is to provision the CD infra.

### Get Going
1. Make a directory or clone a git repository in you local machine.
6. cd into  the directory.
7. If it's a new directory, code the app and create the deployment manifests.
8. Execute `cndev init` to provision the _Git_ server and _ArgoCD_ server.
9. Run `cndev pfa -f git` and in another terminal `cndev pfa -f argocd`. This will open up the Git UI and ArgoCD UI in your default browser.
10. A one time initial setup is necessary for Git Server.
    1. Choose `SQLite3` for the database.
    1. Change the port number on `Application URL` text box to the port number on your address bar and click `Install Gogs`.
    1. Now click `Sign up now` link and create a user for yourself and login. That's your local git server.
11. Create a repository in the private git and push your changes from local machine.
12. Once the app pushed into the Git repo, you can stop the port forwarding executed by `cndev pfa -f git` command.
13. Head over to the terminal where you executed the `cndev pfa -f argocd` command.
    1. Copy the username and password from the terminal.
    1. Use that to login to ArgoCD UI.
    1. Create a new application by clicking `+ NEW APP`.
    1. Enter the `Application Name` as your app name.
    1. Choose `default` for `Project`.
    1. For the `Repository URL`, copy the `http` url from your local git and replace `http://localhost:<port>` to `http://gogs-svc.gogs.svc:18080`.
    1. Choose the path for the deployment manifests.
    1. Choose the default `https://kubernetes.default.svc` as the `Cluster URL`.
    1. Desired namespace to deploy, and click `Create`.
14. Now for every change in the manifests, you can sync them from your local machine by running the following command. `argocd app sync <Application Name> --local <Path to manifests>`.

<hr/>

## Legends
Private (Cluster/Environment): A local environment or a kubernetes cluster that is running in your local machine.
