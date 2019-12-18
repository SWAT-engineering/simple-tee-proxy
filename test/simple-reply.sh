#!/bin/bash

dd if=/dev/urandom of=/tmp/reply bs=1M count=1 2>/dev/null
cat /tmp/reply | nc -w 5 -l -p 8811 &
EXPECTED_SUM=$(cat /tmp/reply | md5sum )

"$1" -hosting "localhost:8800" -forward "localhost:8811" &
proxy=$!

sleep 1
SUM=$( nc localhost 8800 | md5sum)

echo -e "Expected: \t $EXPECTED_SUM"
echo -e "Gotten: \t $SUM"

kill $proxy || true


if [ "$SUM" = "$EXPECTED_SUM" ]; then
    exit 0
else
    exit 1
fi
