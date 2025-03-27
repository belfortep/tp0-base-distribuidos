#!/bin/bash

TEST_MESSAGE="De aquel amor, De musica ligera..."
PORT="12345"
RESPONSE=$(docker run --rm --network="tp0_testing_net" alpine sh -c "echo $TEST_MESSAGE | nc server $PORT")

if [ "$RESPONSE" = "$TEST_MESSAGE" ]; then
    RESULT="success"
else
    RESULT="fail"
fi

echo "action: test_echo_server | result: $RESULT"