package kubernetes

import (
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Cloud Provider interface
// The cloud provider like AWS, Google, Azure, etc. should implement this interface
// They should have their own region and cluster management
type CloudProvider interface {
	Init() error
	ListRegions() ([]string, error)

	// List clusters in all regions
	ListClusters() ([]Cluster, error)

	// Get the cluster config
	GetClusterConfig(cluster Cluster) error
	VerifyCluster(cluster Cluster) (bool, string)
}

type K8sProvider int

const (
	None K8sProvider = iota
	AWS
	TencentCloud
	Google
	Azure
	DigitalOcean
)

// Cluster is the representation of a K8S Cluster
// For now it is tailored to AWS, more specifically eks clusters
type Cluster struct {
	Provider K8sProvider
	Name     string

	// For AWS/Tencent Cloud, the region means the region of the cluster
	// But for BlueKing, the region means
	Region string

	ID                       string
	Endpoint                 string
	CertificateAuthorityData string
	Status                   string

	GenerateClusterConfig func(cls *Cluster) *clientcmdapi.Cluster
	GenerateAuthInfo      func(cls *Cluster) *clientcmdapi.AuthInfo
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

//	func (cls *Cluster) GetConfigAuthInfo() *clientcmdapi.AuthInfo {
//		return cls.GenerateAuthInfo(cls)
//	}
