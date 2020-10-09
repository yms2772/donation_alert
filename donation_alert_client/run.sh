#!/bin/bash
echo "Builing wasm..."
GOARCH=wasm GOOS=js go build -v -o ../donation_alert/web/app.wasm
