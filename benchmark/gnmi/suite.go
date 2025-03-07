// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package gnmi

import (
	"context"
	"github.com/onosproject/helmit/pkg/benchmark"
	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/input"
	"github.com/onosproject/helmit/pkg/util/random"
	"github.com/onosproject/onos-config/test/utils/charts"
	"github.com/onosproject/onos-config/test/utils/gnmi"
	"github.com/onosproject/onos-test/pkg/onostest"
	"github.com/openconfig/gnmi/client/gnmi"
	"time"
)

// BenchmarkSuite is an onos-config gNMI benchmark suite
type BenchmarkSuite struct {
	benchmark.Suite
	simulator *helm.HelmRelease
	client    *client.Client
	value     input.Source
}

// SetupSuite :: benchmark
func (s *BenchmarkSuite) SetupSuite(c *input.Context) error {
	umbrella := charts.CreateUmbrellaRelease().
		Set("onos-topo.replicaCount", 2).
		Set("onos-config.replicaCount", 2)
	return umbrella.Install(true)
}

// SetupWorker :: benchmark
func (s *BenchmarkSuite) SetupWorker(c *input.Context) error {
	s.value = input.RandomString(8)
	s.simulator = helm.
		Chart("device-simulator", onostest.OnosChartRepo).
		Release(random.NewPetName(2))
	if err := s.simulator.Install(true); err != nil {
		return err
	}
	gnmiClient, err := getGNMIClient()
	if err != nil {
		return err
	}
	s.client = gnmiClient
	return nil
}

// TearDownWorker :: benchmark
func (s *BenchmarkSuite) TearDownWorker(c *input.Context) error {
	s.client.Close()
	return s.simulator.Uninstall()
}

var _ benchmark.SetupWorker = &BenchmarkSuite{}

// getGNMIClient makes a GNMI client to use for requests
func getGNMIClient() (*client.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dest, err := gnmi.GetOnosConfigDestination()
	if err != nil {
		return nil, err
	}
	gnmiClient, err := client.New(ctx, dest)
	if err != nil {
		return nil, err
	}
	return gnmiClient.(*client.Client), nil
}
