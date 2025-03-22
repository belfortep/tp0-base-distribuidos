package common

import (
	"fmt"
)

type Bet struct {
	agency    int
	name      string
	surname   string
	dni       int
	birthdate string
	number    int
}

func NewBet(agency int, name string, surname string, dni int, birthdate string, number int) *Bet {
	bet := &Bet{
		agency,
		name,
		surname,
		dni,
		birthdate,
		number,
	}

	return bet
}

func (bet *Bet) serialize() string {
	return fmt.Sprintf("%d;%s;%s;%d;%s;%d\000", bet.agency, bet.name, bet.surname, bet.dni, bet.birthdate, bet.number)
}
