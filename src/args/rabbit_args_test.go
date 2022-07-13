package args

import (
	"regexp"
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/data/consts"

	"github.com/stretchr/testify/assert"
)

func TestSetGlobalArgs_Default(t *testing.T) {
	argList := ArgumentList{}
	err := SetGlobalArgs(argList)
	assert.NoError(t, err, "err should be nil")
}

func TestSetGlobalArgs_BadArgs(t *testing.T) {
	argList := ArgumentList{
		Exchanges: "invalid",
	}
	err := SetGlobalArgs(argList)
	assert.Error(t, err, "should have error from bad arg when unmarshaling")

	argList.Exchanges = ""
	argList.ExchangesRegexes = `[,]`
	err = SetGlobalArgs(argList)
	assert.Error(t, err)

	argList.ExchangesRegexes = ""
	argList.Queues = "invalid"
	err = SetGlobalArgs(argList)
	assert.Error(t, err)

	argList.Queues = ""
	argList.QueuesRegexes = `[,]`
	err = SetGlobalArgs(argList)
	assert.Error(t, err)

	argList.QueuesRegexes = ""
	argList.Vhosts = "invalid"
	err = SetGlobalArgs(argList)
	assert.Error(t, err)

	argList.Vhosts = ""
	argList.VhostsRegexes = `[,]`
	err = SetGlobalArgs(argList)
	assert.Error(t, err)

	argList.VhostsRegexes = `["(invalid-group"]`
	err = SetGlobalArgs(argList)
	assert.Error(t, err)
}

func TestSetGlobalArgs_ValidJson(t *testing.T) {
	argList := ArgumentList{
		Queues:        `["test-1", "test-2", "test-3"]`,
		QueuesRegexes: `["one-.*", "two-.*"]`,
	}
	err := SetGlobalArgs(argList)
	assert.NoError(t, err)
	assert.NotNil(t, GlobalArgs)
	assert.Equal(t, 3, len(GlobalArgs.Queues))
	assert.Equal(t, 2, len(GlobalArgs.QueuesRegexes))
	assert.True(t, GlobalArgs.QueuesRegexes[0].MatchString("one-queue"))
	assert.False(t, GlobalArgs.QueuesRegexes[0].MatchString("two-queue"))
	assert.False(t, GlobalArgs.QueuesRegexes[1].MatchString("one-queue"))
	assert.True(t, GlobalArgs.QueuesRegexes[1].MatchString("two-queue"))
}

func TestRabbitMQArguments_include(t *testing.T) {
	testRegex, _ := regexp.Compile("four-.*")
	testArgs := RabbitMQArguments{
		Exchanges:        []string{"one"},
		ExchangesRegexes: []*regexp.Regexp{testRegex},
		Queues:           []string{"two"},
		QueuesRegexes:    []*regexp.Regexp{testRegex},
		Vhosts:           []string{"three"},
		VhostsRegexes:    []*regexp.Regexp{testRegex},
	}
	assert.True(t, testArgs.includeExchange("one"))
	assert.False(t, testArgs.includeExchange("false"))
	assert.True(t, testArgs.includeExchange("four-exchange"))

	assert.True(t, testArgs.includeQueue("two"))
	assert.False(t, testArgs.includeQueue("false"))
	assert.True(t, testArgs.includeQueue("four-queue"))

	assert.True(t, testArgs.includeVhost("three"))
	assert.False(t, testArgs.includeVhost("false"))
	assert.True(t, testArgs.includeVhost("four-vhost"))
}

func TestRabbitMQArguments_IncludeEntity(t *testing.T) {
	testArgs := RabbitMQArguments{
		Exchanges: []string{"one"},
		Queues:    []string{"two"},
		Vhosts:    []string{"three"},
	}
	assert.True(t, testArgs.IncludeEntity("one", consts.ExchangeType, "three"))
	assert.True(t, testArgs.IncludeEntity("two", consts.QueueType, "three"))
	assert.True(t, testArgs.IncludeEntity("five", consts.NodeType, "three"))
	assert.False(t, testArgs.IncludeEntity("one", consts.ExchangeType, ""))

	testArgs = RabbitMQArguments{}
	assert.True(t, testArgs.IncludeEntity("any", consts.VhostType, "any"))
}
