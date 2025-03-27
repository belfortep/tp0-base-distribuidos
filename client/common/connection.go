package common

import (
	"bytes"
	"io"
	"net"
)

// Read up to the delimiter, this function doesn't have short-reads
// Returns a Message object with the logic depending of the message received
func readUpToDelimiter(connection net.Conn, delimiter string) (string, error) {

	buffer := make([]byte, 1024)
	var result bytes.Buffer
	foundDelimiter := false
	delimiterBytes := []byte(delimiter)
	delimiterIndex := 0

	for !foundDelimiter {
		bytesRead, err := connection.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		result.Write(buffer[:bytesRead])

		delimiterIndex = bytes.Index(result.Bytes(), delimiterBytes)
		if delimiterIndex != -1 {
			foundDelimiter = true

		}
	}

	return string(result.Bytes()[:delimiterIndex]), nil
}

// Send a message without short-writes with a \0 at the end
func send(connection net.Conn, messageToSend string) error {
	messageToSend += "\000"
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
