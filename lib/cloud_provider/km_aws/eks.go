// Package internal provides wrapper for creating aws sessions
package km_aws

// Some code copy from: https://github.com/mateimicu/kdiscover

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	cluster "kubemux/lib/cloud_provider"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type AWSProvider struct {
	EKS eksiface.EKSAPI
}

type eksClient struct {
	EKS    eksiface.EKSAPI
	Region string
}

func (c *eksClient) String() string {
	return fmt.Sprintf("EKS Client for region %v", c.Region)
}

func (c *eksClient) GetClusters(ch chan<- *cluster.CPCluster) {
	input := &eks.ListClustersInput{}
	defer close(ch)

	err := c.EKS.ListClustersPages(input,
		func(page *eks.ListClustersOutput, lastPage bool) bool {
			log.WithFields(log.Fields{
				"svc":  c.String(),
				"page": page.GoString(),
			}).Debug("Parse page")
			for _, cluster := range page.Clusters {
				log.WithFields(log.Fields{
					"svc":     c.String(),
					"cluster": *cluster,
				}).Debug("Found cluster")
				if cls, err := c.detailCluster(*cluster); err == nil {
					ch <- cls
				} else {
					log.WithFields(log.Fields{
						"svc":     c.String(),
						"cluster": *cluster,
						"err":     err,
					}).Warn("Can't get details on the cluster")
				}
			}

			if lastPage {
				log.WithFields(log.Fields{
					"svc": c.String(),
				}).Debug("hit last page")
			}
			return !lastPage
		})

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"svc": c.String(),
		}).Warn("Can't list clusters")
	}
}

func (c *eksClient) detailCluster(cName string) (*cluster.CPCluster, error) {
	input := &eks.DescribeClusterInput{
		Name: aws.String(cName),
	}
	result, err := c.EKS.DescribeCluster(input)

	if err != nil {
		// TODO(mmicu): handle errors better here
		if aerr, ok := err.(awserr.Error); ok {
			log.Warn(aerr.Error())
		} else {
			log.Warn(err.Error())
		}
		msg := fmt.Sprintf("Can't fetch more details for the cluster %v", cName)
		log.WithFields(log.Fields{
			"cluster-name": cName,
			"svc":          c.String(),
		}).Warn(msg)
		return nil, errors.New(msg)
	}
	if result.Cluster.CertificateAuthority.Data == nil {
		return nil, errors.New("certificate authority is nil")
	}

	certificatAuthorityData, err := base64.StdEncoding.DecodeString(*result.Cluster.CertificateAuthority.Data)
	if err != nil {
		log.WithFields(log.Fields{
			"cluster-name":               *result.Cluster.Name,
			"arn":                        *result.Cluster.Arn,
			"certificate-authority-data": *result.Cluster.CertificateAuthority.Data,
			"svc":                        c.String(),
		}).Error("Can't decode the Certificate Authority Data")
		return nil, err
	}

	cls := cluster.NewCluster()
	cls.Name = *result.Cluster.Name
	cls.ID = *result.Cluster.Arn
	cls.Endpoint = *result.Cluster.Endpoint
	cls.CertificateAuthorityData = string(certificatAuthorityData)
	cls.Status = *result.Cluster.Status
	cls.Region = c.Region

	return cls, nil
}

func NewEKS(region string) (*eksClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.WithFields(log.Fields{
			"region": region,
			"error":  err.Error(),
		}).Error("Failed to create AWS SDK session")
		return nil, err
	}
	return &eksClient{
		EKS:    eks.New(sess),
		Region: region,
	}, nil
}

func (c *AWSProvider) Init() error {
	return nil
}

func (c *AWSProvider) ListRegions() ([]string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"),
	}))
	svc := ec2.New(sess)
	result, err := svc.DescribeRegions(nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to query aws regions")
		return []string{}, err
	}

	regions := []string{}
	for _, region := range result.Regions {
		regions = append(regions, *region.RegionName)
	}
	return regions, nil
}

type ClusterGetter interface {
	GetClusters(ch chan<- *cluster.CPCluster)
}

func getEKSClients(regions []string) []ClusterGetter {
	clients := make([]ClusterGetter, 0, len(regions))

	for _, region := range regions {
		log.WithFields(log.Fields{
			"region": region,
		}).Debug("Initialize client")
		eks, err := NewEKS(region)
		if err != nil {
			log.WithFields(log.Fields{
				"region": region,
				"error":  err.Error(),
			}).Error("Failed to create AWS SDK session")
			continue
		}

		clients = append(clients, ClusterGetter(eks))
	}
	return clients
}

func (c *AWSProvider) ListClusters(regions []string, setProgress func(int)) ([]*cluster.CPCluster, error) {
	clients := getEKSClients(regions)
	totalRegionClients := len(clients)
	completedRegionClients := 0
	clustersChan := make(chan *cluster.CPCluster)

	var wg sync.WaitGroup
	wg.Add(totalRegionClients)
	for _, c := range clients {
		regionCh := make(chan *cluster.CPCluster)
		go c.GetClusters(regionCh)

		go func(clusterCh chan *cluster.CPCluster) {
			defer wg.Done()
			for cls := range regionCh {
				clusterCh <- cls
			}
		}(clustersChan)
	}

	go func() {
		// close Goroutine, then close channel
		wg.Wait()
		close(clustersChan)
	}()

	var allClusters []*cluster.CPCluster
	for cluster := range clustersChan {
		allClusters = append(allClusters, cluster)
		completedRegionClients++
		progress := (completedRegionClients * 100) / totalRegionClients
		if setProgress != nil {
			setProgress(progress)
		}
	}
	c.GetKubeconfig(allClusters[0])
	return allClusters, nil
}

func (provider *AWSProvider) GetKubeconfig(c *cluster.CPCluster) (*clientcmdapi.Config, error) {
	cfg := c.GenerateClusterConfig(c)

	config := clientcmdapi.NewConfig()
	clusterName := c.Name
	authName := fmt.Sprintf("auth-user-%s", c.Name)
	contextName := c.Name

	config.Clusters[clusterName] = cfg

	authType := getAuthType()
	config.AuthInfos[authName] = getConfigAuthInfo(c, authType)
	config.Contexts[contextName] = &clientcmdapi.Context{
		Cluster:  clusterName,
		AuthInfo: authName,
	}
	config.CurrentContext = contextName

	err := clientcmd.WriteToFile(*config, "/tmp/kubeconfig")
	if err != nil {
		return nil, err
	}
	return config, nil
}
