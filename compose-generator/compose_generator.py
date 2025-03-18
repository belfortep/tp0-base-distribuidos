import sys

def generate_docker_compose(number_of_clients):
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
    networks:
      - testing_net
    depends_on:
      - server
    volumes:
      - ./client/config.yaml:/config.yaml
    """

    for i in range(1, number_of_clients + 1):
        template += client_template.format(id=i)

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
        print("Number of clients must be an integer bigger than 0")
        sys.exit(-1)

    if number_of_clients <= 0:
        print("Must have a number of clients bigger than 0") 
        sys.exit(-1)
    
    create_docker_compose_file(output_filename, number_of_clients)

if __name__ == "__main__":
    main()