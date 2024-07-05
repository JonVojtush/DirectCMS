#!/bin/bash
cd ../web || exit 1;

python3 -m http.server "8080" & PID=$!;

sleep 2; # Make sure http-server has enough time to start.

if ps -p "$PID" >/dev/null; then
    echo "Success. Access at $(pwd)";
else
    echo "[Error] Could not start Python HTTP server." 1>&2;
fi

exit