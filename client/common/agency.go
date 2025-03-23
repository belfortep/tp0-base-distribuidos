package common

import (
	"fmt"
)

type Bet struct {
	agency    string
	name      string
	surname   string
	dni       string
	birthdate string
	number    string
}

func (bet *Bet) Serialize() string {
	return fmt.Sprintf("%s;%s;%s;%s;%s;%s\n", bet.agency, bet.name, bet.surname, bet.dni, bet.birthdate, bet.number)
}

func (bet *Bet) MemorySize() int {
	return len(bet.Serialize())
}
