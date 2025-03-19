#!/bin/bash

RESPONSE=$(docker run --rm --network="tp0_testing_net" alpine sh -c "echo TEST | nc server 12345")

if [ "$RESPONSE" == "TEST" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi
