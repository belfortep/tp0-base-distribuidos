def send(connection, message):
    byte_to_write = len(message)
    bytes_already_written = 0
    while bytes_already_written < byte_to_write:
        bytes_written = connection.send(message[bytes_already_written:].encode('utf-8'))   
        if bytes_written == 0:
            raise Exception("Socket is closed in sending")
        bytes_already_written = bytes_already_written + bytes_written


def read_up_to_delimiter(connection, delimiter):
    buffer = bytearray()
    fund_delimiter = False

    while not fund_delimiter:
        chunk = connection.recv(1024)
        if not chunk:
            raise Exception("Socket is closed in reading")
        buffer.extend(chunk)

        if delimiter in chunk:
            fund_delimiter = True
        
    return buffer.decode("utf-8").rstrip()