package common

import (
	"fmt"
	"os"
)

type Bet struct {
	agency    string
	name      string
	surname   string
	dni       string
	birthdate string
	number    string
}

func GetBet(agency string) Bet {

	name := os.Getenv("nombre")
	surname := os.Getenv("apellido")
	dni := os.Getenv("documento")
	birthdate := os.Getenv("nacimiento")
	number := os.Getenv("numero")

	return Bet{
		agency,
		name,
		surname,
		dni,
		birthdate,
		number,
	}
}

func (bet *Bet) serialize() string {
	return fmt.Sprintf("%s;%s;%s;%s;%s;%s\000", bet.agency, bet.name, bet.surname, bet.dni, bet.birthdate, bet.number)
}
