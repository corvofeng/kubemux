package common

import (
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type ClusterCommon interface {
	GetClusterName() string
	GetClusterConfig(loc string) *clientcmdapi.Cluster
}
