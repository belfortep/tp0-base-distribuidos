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

const ACK_MESSAGE = "ACK"

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
	is_running    bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {

	client := &Client{
		config:        config,
		signalChannel: make(chan os.Signal, 1),
		lastBatchLine: "",
		is_running:    true,
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

		err := client.createClientSocket()

		if err != nil {
			log.Criticalf(
				"action: connect | result: fail | client_id: %v | error: %v",
				client.config.ID,
				err,
			)
			return
		}

		batch, err := client.createBatch(reader)

		if err != nil {
			log.Errorf("action: create_batch | result: fail | client_id: %v | error: %v",
				client.config.ID,
				err,
			)
			client.conn.Close()
			return
		}

		if batch.isEmpty() {
			client.conn.Close()
			break
		}

		err = send(client.conn, batch.Serialize())

		if err != nil {
			log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
				client.config.ID,
				err,
			)
			client.conn.Close()
			return
		}

		message, err := readUpToDelimiter(client.conn, "\000")

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				client.config.ID,
				err,
			)
			client.conn.Close()
			return
		}
		client.conn.Close()
		if message == ACK_MESSAGE {
			log.Infof("action: batch_send | result: success")
		} else {
			log.Errorf("action: receive_message | result: fail | client_id: %v | message: %v",
				client.config.ID,
				message,
			)
			return
		}

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
