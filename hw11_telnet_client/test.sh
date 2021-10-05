#!/usr/bin/env bash
set -xeuo pipefail
NC=/tmp/nc.out
TELNET=/tmp/telnet.out

go build -o go-telnet

(echo -e "Hello\nFrom\nNC\n" && cat 2>/dev/null) | nc -l localhost 4242 >$NC &
NC_PID=$!

sleep 5
(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | ./go-telnet --timeout=5s localhost 4242 >$TELNET &
TL_PID=$!

sleep 5
kill ${TL_PID} 2>/dev/null || true
kill ${NC_PID} 2>/dev/null || true
cat $NC 
cat $TELNET
echo "path=$NC data=$(cat $NC)"
echo "path=$TELNET data=$(cat $TELNET)"
sleep 1

function fileEquals() {
  local fileData
  cat $1
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "unexpected output, $1:\n${fileData}" && exit 1)
}

expected_nc_out='I
am
TELNET client'
fileEquals $NC "${expected_nc_out}"

expected_telnet_out='Hello
From
NC'
fileEquals $TELNET "${expected_telnet_out}"

rm -f go-telnet
echo "PASS"