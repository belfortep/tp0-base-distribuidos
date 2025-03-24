package common

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	Batchs        int
}

// Client Entity that encapsulates how
type Client struct {
	config        ClientConfig
	conn          net.Conn
	signalChannel chan os.Signal
	lastBatchLine string
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {

	client := &Client{
		config:        config,
		signalChannel: make(chan os.Signal, 1),
		lastBatchLine: "",
	}

	signal.Notify(client.signalChannel, syscall.SIGTERM)

	return client
}

func (client *Client) createBatch(reader *bufio.Reader) (Batch, error) {
	batch := NewBatch(client.config.Batchs)

	for {
		var line string
		var err error

		if client.lastBatchLine != "" {
			line = client.lastBatchLine
			client.lastBatchLine = ""
		} else {
			line, err = reader.ReadString('\n')
		}

		if err == io.EOF {
			log.Infof("action: EOF | result: success | message: returning the rest of the batch")
			return batch, nil
		}

		if err != nil {
			return batch, err
		}

		betValues := strings.Split(strings.TrimSpace(line), ",")

		if len(betValues) != 5 {
			continue
		}
		bet := Bet{
			agency:    client.config.ID,
			name:      betValues[0],
			surname:   betValues[1],
			dni:       betValues[2],
			birthdate: betValues[3],
			number:    betValues[4],
		}

		if batch.CantAppend(bet) {
			client.lastBatchLine = line
			log.Infof("action: cant_append | result: success | message: %v ",
				bet.Serialize(),
			)
			return batch, nil
		}

		batch.Append(bet)
	}

}

func (client *Client) sendBatch(batch Batch) error {
	err := send(client.conn, batch.Serialize())

	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
			client.config.ID,
			err,
		)
		return err
	}

	message, err := readUpToDelimiter(client.conn, "\000")

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			client.config.ID,
			err,
		)
		return err
	}
	message.ActionForClient()

	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (client *Client) StartClientLoop() {

	go client.shutdownClientHandler()

	filepath := fmt.Sprintf("/.data/agency-%v.csv", client.config.ID)
	file, err := os.Open(filepath)
	if err != nil {
		log.Errorf("action: open_file | result: fail | client_id: %v | error: %v",
			client.config.ID,
			err,
		)
		return
	}

	reader := bufio.NewReader(file)
	defer file.Close()

	for {
		client.createClientSocket()
		batch, err := client.createBatch(reader)

		if err != nil {
			log.Errorf("action: create_batch | result: fail | client_id: %v | error: %v",
				client.config.ID,
				err,
			)
			return
		}

		if batch.isEmpty() {
			client.conn.Close()
			break
		}

		err = client.sendBatch(batch)

		if err != nil {
			log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v",
				client.config.ID,
				err,
			)
			return
		}

		client.conn.Close()
		// Wait a time between sending one message and the next one
		time.Sleep(client.config.LoopPeriod)
	}

	for {
		client.createClientSocket()
		message, err := client.getWinners()

		if err != nil {
			log.Errorf("action: get_winners | result: fail | client_id: %v | error: %v",
				client.config.ID,
				err,
			)
			return
		}

		message.ActionForClient()
		time.Sleep(client.config.LoopPeriod)

		if message.MessageType() == WINNERS_MESSAGE {
			return
		}

		client.conn.Close()
	}

}

func (client *Client) getWinners() (Message, error) {
	err := send(client.conn, fmt.Sprintf("GETWINNERS;%v", client.config.ID))

	if err != nil {
		log.Errorf("action: get_winners | result: fail | client_id: %v",
			client.config.ID,
		)
		return nil, err
	}

	message, err := readUpToDelimiter(client.conn, "\000")

	if err != nil {
		log.Errorf("action: get_winners | result: fail | client_id: %v",
			client.config.ID,
		)
		return nil, err
	}
	return message, nil
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
