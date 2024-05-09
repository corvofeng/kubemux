// Package internal provides wrapper for creating aws sessions
package aws

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"testing"

	cluster "kubemux/lib/kubernetes"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

//nolint:unused, varcheck, deadcode
var update = flag.Bool("update", false, "update .golden files")

type mockEKSClient struct {
	eksiface.EKSAPI
	// The clusters we need to return
	Clusters []*cluster.Cluster

	// What is the size of the pages
	PageSize int

	// DescribeCluster will fail in the calls specified here (0, 3) means
	// that the first call and the fourth one will fail
	ErrorOnDescribe map[int]error

	// ListClustersPages will fail in the calls specified here (0, 3) means
	// that the first call and the fourth one will fail
	ErrorOnList map[int]error

	ListCallCount     int
	DescribeCallCount int
}

// TODO(mmicu): implement DescribeCluster and ListClustersPages

func (c *mockEKSClient) ListClustersPages(_ *eks.ListClustersInput, fn func(*eks.ListClustersOutput, bool) bool) error {
	if err, ok := c.ErrorOnList[c.ListCallCount]; ok {
		c.ListCallCount++
		return err
	}
	start := 0
	end := min(c.PageSize, len(c.Clusters))
	for end <= len(c.Clusters) {
		o := eks.ListClustersOutput{}
		clusters := []*string{}

		// prepare clusters
		for _, cls := range c.Clusters[start:end] {
			localCluster := *cls
			clusters = append(clusters, &localCluster.Name)
		}
		o.Clusters = clusters
		lastPage := end == len(c.Clusters)
		fn(&o, lastPage)

		if lastPage {
			break
		}

		step := min(c.PageSize, len(c.Clusters)-end)
		start += c.PageSize
		end += step
	}

	return nil
}

func (c *mockEKSClient) DescribeCluster(input *eks.DescribeClusterInput) (*eks.DescribeClusterOutput, error) {
	defer func() {
		c.DescribeCallCount++
	}()
	if err, ok := c.ErrorOnDescribe[c.DescribeCallCount]; ok {
		return nil, err
	}

	for _, cls := range c.Clusters {
		localCluster := *cls
		if *input.Name == cls.Name {
			cluster := eks.Cluster{}
			cluster.Arn = &localCluster.ID
			cluster.Endpoint = &localCluster.Endpoint
			cluster.Name = &localCluster.Name
			cluster.Status = &localCluster.Status

			cert := eks.Certificate{}
			data := base64.StdEncoding.EncodeToString([]byte(cls.CertificateAuthorityData))
			cert.Data = &data
			cluster.CertificateAuthority = &cert

			// TODO(mmicu): populate cluster data here
			out := eks.DescribeClusterOutput{}
			out.Cluster = &cluster
			return &out, nil
		}
	}
	return nil, fmt.Errorf("can't find cluster %v", input.Name)
}

type testCase struct {
	Client mockEKSClient
	Region string
}

var cases = []testCase{
	// Happy flows
	{
		Client: mockEKSClient{
			Clusters:          cluster.GetMockClusters(1),
			PageSize:          1,
			ErrorOnDescribe:   map[int]error{},
			ErrorOnList:       map[int]error{},
			ListCallCount:     0,
			DescribeCallCount: 0,
		},
		Region: "fakeRegion",
	},
	{
		Client: mockEKSClient{
			Clusters:          cluster.GetMockClusters(3),
			PageSize:          1,
			ErrorOnDescribe:   map[int]error{},
			ErrorOnList:       map[int]error{},
			ListCallCount:     0,
			DescribeCallCount: 0,
		},
		Region: "fakeRegion",
	},
	{
		Client: mockEKSClient{
			Clusters:          cluster.GetMockClusters(3),
			PageSize:          3,
			ErrorOnDescribe:   map[int]error{},
			ErrorOnList:       map[int]error{},
			ListCallCount:     0,
			DescribeCallCount: 0,
		},
		Region: "fakeRegion",
	},
	{
		Client: mockEKSClient{
			Clusters:          cluster.GetMockClusters(3),
			PageSize:          4,
			ErrorOnDescribe:   map[int]error{},
			ErrorOnList:       map[int]error{},
			ListCallCount:     0,
			DescribeCallCount: 0,
		},
		Region: "fakeRegion",
	},
	{
		Client: mockEKSClient{
			Clusters:          cluster.GetMockClusters(100),
			PageSize:          7,
			ErrorOnDescribe:   map[int]error{},
			ErrorOnList:       map[int]error{},
			ListCallCount:     0,
			DescribeCallCount: 0,
		},
		Region: "fakeRegion",
	},
	{
		Client: mockEKSClient{
			Clusters:          cluster.GetMockClusters(1),
			PageSize:          1,
			ErrorOnDescribe:   map[int]error{0: fmt.Errorf("can't Describe Cluster 0")},
			ErrorOnList:       map[int]error{},
			ListCallCount:     0,
			DescribeCallCount: 0,
		},
		Region: "fakeRegion",
	},
	{
		Client: mockEKSClient{
			Clusters: cluster.GetMockClusters(3),
			PageSize: 1,
			ErrorOnDescribe: map[int]error{
				0: fmt.Errorf("can't Describe Cluster 0"),
				2: fmt.Errorf("can't Describe Cluster 2"),
			},
			ErrorOnList:       map[int]error{},
			ListCallCount:     0,
			DescribeCallCount: 0,
		},
		Region: "fakeRegion",
	},
	{
		Client: mockEKSClient{
			Clusters: cluster.GetMockClusters(3),
			PageSize: 3,
			ErrorOnDescribe: map[int]error{
				0: fmt.Errorf("can't Describe Cluster 0"),
				2: fmt.Errorf("can't Describe Cluster 2"),
			},
			ErrorOnList:       map[int]error{},
			ListCallCount:     0,
			DescribeCallCount: 0,
		},
		Region: "fakeRegion",
	},
	{
		Client: mockEKSClient{
			Clusters: cluster.GetMockClusters(3),
			PageSize: 4,
			ErrorOnDescribe: map[int]error{
				0: fmt.Errorf("can't Describe Cluster 0"),
				2: fmt.Errorf("can't Describe Cluster 2"),
			},
			ErrorOnList:       map[int]error{},
			ListCallCount:     0,
			DescribeCallCount: 0,
		},
		Region: "fakeRegion",
	},
	{
		Client: mockEKSClient{
			Clusters: cluster.GetMockClusters(100),
			PageSize: 7,
			ErrorOnDescribe: map[int]error{
				0:  fmt.Errorf("can't Describe Cluster 0"),
				2:  fmt.Errorf("can't Describe Cluster 2"),
				66: fmt.Errorf("can't Describe Cluster 66"),
				97: fmt.Errorf("can't Describe Cluster 87"),
			},
			ErrorOnList:       map[int]error{},
			ListCallCount:     0,
			DescribeCallCount: 0,
		},
		Region: "fakeRegion",
	},
}

