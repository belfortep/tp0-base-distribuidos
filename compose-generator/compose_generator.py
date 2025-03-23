import sys

def generate_docker_compose(number_of_clients):
    

    names = ["ALAN", "DAN", "ADELE", "LINUS", "GUIDO"]
    surnames = ["KAY", "INGALLS", "GOLDBERG", "TORVALDS", "VAN ROSSUM"]
    dnis = ["123", "456", "789", "012", "345"]
    birthdates = ["1940-05-17", "1944-01-01", "1945-07-07", "1969-12-28", "1956-01-31"]
    numbers = ["7574", "1234", "5678", "9012", "3456"]



    template = """
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
    networks:
      - testing_net
    volumes:
      - ./server/config.ini:/config.ini
    """

    client_template = """
  client{id}:
    container_name: client{id}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={id}
      - CLI_NOMBRE={name}
      - CLI_APELLIDO={surname}
      - CLI_DOCUMENTO={dni}
      - CLI_NACIMIENTO={birthdate}
      - CLI_NUMERO={number}
    networks:
      - testing_net
    depends_on:
      - server
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/agency-{id}.csv:/agency-{id}.csv
    """

    for i in range(1, number_of_clients + 1):
        template += client_template.format(id=i, name=names[i-1], surname=surnames[i-1], dni=dnis[i-1], birthdate=birthdates[i-1], number=numbers[i-1])

    template += """
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
    """

    return template
    
def create_docker_compose_file(output_filename, number_of_clients):
    
    docker_compose_data = generate_docker_compose(number_of_clients)

    with open(output_filename, "w") as output_file:
        output_file.write(docker_compose_data.strip())


def main():
    if len(sys.argv) != 3:
        print("Use: python3 compose_generator.py <output filename> <number of clients>")
        sys.exit(-1)

    output_filename = sys.argv[1]
    try:
        number_of_clients = int(sys.argv[2])
    except ValueError:
        print("Number of clients must be a positive integer ")
        sys.exit(-1)

    
    create_docker_compose_file(output_filename, number_of_clients)

if __name__ == "__main__":
    main()