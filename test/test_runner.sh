#!/bin/bash

CI=$1

if $CI; then
  cd test
  TEST_BASE=.
else
  TEST_BASE=$GOPATH/src/github.com/zeroFruit/vnet/test
fi

for TEST_DIR in $(find $TEST_BASE -type d -name scenario*)
do
  echo "TEST RUNNING: $TEST_DIR"
  echo ""

  go run $TEST_DIR/main.go
  echo ""
  EXIT_CODE=$?
  if [ $EXIT_CODE -ne 0 ]; then
    echo "RESULT: FAILED"
    exit $EXIT_CODE
  fi
  echo "RESULT: SUCCESS"
  echo "---------------------------------------------------"
done