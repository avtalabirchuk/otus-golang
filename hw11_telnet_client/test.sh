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
kill ${TL_PID} 2>/dev/null && \
kill ${NC_PID} 2>/dev/null
#echo 123 >/tmp/123 && cat /tmp/123
# cat $NC 
# cat $TELNET
# echo "path=$NC data=$(cat $NC)"
# echo "path=$TELNET data=$(cat $TELNET)"

function fileEquals() {
  local fileData
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