func TestGetClustersNoFailure(t *testing.T) {
	t.Parallel()
	log.SetOutput(ioutil.Discard)
	for _, tt := range cases {
		client := tt.Client
		describeErrorCount := 0
		for k := range tt.Client.ErrorOnDescribe {
			if k > len(tt.Client.Clusters) {
				continue
			}
			describeErrorCount++
		}

		listErrorCount := 0
		for range tt.Client.ErrorOnList {
			listErrorCount++
		}
		testname := fmt.Sprintf(
			"Get %v clusters in batches of %v [descErr: %v, listErr: %v]",
			len(tt.Client.Clusters), tt.Client.PageSize,
			describeErrorCount, listErrorCount,
		)
		t.Run(testname, func(t *testing.T) {
			ch := make(chan *cluster.Cluster)
			c := EKSClient{
				EKS:    &client,
				Region: tt.Region,
			}
			go c.GetClusters(ch)
			clusters := []*cluster.Cluster{}
			for c := range ch {
				clusters = append(clusters, c)
			}

			assert.Equal(t, len(tt.Client.Clusters)-describeErrorCount, len(clusters))

			// fix Regions
			for _, c := range tt.Client.Clusters {
				c.Region = tt.Region
			}

			// nillify function fields
			for _, c := range tt.Client.Clusters {
				c.GenerateClusterConfig = nil
				c.GenerateAuthInfo = nil
			}

			for _, c := range clusters {
				c.GenerateClusterConfig = nil
				c.GenerateAuthInfo = nil
			}

			if describeErrorCount == 0 {
				assert.ElementsMatch(t, tt.Client.Clusters, clusters)
			}
			assert.Subset(t, tt.Client.Clusters, clusters)
		})
	}
}

func TestGetClustersListFailure(t *testing.T) {
	t.Parallel()
	log.SetOutput(ioutil.Discard)

	tts := []testCase{
		{
			Client: mockEKSClient{
				Clusters:          cluster.GetMockClusters(3),
				PageSize:          1,
				ErrorOnDescribe:   map[int]error{0: fmt.Errorf("can't Describe Cluster 0")},
				ErrorOnList:       map[int]error{0: fmt.Errorf("can't List Clusters 0")},
				ListCallCount:     0,
				DescribeCallCount: 0,
			},
			Region: "fakeRegion",
		},
	}
	for _, tt := range tts {
		client := tt.Client
		describeErrorCount := 0
		for k := range tt.Client.ErrorOnDescribe {
			if k > len(tt.Client.Clusters) {
				continue
			}
			describeErrorCount++
		}

		listErrorCount := 0
		for range tt.Client.ErrorOnList {
			listErrorCount++
		}
		testname := fmt.Sprintf(
			"Get %v clusters in batches of %v [descErr: %v, listErr: %v]",
			len(tt.Client.Clusters), tt.Client.PageSize,
			describeErrorCount, listErrorCount,
		)
		t.Run(testname, func(t *testing.T) {
			ch := make(chan *cluster.Cluster)
			c := EKSClient{
				EKS:    &client,
				Region: tt.Region,
			}
			go c.GetClusters(ch)
			clusters := []*cluster.Cluster{}
			for c := range ch {
				clusters = append(clusters, c)
			}

			// this test assumes that there is at least one ErrorOnList
			assert.True(t, listErrorCount > 0)

			// Error on list
			// for now this are fatals but in the future we may hande diferite errors
			assert.Equal(t, 0, len(clusters))
		})
	}
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
