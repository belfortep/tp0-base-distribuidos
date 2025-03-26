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
	is_running    bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {

	client := &Client{
		config:        config,
		signalChannel: make(chan os.Signal, 1),
		is_running:    true,
	}

	signal.Notify(client.signalChannel, syscall.SIGTERM)

	return client
}

// StartClientLoop Send messages to the client until some time threshold is met
func (client *Client) StartClientLoop() {

	go client.shutdownClientHandler()

	err := client.createClientSocket()

	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			client.config.ID,
			err,
		)
		return
	}
	defer client.conn.Close()
	bet := GetBet(client.config.ID)
	err = send(client.conn, bet.serialize())

	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
			client.config.ID,
			err,
		)
		return
	}

	message, err := readUpToDelimiter(client.conn, "\000")

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
		client.is_running = false
	} else {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			client.config.ID,
			message,
		)
		return
	}
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (client *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", client.config.ServerAddress)
	if err != nil {
		return err
	}
	client.conn = conn
	return nil
}

func (client *Client) shutdownClientHandler() {
	<-client.signalChannel
	close(client.signalChannel)
	if client.conn != nil {
		client.conn.Close()
	}
	client.is_running = false
	log.Infof("action: shutdown_client | result: success | client_id: %v ",
		client.config.ID,
	)

}
