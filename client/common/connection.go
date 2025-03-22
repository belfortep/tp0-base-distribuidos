package common

import (
	"io"
	"net"
)

func read(connection net.Conn, bufferSize int, messageBuffer []byte) (string, error) {
	bytesAlreadyRead := 0
	for bytesAlreadyRead < bufferSize {
		bytesRead, err := connection.Read(messageBuffer[bytesAlreadyRead:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		bytesAlreadyRead += bytesRead
	}

	log.Infof("action: reading | result: success | message: %v ",
		string(messageBuffer),
	)

	return string(messageBuffer[:bytesAlreadyRead]), nil
}

func send(connection net.Conn, messageToSend string) error {
	bytesToWrite := len(messageToSend)
	bytesAlreadyWritten := 0

	for bytesAlreadyWritten < bytesToWrite {
		bytesWritten, err := connection.Write([]byte(messageToSend[bytesAlreadyWritten:]))
		if err != nil {
			return err
		}

		bytesAlreadyWritten += bytesWritten
	}
	return nil
}
