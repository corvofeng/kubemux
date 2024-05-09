package kubernetes

import (
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type KubernetesClient interface {
	ListClusters() []string
	DescribeCluster(clusterName string) Cluster
}

// Cluster is the representation of a K8S Cluster
// For now it is tailored to AWS, more specifically eks clusters
type Cluster struct {
	// Provider                 K8sProvider
	Name string

	// For AWS/Tencent Cloud, the region means the region of the cluster
	// But for BlueKing, the region means
	Region string

	ID                       string
	Endpoint                 string
	CertificateAuthorityData string
	Status                   string
	GenerateClusterConfig    func(cls *Cluster) *clientcmdapi.Cluster
	GenerateAuthInfo         func(cls *Cluster) *clientcmdapi.AuthInfo
}

func NewCluster() *Cluster {
	return &Cluster{
		GenerateClusterConfig: defaultGenerateClusterConfig,
	}
}

func defaultGenerateClusterConfig(cls *Cluster) *clientcmdapi.Cluster {
	cluster := clientcmdapi.NewCluster()
	cluster.Server = cls.Endpoint
	cluster.CertificateAuthorityData = []byte(cls.CertificateAuthorityData)
	return cluster
}
