/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset

var namespace = "prova"

var config *rest.Config

var address, port string

func main() {
	// creates the in-cluster config
	var err error
	config, err = rest.InClusterConfig()
	var kcfg clientcmd.ClientConfig
	if err != nil {
		kcfg = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{},
		)

		config, err = kcfg.ClientConfig()
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Using external Config")
		namespace, _, err = kcfg.Namespace()

		// Listen only on localhost when not running in a Pod
		address = "localhost"

	} else {

		fmt.Println("Using in-cluster config")
		ns, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err != nil {
			fmt.Println("ERR: ", err)
		} else {
			namespace = string(ns)
		}

		// Listen on ANY when inside a Pod
		address = "0.0.0.0"

	}

	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Current namespace: ", namespace)

	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go webServer()
	go checkPods()

	fmt.Println("Starting main loop.")
	for {
		time.Sleep(10 * time.Second)
	}
}

func checkPods() {

	podselector := "name=ocp-ovs-nodecheck"
	fmt.Println("Looking for Pods with label ", podselector)

	for {
		pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{
			LabelSelector: podselector,
		})

		if err != nil {
			fmt.Println("ERR: ", err)
		} else {

			fmt.Printf("\nFound %d Pods in the namespace:\n\n", len(pods.Items))

			for _, pod := range pods.Items {

				if pod.Status.Phase == "Running" {
					url := fmt.Sprintf(
						"http://%v:%v/",
						pod.Status.PodIP,
						pod.Spec.Containers[0].Ports[0].ContainerPort)
					fmt.Printf("Pod: name:%v state:%v ip:%v, attempting to GET %s ...", pod.Name, pod.Status.Phase, pod.Status.PodIP, url)
					resp, err := http.Get(url)
					if err != nil {
						fmt.Println("ERR: ", err)
					} else {
						fmt.Println(resp.Status)
					}
					resp.Body.Close()
				} else {
					fmt.Printf("Pod: name:%v state:%v\n", pod.Name, pod.Status.Phase)
				}
				fmt.Println("")
			}

			if err != nil {
				panic(err.Error())
			}

		}

		time.Sleep(5 * time.Second)
	}

}

func webServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "It works.\n")
	})
	url := fmt.Sprintf("%s:%s", address, port)
	fmt.Printf("Starting web server on %s ...\n", url)
	http.ListenAndServe(url, nil)
}
