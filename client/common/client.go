package common

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

const ACK_MESSAGE = "ACK"

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config        ClientConfig
	conn          net.Conn
	signalChannel chan os.Signal
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {

	client := &Client{
		config:        config,
		signalChannel: make(chan os.Signal, 1),
	}

	signal.Notify(client.signalChannel, syscall.SIGTERM)

	return client
}

// StartClientLoop Send messages to the client until some time threshold is met
func (client *Client) StartClientLoop() {

	go client.shutdownClientHandler()
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= client.config.LoopAmount; msgID++ {
		// Create the connection the server in every loop iteration. Send an
		client.createClientSocket()
		bet := GetBet(client.config.ID)
		err := send(client.conn, bet.serialize())

		if err != nil {
			log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
				client.config.ID,
				err,
			)
			return
		}

		message, err := client.readACK()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				client.config.ID,
				err,
			)
			return
		}

		if message == ACK_MESSAGE {
			log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
				bet.dni,
				bet.number,
			)
		} /* else {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				client.config.ID,
				message,
			)
			return
		}*/

		client.conn.Close()
		// Wait a time between sending one message and the next one
		time.Sleep(client.config.LoopPeriod)

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", client.config.ID)
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (client *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", client.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			client.config.ID,
			err,
		)
	}
	client.conn = conn
	return nil
}

func (client *Client) readACK() (string, error) {
	buffer := make([]byte, len(ACK_MESSAGE))
	message, err := read(client.conn, len(ACK_MESSAGE), buffer)

	if err != nil {
		return "", err
	}

	return message, nil
}

func (client *Client) shutdownClientHandler() {
	<-client.signalChannel
	close(client.signalChannel)
	if client.conn != nil {
		client.conn.Close()
	}
	log.Infof("action: shutdown_client | result: success | client_id: %v ",
		client.config.ID,
	)

	os.Exit(0)
}
