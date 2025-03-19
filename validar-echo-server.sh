#!/bin/bash

RESPONSE=$(docker run --rm --network="tp0_testing_net" gophernet/netcat sh -c "echo TEST | nc server 12345")


echo $RESPONSE