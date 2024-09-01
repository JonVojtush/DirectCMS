#!/bin/bash
LOG_FILE="scripts/wasmBuild.log"
touch $LOG_FILE 2>&1 | tee /dev/stderr
echo "Script started at: $(date)" | tee -a $LOG_FILE

# Remove go.wasm if it exists
echo "Removing web/go.wasm if it exists..." | tee -a $LOG_FILE
rm -f web/go.wasm 2>&1 | tee -a $LOG_FILE >&2

# Compile Go code to WebAssembly
echo "Building a new file..." | tee -a $LOG_FILE
GOOS=js GOARCH=wasm go build -o=web/go.wasm -buildvcs=false main.go 2>&1 | tee -a $LOG_FILE >&2
status=$?
if [ $status -ne 0 ]; then
    echo "Failed to build web/go.wasm" | tee -a $LOG_FILE >&2
    exit 1
fi
echo "web/go.wasm was built." | tee -a $LOG_FILE

# Remove wasm_exec.js if it exists
echo "Removing web/wasm_exec.js if it exists..." | tee -a $LOG_FILE
rm -f web/wasm_exec.js 2>&1 | tee -a $LOG_FILE >&2

# Copy wasm_exec.js from GOROOT/misc/wasm to the current working directory's web subdirectory
echo "Fetching a new file..." | tee -a $LOG_FILE 
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" web/ 2>&1 | tee -a $LOG_FILE >&2  
echo "web/wasm_exec.js was fetched from \$GOROOT/misc/wasm/" | tee -a $LOG_FILE

echo "Script ended at: $(date)" | tee -a $LOG_FILE
echo "----------------------------------------------------------------------------------------------------" | tee -a $LOG_FILE

exit
