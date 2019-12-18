#!/bin/bash

nc -w 5 -l -p 8811 | md5sum > /tmp/reply-hash &

"$1" -hosting "localhost:8800" -forward "localhost:8811" &
proxy=$!

sleep 1
dd if=/dev/urandom of=/tmp/answer bs=1M count=1 2>/dev/null
cat /tmp/answer | nc localhost 8800

SUM="$(cat /tmp/reply-hash)"
echo -e "Gotten: \t $SUM"
EXPECTED_SUM=$(cat /tmp/answer | md5sum )
echo -e "Expected: \t $EXPECTED_SUM"

kill $proxy || true

if [ "$SUM" = "$EXPECTED_SUM" ]; then
    exit 0
else
    exit 1
fi