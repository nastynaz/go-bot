package strategy

type rule interface {
	check() bool
}

type ExchangeProvider interface {
	getPrice() float64
}

type UniswapProvider struct {
	price float64
}

func (u *UniswapProvider) getPrice() float64 {
	return u.price
}

type priceRule struct {
	lower    float64
	upper    float64
	provider ExchangeProvider
}

func (p *priceRule) check() bool {
	price := p.provider.getPrice()
	return price >= p.lower && price <= p.upper
}

type command interface {
	execute() bool
}

type Bot interface {
	sell(float64) bool
}

type SellCommand struct {
	amount float64
	bot    Bot
}

func (c *SellCommand) execute() bool {
	return c.bot.sell(c.amount)
}

type Strategy struct {
	rules   []rule
	command command
}

func (s *Strategy) run() bool {
	for _, r := range s.rules {
		if !r.check() {
			return false
		}
	}
	return s.command.execute()
}
