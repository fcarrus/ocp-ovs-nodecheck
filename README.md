
ocp-ovs-nodecheck
===

A simple tool to check intra-node connection issues.

## How it works

This tool creates a Pod for each Node in your Cluster (DaemonSet).
Each Pod gets all its siblings and retrieve their IP addresses.
It then attempts to connect to port 8080 with an http GET.

If there are any issues in the Cluster's intra-node connectivity, you'll also see the issue in the logs.

## What you need

* Have an OpenShift cluster at hand.
* An OpenShift project and admin privileges on it
* The go-toolset-rhel7 ImageStream in the openshift project

## How to use it

* Import the go-toolset-rhel7 ImageStream

```
oc -n openshift import-image go-toolset-rhel7 --from=registry.redhat.io/devtools/go-toolset-rhel7 --confirm
```

* Load the template.yaml file 

```
oc create -f template.yaml
```

* Wait for the build to complete and wait until all Pods are running

* For each Pod, look at its logs, you should see something like:

