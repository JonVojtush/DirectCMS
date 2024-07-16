#!/bin/bash
cd ../

LOG_FILE="web/wasmBuild.log"
touch $LOG_FILE 2>&1 | tee /dev/stderr
echo "Script started at: $(date)" | tee -a $LOG_FILE

# Remove go.wasm
if [ -f web/go.wasm ]; then
    echo "web/go.wasm exists, removing it..." | tee -a $LOG_FILE

    if ! rm -f web/go.wasm; then
        echo "Failed to remove web/go.wasm" | tee -a $LOG_FILE >&2
        exit 1
    fi
else
    echo "web/go.wasm doesn't exist." | tee -a $LOG_FILE
fi

# Compile Go code to WebAssembly
echo "Building a new file..." | tee -a $LOG_FILE  
GOOS=js GOARCH=wasm go build -o=web/go.wasm -buildvcs=false 2>&1 | tee -a $LOG_FILE >&2  
echo "web/go.wasm was built." | tee -a $LOG_FILE

# Remove wasm_exec.js
if [ -f web/wasm_exec.js ]; then
    echo "web/wasm_exec.js exists, removing it..." | tee -a $LOG_FILE

    if ! rm -f web/wasm_exec.js; then
        echo "Failed to remove web/wasm_exec.js" | tee -a $LOG_FILE >&2
        exit 1
    fi
else
    echo "web/wasm_exec.js doesn't exist." | tee -a $LOG_FILE  
fi

# Copy wasm_exec.js from GOROOT/misc/wasm to the current working directory's web subdirectory
echo "Fetching a new file..." | tee -a $LOG_FILE 
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" web/ 2>&1 | tee -a $LOG_FILE >&2  
echo "web/wasm_exec.js was fetched from \$GOROOT/misc/wasm/" | tee -a $LOG_FILE

echo "Script ended at: $(date)" | tee -a $LOG_FILE
echo "----------------------------------------------------------------------------------------------------" | tee -a $LOG_FILE

exit