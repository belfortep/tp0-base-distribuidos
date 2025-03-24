package common

import "strings"

const ACK_MESSAGE = "ACK"
const NOTYET_MESSAGE = "NOTYET"
const ERR_MESSAGE = "ERR"
const WINNERS_MESSAGE = "WINNERS"
const WRONG_MESSAGE = "WRONG"

type Message interface {
	ActionForClient()
	MessageType() string
}

func CreateMessage(message string) Message {

	if strings.HasPrefix(message, ACK_MESSAGE) {
		return &MessageACK{}
	} else if strings.HasPrefix(message, NOTYET_MESSAGE) {
		return &MessageNotYet{}
	} else if strings.HasPrefix(message, ERR_MESSAGE) {
		errors := strings.Split(message, ";")
		if len(errors) != 2 {
			return nil
		}
		number_of_errors := errors[1]
		return &MessageError{

			number_of_errors: number_of_errors,
		}
	} else if strings.HasPrefix(message, WINNERS_MESSAGE) {
		winners := strings.Split(message, ";")

		return &MessageWinners{
			number_of_winners: len(winners) - 1,
		}
	} else {
		return &MessageWrong{}
	}
}

type MessageWrong struct {
}

func (message *MessageWrong) ActionForClient() {
	log.Errorf("action: receive_message | result: fail | message: wrong message received")
}

func (message *MessageWrong) MessageType() string {
	return WRONG_MESSAGE
}

type MessageWinners struct {
	number_of_winners int
}

func (message *MessageWinners) ActionForClient() {
	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", message.number_of_winners)
}

func (message *MessageWinners) MessageType() string {
	return WINNERS_MESSAGE
}

type MessageNotYet struct {
}

func (message *MessageNotYet) ActionForClient() {
	log.Infof("action: sleeping | result: success")
}

func (message *MessageNotYet) MessageType() string {
	return NOTYET_MESSAGE
}

type MessageACK struct {
}

func (message *MessageACK) ActionForClient() {
	log.Infof("action: batch_send | result: success")
}

func (message *MessageACK) MessageType() string {
	return ACK_MESSAGE
}

type MessageError struct {
	number_of_errors string
}

func (message *MessageError) ActionForClient() {
	log.Errorf("action: receive_message | result: fail | number_of_errors: %d",
		message.number_of_errors,
	)
}

func (message *MessageError) MessageType() string {
	return ERR_MESSAGE
}
