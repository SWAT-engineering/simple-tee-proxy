#!/bin/sh

dd if=/dev/urandom of=/tmp/reply bs=1M count=1
cat /tmp/reply | nc -w 5 -L -p 8811 &
backend=$!

"$1" -hosting "localhost:8800" -forward "localhost:8811" &
proxy=$!

sleep 1
echo "Making request"
SUM=$( nc localhost 8800 | md5sum)

echo -e "Gotten: \t $SUM"
EXPECTED_SUM=$(cat /tmp/reply | md5sum )
echo -e "Expected: \t $EXPECTED_SUM"

kill $proxy || true
kill $backend || true


if [ "$SUM" = "$EXPECTED_SUM" ]; then
    exit 0
else
    exit 1
fi