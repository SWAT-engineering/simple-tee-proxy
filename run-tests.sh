#!/bin/bash
RESULT=0

PROXY_PROGRAM="$1"

function runTest() {
    printf "🧪 \tRunning: $1 \n"
    bash "$2" "$PROXY_PROGRAM"
    if [ $? -eq 0 ]; then
        printf "🎉 \tSuccess  \n"
    else
        printf "🔥 \tFailed  \n"
        RESULT=1
    fi
}

runTest "simple reply" "test/simple-reply.sh"
runTest "simple answer" "test/simple-answer.sh"

exit $RESULT