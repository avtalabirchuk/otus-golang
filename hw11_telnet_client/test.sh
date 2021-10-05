#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-telnet

(echo -e "Hello" && cat 2>/dev/null) | nc -l 127.0.0.1 4242 >./nc.out &
NC_PID=$!
ip a

sleep 5
(echo -e "I" && cat 2>/dev/null) | ./go-telnet --timeout=5s 127.0.0.1 4242 >./telnet.out &
TL_PID=$!
ip a

sleep 5
kill ${TL_PID} 2>/dev/null || true
kill ${NC_PID} 2>/dev/null || true
cat ./telnet.out
cat ./nc.out


function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "unexpected output, $1:\n${fileData}" && exit 1)
}

expected_nc_out='I'
fileEquals ./nc.out "${expected_nc_out}"

expected_telnet_out='Hello'
fileEquals ./telnet.out "${expected_telnet_out}"

rm -f go-telnet
echo "PASS"