#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-telnet

(echo -e "Hello" && cat 2>/dev/null) | nc -l localhost 4242 >/tmp/nc.out &
NC_PID=$!

#sleep 2
(echo -e "I" && cat 2>/dev/null) | ./go-telnet --timeout=5s localhost 4242 >/tmp/telnet.out &
TL_PID=$!

sleep 5
kill ${TL_PID} 2>/dev/null || true
kill ${NC_PID} 2>/dev/null || true

function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "unexpected output, $1:\n${fileData}" && exit 1)
}

expected_nc_out='I'
fileEquals /tmp/nc.out "${expected_nc_out}"

expected_telnet_out='Hello'
fileEquals /tmp/telnet.out "${expected_telnet_out}"

rm -f go-telnet
echo "PASS"