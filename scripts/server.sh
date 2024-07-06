#!/bin/bash
cd ../web || { echo "Failed to cd into ../web"; exit 1; }

python -m http.server 8080 &
PID=$!

# Wait a bit for the server to start
for i in {1..5}; do
    if ps -p "$PID" >/dev/null; then
        echo "Success in $i seconds... Access at $(pwd)"
        exit 0
    fi
    sleep 2
done

echo "[Error] Could not start Python HTTP server." 1>&2
exit 1