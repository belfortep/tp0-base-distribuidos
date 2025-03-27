package common

import (
	"fmt"
	"os"
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

// Get a bet using the env-vars
func GetBet(agency string) Bet {

	name := os.Getenv("CLI_NOMBRE")
	surname := os.Getenv("CLI_APELLIDO")
	dni := os.Getenv("CLI_DOCUMENTO")
	birthdate := os.Getenv("CLI_NACIMIENTO")
	number := os.Getenv("CLI_NUMERO")

	return Bet{
		agency,
		name,
		surname,
		dni,
		birthdate,
		number,
	}
}

// Serialize all the fields of the bet
func (bet *Bet) serialize() string {
	return fmt.Sprintf("%s;%s;%s;%s;%s;%s", bet.agency, bet.name, bet.surname, bet.dni, bet.birthdate, bet.number)
}
