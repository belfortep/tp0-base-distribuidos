import sys



def main():
    if len(sys.argv) != 3:
        print("Use: python3 compose_generator.py <output filename> <number of clients>")
        sys.exit(-1)

    output_filename = sys.argv[1]
    try:
        number_of_clients = int(sys.argv[2])
    except ValueError:
        print("Number of clients must be an integer")
        sys.exit(-1)

    print(f"filename: {output_filename}")
    print(f"number: {number_of_clients}")



if __name__ == "__main__":
    main()