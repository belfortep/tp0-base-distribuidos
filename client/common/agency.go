package common

import (
	"fmt"
)

// Bet entity, encapsulates the serialization of bets logic
type Bet struct {
	agency    string
	name      string
	surname   string
	dni       string
	birthdate string
	number    string
}

// Serialize all the fields of the bet
func (bet *Bet) Serialize() string {
	return fmt.Sprintf("%s;%s;%s;%s;%s;%s\n", bet.agency, bet.name, bet.surname, bet.dni, bet.birthdate, bet.number)
}

// Returns the memory needed for the serialization of this bet
func (bet *Bet) MemorySize() int {
	return len(bet.Serialize())
}
