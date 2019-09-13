
ocp-ovs-nodecheck
===

A simple tool to test for intra-node connection issues.

## How it works

This tool creates a Pod for each Node in your Cluster (DaemonSet).

Each Pod tries to connect to all other Pods of the DaemonSet on port 8080 with an http GET.

If there are any issues (like [this one](https://access.redhat.com/solutions/3083121)) in the Cluster's intra-node connectivity, you'll likely see it in the logs.


## What you need

* Have an OpenShift cluster at hand.
* An OpenShift project and admin privileges on it
* The go-toolset-rhel7 ImageStream in the openshift project

## How to use it

* Import the go-toolset-rhel7 ImageStream

```
oc -n openshift import-image go-toolset-rhel7 --from=registry.redhat.io/devtools/go-toolset-rhel7 --confirm
```

* Create your project

```
oc new-project myproject
```

* Load the [template.yml](template.yml) file 

```
oc create -f template.yml
```

* Give the project's default serviceaccount the rights to read the Pods' info

```
oc adm policy add-role-to-user view -z default
```

* Wait for the build to complete and until all the Pods are running.

* For each Pod, look at its logs

```
oc logs -f pods/ocp-ovs-nodecheck-gbx8m
```

You should see something like:

```
Found 5 Pods in the namespace:
Pod: name:ocp-ovs-nodecheck-4d55c state:Running ip:10.130.1.16, attempting to GET http://10.130.1.16:8080/ ...200 OK
Pod: name:ocp-ovs-nodecheck-5xv6m state:Running ip:10.128.3.65, attempting to GET http://10.128.3.65:8080/ ...200 OK
Pod: name:ocp-ovs-nodecheck-q7wsk state:Running ip:10.130.2.106, attempting to GET http://10.130.2.106:8080/ ...200 OK
Pod: name:ocp-ovs-nodecheck-tx4kl state:Running ip:10.129.2.58, attempting to GET http://10.129.2.58:8080/ ...200 OK
Pod: name:ocp-ovs-nodecheck-wwg7v state:Running ip:10.129.0.248, attempting to GET http://10.129.0.248:8080/ ...200 OK
```

And, in case something's wrong:

```
Pod: name:ocp-ovs-nodecheck-tx4kl state:Running ip:10.129.2.58, attempting to GET http://10.129.2.58:8080/ ...
  ERR:  Get http://10.129.2.58:8080/: dial tcp 10.129.2.58:8080: i/o timeout
```

YMMV of course, depending on the issue.

