package cluster

import (
	"fmt"

	"github.com/exoscale/exoscale/csi-driver/internal/integ/flags"
)

func (c *Cluster) tearDownCluster() error {
	id, err := c.getClusterID()
	if err != nil {
		return fmt.Errorf("error getting cluster ID: %w", err)
	}

	cluster, err := c.Ego.GetSKSCluster(c.exoV2Context, *flags.Zone, id)
	if err != nil {
		return err
	}

	if len(cluster.Nodepools) > 0 {
		if err := c.Ego.DeleteSKSNodepool(c.exoV2Context, *flags.Zone, cluster, cluster.Nodepools[0]); err != nil {
			return fmt.Errorf("error deleting nodepool: %w", err)
		}
	}

	return c.Ego.DeleteSKSCluster(c.exoV2Context, *flags.Zone, cluster)
}

func (c *Cluster) TearDown() error {
	if *flags.TearDownCSI {
		if err := c.tearDownCSI(); err != nil {
			return err
		}
	}

	if *flags.TearDownCluster {
		if err := c.tearDownCluster(); err != nil {
			return err
		}
	}

	c.exoV2ContextCancel()

	return nil
}

func (c *Cluster) tearDownCSI() error {
	err := c.K8s.DeleteManifest(c.exoV2Context, mainManifest)
	if err != nil {
		return err
	}

	return err

	// TODO (sauterp) reenable once we have a non-legacy key in GH
	// return c.deleteAPIKeyAndRole()
}