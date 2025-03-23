import socket
import logging
import signal
from common.utils import Bet,store_bets
from common.connection import send, read_up_to_delimiter

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
            bet = self.__read_bet()
            store_bets([bet])
            logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")
            send(self._last_client_socket, "ACK")
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            self._last_client_socket.close()
            self._last_client_socket = None
    
    def __read_bet(self):   
        message = read_up_to_delimiter(self._last_client_socket, "\0")
        logging.info(f"action: reading_bet | result: success | message: {message}")
        data_list = message.split(";")
        if len(data_list) != 6:
            raise Exception("Corrupted message read")
        
        agency = data_list[0]
        name = data_list[1]
        surname = data_list[2]
        dni = data_list[3]
        birthdate = data_list[4]
        number = data_list[5]

        if not agency.isdigit() or not number.isdigit():
            raise Exception("Wrong data types in number or agency")
        
        return Bet(agency, name, surname, dni, birthdate, number)
        
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
    