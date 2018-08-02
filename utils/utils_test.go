package utils

import (
	"errors"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-rabbitmq/args"
	"github.com/newrelic/nri-rabbitmq/logger"
	"github.com/newrelic/nri-rabbitmq/testutils"
	"github.com/newrelic/nri-rabbitmq/utils/consts"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateEntity(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{}
	i := testutils.GetTestingIntegration(t)

	expectedEntityName := "/firstEntity"
	expectedMetricEntityName := "queue:" + expectedEntityName
	e1, metricNS, err := CreateEntity(i, "firstEntity", consts.QueueType, "/")
	assert.NoError(t, err)
	assert.NotNil(t, e1)
	assert.Equal(t, expectedEntityName, e1.Metadata.Name)
	assert.Equal(t, consts.QueueType, e1.Metadata.Namespace)
	assert.Equal(t, 2, len(metricNS))
	assert.Contains(t, metricNS, metric.Attribute{
		Key: "displayName", Value: expectedEntityName,
	})
	assert.Contains(t, metricNS, metric.Attribute{
		Key: "entityName", Value: expectedMetricEntityName,
	})

	expectedEntityName = "/vhost2/" + consts.DefaultExchangeName
	expectedMetricEntityName = "exchange:" + expectedEntityName
	e2, metricNS, err := CreateEntity(i, "", consts.ExchangeType, "/vhost2")
	assert.NoError(t, err)
	assert.NotNil(t, e2)
	assert.NotNil(t, metricNS)
	assert.Equal(t, expectedEntityName, e2.Metadata.Name)
	assert.Equal(t, consts.ExchangeType, e2.Metadata.Namespace)
	assert.Equal(t, 2, len(metricNS))
	assert.Contains(t, metricNS, metric.Attribute{
		Key: "displayName", Value: expectedEntityName,
	})
	assert.Contains(t, metricNS, metric.Attribute{
		Key: "entityName", Value: expectedMetricEntityName,
	})

	existingArgs := args.GlobalArgs
	defer func() {
		args.GlobalArgs = existingArgs
	}()
	args.GlobalArgs = args.RabbitMQArguments{
		Queues: []string{"missing-queue"},
	}

	e3, metricNS, err := CreateEntity(i, "actual-queue", consts.QueueType, "/")
	assert.Nil(t, e3)
	assert.Nil(t, metricNS)
	assert.Nil(t, err)
}

func TestConvertBoolToIntTrue(t *testing.T) {
	actual := ConvertBoolToInt(true)
	expected := 1
	assert.Equal(t, expected, actual)
}

func TestConvertBoolToIntFalse(t *testing.T) {
	actual := ConvertBoolToInt(false)
	expected := 0
	assert.Equal(t, expected, actual)
}

func TestCheckErr(t *testing.T) {
	l := new(testutils.MockedLogger)
	logger.SetLogger(l)
	l.On("Errorf", mock.Anything, mock.Anything).Once()
	CheckErr(func() error {
		return nil
	})
	CheckErr(func() error {
		return errors.New("Test Error")
	})
	l.AssertExpectations(t)
}

func TestPanicErr(t *testing.T) {
	assert.NotPanics(t, func() {
		PanicOnErr(nil)
	}, "panicOnEror should not panic when not given an error")

	assert.Panics(t, func() {
		PanicOnErr(errors.New("Panic Error"))
	}, "panicOnEror should panic when given an error")
}
