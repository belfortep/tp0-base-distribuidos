#!/bin/bash
echo "Output filename: $1"
echo "Number of clients: $2"
python3 ./compose-generator/compose_generator.py $1 $2