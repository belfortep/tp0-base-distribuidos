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
        self._clients = []
        self._shutdown = False

        signal.signal(signal.SIGTERM, self.__shutdown_server_handler)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        while self._is_running:
            try:
                client_socket = self.__accept_new_connection()
                if client_socket:
                    self.__handle_client_connection(client_socket)
            except Exception as exception:
                logging.error(f"action: server_exception | result: fail | exception: {exception}")
            finally:
                if not self._shutdown:
                    self.__shutdown_server_handler(None, None)
            

    def __handle_client_connection(self, client_socket):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # TODO: Modify the receive to avoid short-reads
            msg = client_socket.recv(1024).rstrip().decode('utf-8')
            addr = client_socket.getpeername()
            self._clients.append(addr)
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # TODO: Modify the send to avoid short-writes
            client_socket.send("{}\n".format(msg).encode('utf-8'))
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_socket.close()

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
    
    def __shutdown_server_handler(self, signum, frame):
        """
        Closes the server and the client socket
        """
        self._is_running = False
        self._shutdown = True
        if self._server_socket:
            self._server_socket.close()
            self._server_socket = None
            logging.debug("action: shutdown_server_socket | result: success")

        for client in self._clients:
            client.close()
            logging.debug(f"action: shutdown_client | result: success | {client[0]}")
        

    

