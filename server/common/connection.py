def send(connection, message):

    connection.sendall(message.encode("utf-8"))


def read_up_to_delimiter(connection, delimiter):
    buffer = bytearray()
    fund_delimiter = False
    delimiter = delimiter.encode("utf-8")

    while not fund_delimiter:
        chunk = connection.recv(1024)
        if not chunk:
            raise Exception("Socket is closed in reading")
        buffer.extend(chunk)

        if delimiter in chunk:
            fund_delimiter = True

    return buffer[:-len(delimiter)].decode("utf-8")