#!/bin/bash
cd ../;

# Remove build resources
rm web/go.wasm web/wasm_exec.js build.log;
echo "Build resources cleaned";

# Check if first argument is "go", then clean go environment
if [ "$1" == "go" ]; then 
    go clean --modcache;
    echo "Go environment cleaned.";
fi

exit