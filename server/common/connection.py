def send(connection, message):
    message = message + "\0"
    connection.sendall(message.encode("utf-8"))


def read_up_to_delimiter(connection, delimiter):
    buffer = bytearray()
    found_delimiter = False
    delimiter = delimiter.encode("utf-8")
    delimiter_index = 0

    while not found_delimiter:
        chunk = connection.recv(1024)
        if not chunk:
            raise ConnectionResetError("Socket is closed in reading")
        buffer.extend(chunk)

        delimiter_index = buffer.find(delimiter)
        if delimiter_index != -1:
            found_delimiter = True

    return buffer[:delimiter_index].decode("utf-8")