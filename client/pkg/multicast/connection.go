// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package multicast

import (
	"cloud.google.com/go/pubsub"
	"fmt"
	"github.com/GoogleCloudPlatform/bq-datashare-toolkit/client/internal/pubsubutil"
	log "github.com/sirupsen/logrus"
	"net"
)

const (
	maxDatagramSize = 4096
)

// Client for Multicast
type Client struct {
	// UDP multicast address in <HOST:PORT> format
	Address string
	// UDP multicast interface name
	IfName string
	// UDP multicast network group name
	Network string

	// Connection instance
	Conn *net.UDPConn
	// PubSub Topic instance
	Topic *pubsub.Topic
}

// Create Conn binds to the Client UDP network with address:port and returns Connection or error
func (c *Client) CreateConn() error {

	addr, err := net.ResolveUDPAddr(c.Network, c.Address)
	if err != nil {
		return fmt.Errorf("net.ResolveUDPAddr: %s", err)
	}

	ifi, err := net.InterfaceByName(c.IfName)
	if err != nil {
		return fmt.Errorf("net.InterfaceByName: %s", err)
	}

	conn, err := net.ListenMulticastUDP(c.Network, ifi, addr)
	if err != nil {
		return fmt.Errorf("net.ListenMulticastUDP: %s", err)
	}

	log.Debugf("Multicast Listener created: network '%s', address '%s', ifName '%s'",
		c.Network, c.Address, c.IfName)
	c.Conn = conn
	return nil
}

// Create Topic creates a topic client to connect
func (c *Client) CreateTopicClient(projectID, topicID string) error {
	topic, err := pubsubutil.CreateTopicClient(projectID, topicID)
	if err != nil {
		return fmt.Errorf("pubsubutil.CreateTopicClient: %s", err)
	}
	c.Topic = topic
	return nil
}