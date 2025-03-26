import socket
import logging
import signal
from common.utils import Bet,store_bets, load_bets, has_won
from common.connection import send, read_up_to_delimiter

class Server:
    def __init__(self, port, listen_backlog, number_of_clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._is_running = True
        self._last_client_socket = None
        self._completed_agencies = set()
        self._number_of_clients = number_of_clients

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
            logging.error(f"action: server_run | result: fail | error: {e}")
        finally:
            self.__shutdown_server(None, None)
            

    def __handle_client_connection(self):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            message = read_up_to_delimiter(self._last_client_socket, "\0")
            if message.startswith("GETWINNERS"):
                message = self.__get_winners(message)
                send(self._last_client_socket, message)
            else:
                bets, errors = self.__get_bets(message)
                store_bets(bets)
                
                if errors > 0:
                    logging.error(f"action: apuesta_recibida | result: fail  | cantidad: {errors}")
                    send(self._last_client_socket, f"ERR;{errors}")
                else:
                    logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
                    send(self._last_client_socket, "ACK")
        except ConnectionResetError as e:
            logging.info(f"action: server_run | result: success | message: the socket is now closed")
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            self._last_client_socket.close()
            self._last_client_socket = None
    
    def __get_winners(self, message): 
        
        values = message.split(";")
        logging.info(f"action: trying to get winners | result: success | agency: {values[1]}")
        self._completed_agencies.add(values[1])
        message = "WINNERS;"
        if len(self._completed_agencies) == self._number_of_clients:
            
            for bet in load_bets():
                if has_won(bet):
                    if bet.agency == int(values[1]):
                        message = message + bet.document + ";"
            logging.info(f"action: sorteo | result: success")
            message = message[:-1]
            return message
        else:
            logging.info(f"action: not_yet_winner | result: success")
            return "NOTYET"

    def __get_bets(self, message):   
        errors = 0
        bets = []
        for bet_message in message.split("\n"):
            data_list = bet_message.split(";")
            if len(data_list) != 6:
                errors += 1
                continue
        
            agency = data_list[0]
            name = data_list[1]
            surname = data_list[2]  
            dni = data_list[3]
            birthdate = data_list[4]
            number = data_list[5]

            if not agency.isdigit() or not number.isdigit():
                errors += 1
                continue

            bet = Bet(agency, name, surname, dni, birthdate, number)
            bets.append(bet)
        
        return bets, errors
    
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
    