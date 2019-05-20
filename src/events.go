package main

import (
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/event"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/data/consts"
)

const (
	success = "ok"
)

func alivenessTest(rabbitmqIntegration *integration.Integration, vhostTests []*data.VhostTest, clusterName string) {
	if rabbitmqIntegration != nil {
		for _, vhostTest := range vhostTests {
			if vhostTest.Test.Status != success {
				e, _, err := data.CreateEntity(rabbitmqIntegration, vhostTest.Vhost.Name, consts.VhostType, vhostTest.Vhost.Name, clusterName)
				if err != nil {
					log.Error("Error creating vhost entity [%s]: %v", vhostTest.Vhost.Name, err)
					continue
				}

        // Don't add events for the entity if we are skipping its collection
        if e != nil {
          description := fmt.Sprintf("Response [%s] for vhost [%s]: %s", vhostTest.Test.Status, vhostTest.Vhost.Name, vhostTest.Test.Reason)
          warnIfError(e.AddEvent(event.New(description, "integration")), "Error adding event: %v")
        }
			}
		}
	}
}

func healthcheckTest(rabbitmqIntegration *integration.Integration, nodeTests []*data.NodeTest, clusterName string) {
	if rabbitmqIntegration != nil {
		for _, nodeTest := range nodeTests {
			if nodeTest.Test.Status != success {
				e, _, err := nodeTest.Node.GetEntity(rabbitmqIntegration, clusterName)
				if err != nil {
					log.Error("Error creating node entity [%s]: %v", nodeTest.Node.Name, err)
					return
				}

        // Don't add events for the entity if we are skipping its collection
        if e != nil {
          description := fmt.Sprintf("Response [%s] for node [%s]: %s", nodeTest.Test.Status, nodeTest.Node.Name, nodeTest.Test.Reason)
          warnIfError(e.AddEvent(event.New(description, "integration")), "Error adding event: %v")
        }
			}
		}
	}
}
