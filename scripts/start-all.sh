# node
make init &

# rest
sleep 5
bin/lycli rest-server --chain-id longychain --trust-node --laddr "tcp://0.0.0.0:1317" &

# key service
sleep 5
bin/ks --aws-dynamo-url "http://dynamodb:8000/"