// Package internal provides function for working with EKS cluseters
package kubernetes

import (
	"fmt"
	"math/rand"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lenID   = 10
)

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		// NOTE(mmicu): this is for testing only
		// #nosec
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int) string {
	return stringWithCharset(length, charset)
}

func dummyGenerateAuthInfo(_ *CPCluster) *clientcmdapi.AuthInfo {
	return clientcmdapi.NewAuthInfo()
}

func getMockClusters(i int, r string) *CPCluster {
	c := NewCluster()
	c.Name = fmt.Sprintf("clucster-name-%v-%v", r, i)
	c.Region = fmt.Sprintf("clucster-region-%v-%v", r, i)
	c.ID = fmt.Sprintf("clucster-id-%v-%v", r, i)
	c.Status = fmt.Sprintf("clucster-status-%v-%v", r, i)
	c.Endpoint = fmt.Sprintf("clucster-endpoint-%v-%v", r, i)
	c.CertificateAuthorityData = fmt.Sprintf("clucster-certificate-authority-data-%v-%v", r, i)
	c.GenerateClusterConfig = defaultGenerateClusterConfig
	c.GenerateAuthInfo = dummyGenerateAuthInfo
	c.Status = "ACTIVE"
	return c
}

func GetMockClusters(c int) []*CPCluster {
	d := make([]*CPCluster, 0, c)
	for i := 0; i < c; i++ {
		r := randomString(lenID)
		c := getMockClusters(i, r)
		d = append(d, c)
	}
	return d
}

func GetPredictableMockClusters(c int) []*CPCluster {
	d := make([]*CPCluster, 0, c)
	for i := 0; i < c; i++ {
		c := getMockClusters(i, "")
		d = append(d, c)
	}
	return d
}
