package strategy

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestPriceRule(t *testing.T) {
	pr := priceRule{
		lower:    1.0,
		upper:    2.0,
		provider: &UniswapProvider{price: 1.5},
	}
	if !pr.check() {
		t.Error("Price rule failed")
	}
}

type BotMock struct {
	mock.Mock
}

func (b *BotMock) sell(price float64) bool {
	args := b.Called(price)
	return args.Bool(0)
}

func TestSellCommand(t *testing.T) {
	bot := &BotMock{}
	bot.On("sell", 2.0).Return(true)

	sc := SellCommand{
		amount: 2.0,
		bot:    bot,
	}

	if !sc.execute() {
		t.Error("Sell command failed")
	}
}

func TestStrategy(t *testing.T) {
	bot := &BotMock{}
	bot.On("sell", 7.3).Return(true)

	sc := SellCommand{
		amount: 7.3,
		bot:    bot,
	}

	pr := priceRule{
		lower:    1.0,
		upper:    200.0,
		provider: &UniswapProvider{price: 1.2},
	}
	s := Strategy{
		rules:   []rule{&pr},
		command: &sc,
	}

	if !s.run() {
		t.Error("Strategy failed")
	}
}

type RuleMockOk struct{}

func (r *RuleMockOk) check() bool {
	return true
}

type RuleMockFail struct{}

func (r *RuleMockFail) check() bool {
	return false
}

type CommandMock struct {
	mock.Mock
}

func (c *CommandMock) execute() bool {
	args := c.Called()
	return args.Bool(0)
}

func TestStrategyOk(t *testing.T) {
	r := &RuleMockOk{}
	command := &CommandMock{}
	command.On("execute").Return(true)

	s := Strategy{
		rules:   []rule{r},
		command: command,
	}

	if !s.run() {
		t.Error("Strategy failed")
	}
}

func TestStrategyOkCommandFail(t *testing.T) {
	r := &RuleMockOk{}
	command := &CommandMock{}
	command.On("execute").Return(false)

	s := Strategy{
		rules:   []rule{r},
		command: command,
	}

	if s.run() {
		t.Error("Strategy succeeded when command failed")
	}
}

func TestStrategyMultipleRulesFail(t *testing.T) {
	r := &RuleMockOk{}
	r2 := &RuleMockFail{}
	command := &CommandMock{}
	command.On("execute").Return(true)

	s := Strategy{
		rules:   []rule{r, r2},
		command: command,
	}

	if s.run() {
		t.Error("Strategy succeeded when rules failed")
	}
}
