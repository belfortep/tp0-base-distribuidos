import socket
import logging
import signal

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._is_running = True
        self._last_client_socket = None

        signal.signal(signal.SIGTERM, self.__shutdown_server)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        try:
            while self._is_running:
                self._last_client_socket = self.__accept_new_connection()
                if self._last_client_socket:
                    self.__handle_client_connection()
        except Exception as e:
            logging.error("action: server_run | result: fail | error: {e}")
        finally:
            self.__shutdown_server(None, None)
            

    def __handle_client_connection(self):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            message = self.__read_up_to_delimiter(b'\0')
            addr = self._last_client_socket.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {message}')
            self.__send("ACK\0")
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            self._last_client_socket.close()
            self._last_client_socket = None
    
    def __send(self, message):
        byte_to_write = len(message)
        bytes_already_written = 0
        while bytes_already_written < byte_to_write:
            bytes_written = self._last_client_socket.send(message[bytes_already_written:].encode('utf-8'))   
            if bytes_written == 0:
                raise Exception("Socket is closed")
            bytes_already_written = bytes_already_written + bytes_written
    
    def __read_up_to_delimiter(self, delimiter):
        buffer = bytearray()
        fund_delimiter = False

        while not fund_delimiter:
            chunk = self._last_client_socket.recv(1024)
            if not chunk:
                raise Exception("Socket is closed in reading")
            buffer.extend(chunk)

            if delimiter in chunk:
                fund_delimiter = True
        
        return buffer.decode("utf-8").rstrip()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
    
    def __shutdown_server(self, signum, frame):
        self._is_running = False    
        if self._server_socket:
            self._server_socket.close()
            self._server_socket = None
            logging.info("action: shutdown_server | result: success")
        if self._last_client_socket:
            self._last_client_socket.close()
            self._last_client_socket = None
            logging.info(f"action: shutdown_client | result: success")
        
        logging.info("action: shutdown | result: success")
    