package common

import (
	"fmt"
	"net"
)

func read(connection net.Conn, bufferSize int, messageBuffer []byte) (string, error) {
	bytesAlreadyRead := 0
	for bytesAlreadyRead < bufferSize {
		bytesRead, err := connection.Read(messageBuffer[bytesAlreadyRead:])
		if err != nil {
			return "", err
		}

		bytesAlreadyRead += bytesRead
	}

	return string(messageBuffer[:bytesAlreadyRead]), nil
}

func send(connection net.Conn, messageToSend string) error {
	bytesToWrite := len(messageToSend)
	bytesAlreadyWritten := 0

	for bytesAlreadyWritten < bytesToWrite {
		bytesWritten, err := fmt.Fprint(
			connection,
			messageToSend,
		)
		if err != nil {
			return err
		}

		bytesAlreadyWritten += bytesWritten
	}
	return nil
}
