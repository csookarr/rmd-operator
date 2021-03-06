package rmd

import (
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes/scheme"
)

var rmdPodPath = "/rmd-manifests/rmd-pod.yaml"

// NewClient creates a new client to RMD for each controller
func NewClient() OperatorRmdClient {
	logger := log.WithName("NewClient")

	tlsEnabled, err := isTLSEnabled(rmdPodPath)
	if err != nil {
		logger.Info("error occurred checking TLS enablement, reverting to default client", "err", err)
		return NewDefaultOperatorRmdClient()
	}
	if tlsEnabled {
		logger.Info("creating TLS client for operator controller")
		operatorRmdClient, err := NewOperatorRmdClient()
		if err != nil {
			logger.Info("error occurred creating TLS client, reverting to default client", "err", err)
			return NewDefaultOperatorRmdClient()
		}
		logger.Info("TLS client created successfully")
		return operatorRmdClient
	}
	logger.Info("returning default client (no TLS)")
	return NewDefaultOperatorRmdClient()
}

func isTLSEnabled(path string) (bool, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(yamlFile, nil, nil)
	if err != nil {
		return false, err
	}

	rmdPod := obj.(*corev1.Pod)
	err = errors.NewServiceUnavailable("container not found in rmd pod manifest")
	if len(rmdPod.Spec.Containers) == 0 {
		return false, err
	}
	if len(rmdPod.Spec.Containers[0].Ports) == 0 {
		return false, err
	}

	// Check if container port is set to TLS enabled port 8443
	if rmdPod.Spec.Containers[0].Ports[0].ContainerPort == 8443 {
		return true, nil
	}

	return false, nil
}
