import logging
def send(connection, message):
    """
    Send a message without short-writes with a \0 at the end
    """
    message = message + "\0"
    
    connection.sendall(message.encode("utf-8"))


def read_up_to_delimiter(connection, delimiter):
    """
    Read up to the delimiter, this function doesn't have short-reads
    """
    buffer = bytearray()
    found_delimiter = False
    delimiter = delimiter.encode("utf-8")
    delimiter_index = 0

    while not found_delimiter:
        chunk = connection.recv(2)
        if not chunk:
            raise ConnectionResetError("Socket is closed in reading")
        buffer.extend(chunk)
        logging.info(f"action: SERVER READING | result: success | message {buffer.decode("utf-8")}")

        delimiter_index = buffer.find(delimiter)
        if delimiter_index != -1:
            found_delimiter = True

    return buffer[:delimiter_index].decode("utf-8")