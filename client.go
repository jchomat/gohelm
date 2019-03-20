/*
Package gohelm implements a Helm client.

  package main

  import (
  	"log"

  	"gitlab.com/pablo-ruth/gohelm"
  	"gitlab.com/pablo-ruth/gok8s"
  )

  func main() {
  	client, config, err := gok8s.NewClient("https://api.mycluster.com", "mytoken", []byte(""), true)
  	if err != nil {
  		log.Fatal(err)
  	}

  	helm, err := gohelm.NewClient("v2.9.0", client, config)
  	if err != nil {
  		log.Fatal(err)
  	}
  	helm.Close()
  }

*/
package gohelm

import (
	"errors"
	"fmt"

	"github.com/pablo-ruth/gok8s"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Represents an Helm client
type Client struct {
	Version string
	Tunnel  *gok8s.Tunnel
	Conn    *grpc.ClientConn
}

// NewClient initializes a new Helm client
func NewClient(version, namespace string, k8sclient *kubernetes.Clientset, k8sconfig *rest.Config) (*Client, error) {
	client := &Client{Version: version}

	// Find Tiller pod
	podName, err := getTillerPodName(k8sclient, namespace)
	if err != nil {
		return nil, err
	}

	// Open tunnel to Tiller(helm)
	tunnel := gok8s.NewTunnel(k8sclient.CoreV1().RESTClient(), k8sconfig, namespace, podName, 44134)
	err = tunnel.Open()
	if err != nil {
		return nil, err
	}
	client.Tunnel = tunnel

	// Init GRPC client to Tiller
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", tunnel.Local), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client.Conn = conn

	return client, nil
}

// NewClientFromCluster initializes a new Helm client from inside
func NewClientFromCluster(version, namespace string, k8sclient *kubernetes.Clientset, k8sconfig *rest.Config) (*Client, error) {
	client := &Client{Version: version}

	// Find Tiller pod
	podIP, err := getTillerPodName(k8sclient, namespace)
	if err != nil {
		return nil, err
	}

	// Init GRPC client to Tiller
	conn, err := grpc.Dial(fmt.Sprintf("%s:44134", podIP), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client.Conn = conn

	return client, nil
}

// Close connection to helm
func (c *Client) Close() error {
	err := c.Conn.Close()
	c.Tunnel.Close()

	return err
}

// getTillerPodName returns the name of the first Tiller pod in this cluster/namespace
func getTillerPodName(k8sclient *kubernetes.Clientset, namespace string) (string, error) {
	listOpts := metav1.ListOptions{
		LabelSelector: "app=helm,name=tiller",
	}

	pods, err := k8sclient.CoreV1().Pods(namespace).List(listOpts)
	if err != nil {
		return "", err
	}

	if len(pods.Items) == 0 {
		return "", errors.New("No tiller pod found")
	}

	return pods.Items[0].GetName(), nil
}

// getTillerPodIP returns the IP of the first Tiller pod in this cluster/namespace
func getTillerPodIP(k8sclient *kubernetes.Clientset, namespace string) (string, error) {
	listOpts := metav1.ListOptions{
		LabelSelector: "app=helm,name=tiller",
	}

	pods, err := k8sclient.CoreV1().Pods(namespace).List(listOpts)
	if err != nil {
		return "", err
	}

	if len(pods.Items) == 0 {
		return "", errors.New("No tiller pod found")
	}

	return pods.Items[0].Status.PodIP, nil
}
