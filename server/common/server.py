import socket
import logging
import signal
from multiprocessing import Process, Barrier, Lock
from common.utils import Bet,store_bets, load_bets, has_won
from common.connection import send, read_up_to_delimiter

class Server:
    def __init__(self, port, listen_backlog, number_of_clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._is_running = True
        self._clients_sockets = []
        self._barrier = Barrier(number_of_clients)
        
        
        self._bet_lock = Lock()

        signal.signal(signal.SIGTERM, self.__shutdown_server)

    def run(self):
        """
        Server loop

        Server that accept a new connections and establishes a
        communication with a client. Creates a new process for every new Client
        """
        client_processes = []
        try:
            while self._is_running:
                client_socket = self.__accept_new_connection()
                if client_socket:
                    self._clients_sockets.append(client_socket)
                    new_process = Process(target=self.__handle_client_connection, args=(client_socket, ))
                    client_processes.append(new_process)
                    new_process.start()
            
            for process in client_processes:
                process.join()
                logging.info(f"action: joining process | result: success")

        except Exception as e:
            logging.error(f"action: server_run | result: fail | error: {e}")
        finally:
            self.__shutdown_server(None, None)
            

    def __handle_client_connection(self, client_socket):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        is_client_running = True
        try:
            while is_client_running:
                message = read_up_to_delimiter(client_socket, "\0")
                if message.startswith("GETWINNERS"):
                    logging.info(f"action: waiting_barrier | result: success")
                    self._barrier.wait()
                    logging.info(f"action: end in barrier! | result: success")
                    message = self.__get_winners(message)
                    send(client_socket, message)
                    is_client_running = False
                else:
                    bets, errors = self.__get_bets(message)
                    with self._bet_lock:
                        store_bets(bets)
                    
                    if errors > 0:
                        logging.error(f"action: apuesta_recibida | result: fail  | cantidad: {errors}")
                        send(client_socket, f"ERR;{errors}")
                    else:
                        logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
                        send(client_socket, "ACK")
        except ConnectionResetError as e:
            logging.info(f"action: server_run | result: success | message: the socket is now closed")
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            logging.info(f"action: closing client | result: success ")
            client_socket.close()
    
    def __get_winners(self, message): 
        """
        Get the Winners bets from a Message in the format WINNERS;AGENCY_NUMBER
        from the bets stored
        """
        values = message.split(";")
        message = "WINNERS;"
        for bet in load_bets():
            if has_won(bet):
                if bet.agency == int(values[1]):
                    message = message + bet.document + ";"
        message = message[:-1]
        logging.info(f"action: sorteo | result: success")
        return message

    def __get_bets(self, message):   
        """
        Get the bets Message in the format bet1\nbet2\n...betn
        returns an array with all the bets and the number of errors found while parsing the bets
        """
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
        try:
            logging.info('action: accept_connections | result: in_progress')
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return c
        except:
            logging.info("Server closed")
            return None
            
    def __shutdown_server(self, signum, frame):
        """
        Gracefully shutdown the server and the client connections
        """
        self._is_running = False    
        if self._server_socket:
            self._server_socket.close()
            self._server_socket = None
            logging.info("action: shutdown_server | result: success")
        for client_socket in self._clients_sockets:
            client_socket.close()
            logging.info(f"action: shutdown_client | result: success")
        
        logging.info("action: shutdown | result: success")
    