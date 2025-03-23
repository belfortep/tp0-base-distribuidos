def send(connection, message):

    connection.sendall(message.encode("utf-8"))


def read_up_to_delimiter(connection, delimiter):
    buffer = bytearray()
    fund_delimiter = False
    delimiter = delimiter.encode("utf-8")
    delimiter_index = 0

    while not fund_delimiter:
        chunk = connection.recv(3)
        if not chunk:
            raise Exception("Socket is closed in reading")
        buffer.extend(chunk)

        delimiter_index = buffer.find(delimiter)
        if delimiter_index != -1:
            fund_delimiter = True

    return buffer[:delimiter_index].decode("utf-8")