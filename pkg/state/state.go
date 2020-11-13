package state

import (
	"context"
	"fmt"
	"os"

	intelv1alpha1 "github.com/intel/rmd-operator/pkg/apis/intel/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

//var rmdPodPath = "/rmd-manifests/rmd-pod.yaml"
var log = logf.Log.WithName("state")

type RmdNodeData struct {
	stateList map[string]string
}

var rmdStateList = map[string]string{}

// NewRmdNodeData() creates an empty RmdNodeData object if no nodestates are present
// Otherwise, relevant data is extracted from node states and placed in a map
func NewRmdNodeData() RmdNodeData {
	// Create client
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		log.Error(err, "failed to create client")
		os.Exit(1)
	}

	// List RmdNodeStates
	rmdNodeStates := &intelv1alpha1.RmdNodeStateList{}
	err = cl.List(context.TODO(), rmdNodeStates)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("No RMD Node States present")
			log.Info("WE DID IT, Return empty RmdNodeData")
			return RmdNodeData{}
		}
		log.Error(err, "failed to list RMD Node States")
		log.Info("WE DID not DO IT")
		return RmdNodeData{}
	}

	// Assign relevant node state information into map
	for _, rmdNodeState := range rmdNodeStates.Items {
		rmdStateList[rmdNodeState.Spec.Node] = rmdNodeState.GetObjectMeta().GetNamespace()
	}
	log.Info("WE DID IT 2.0")
	log.Info("RETURN NON-ZERO RmdNodeData")
	return RmdNodeData{
		stateList: rmdStateList,
	}
}

/*
// TODO: create GetRmdNodeStateList()
func GetRmdNodeData() RmdNodeData {
	return RmdNodeData{
		stateList: rmdStateList,
	}
}*/

// TODO: create UpdateRmdNodeStateList()
func UpdateRmdNodeData(nodeName string, namespaceName string) {
	log.Info(fmt.Sprintf("States before update: %v", rmdStateList))
	// Check if list needs to be updated
	// if the key-value pair already exists
	rmdStateList[nodeName] = namespaceName
	log.Info(fmt.Sprintf("Updated States      : %v", rmdStateList))
	log.Info("Success! NODE DATA UPDATED")
}

// deleteRmdNodeStateData
func DeleteRmdNodeData(nodeName string, namespaceName string) {
	log.Info(fmt.Sprintf("States before delete: %v", rmdStateList))
	delete(rmdStateList, nodeName)
	log.Info("DATA SUCCESSFULLY DELETED!")
	log.Info(fmt.Sprintf("After Delete: %v", rmdStateList))
}
