package kubernetes

import (
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type EnumCloudProvider int

const (
	None EnumCloudProvider = iota
	AWS
	TencentCloud
	Google
	Azure
	DigitalOcean
)

// Cloud Provider interface
// The cloud provider like AWS, Google, Azure, etc. should implement this interface
// They should have their own region and cluster management
type CloudProvider interface {
	Init() error
	ListRegions() ([]string, error)

	// ListClusters returns a list of clusters in the given regions
	// It also takes a function to set the progress of the operation
	ListClusters(regions []string, setProgress func(int)) ([]*CPCluster, error)

	// Get the cluster config
	GetClusterConfig(cluster *CPCluster) (*clientcmdapi.Config, error)

	// VerifyCluster(cluster CPCluster) (bool, string)
}

// CPCluster is the representation of a K8S Cloud Provider Cluster
// For now it is tailored to AWS, more specifically eks clusters
type CPCluster struct {
	Provider EnumCloudProvider
	Name     string

	// For AWS/Tencent Cloud, the region means the region of the cluster
	// But for BlueKing, the region means
	Region string

	ID                       string
	Endpoint                 string
	CertificateAuthorityData string
	Status                   string

	GenerateClusterConfig func(cls *CPCluster) *clientcmdapi.Cluster
	GenerateAuthInfo      func(cls *CPCluster) *clientcmdapi.AuthInfo
}

func NewCluster() *CPCluster {
	return &CPCluster{
		GenerateClusterConfig: defaultGenerateClusterConfig,
	}
}

func defaultGenerateClusterConfig(cls *CPCluster) *clientcmdapi.Cluster {
	cluster := clientcmdapi.NewCluster()
	cluster.Server = cls.Endpoint
	cluster.CertificateAuthorityData = []byte(cls.CertificateAuthorityData)
	return cluster
}

func (cls *CPCluster) GetConfigAuthInfo() *clientcmdapi.AuthInfo {
	return cls.GenerateAuthInfo(cls)
}
