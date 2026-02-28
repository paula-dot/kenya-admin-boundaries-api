#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

PID_FILE="server.pid"
OUT_FILE="server.out"
BIN_FILE="bin/api"

# Stop existing server if pid file exists
if [ -f "$PID_FILE" ]; then
  PID=$(cat "$PID_FILE" 2>/dev/null || true)
  if [ -n "$PID" ]; then
    if kill -0 "$PID" 2>/dev/null; then
      echo "Stopping existing server (PID $PID) with SIGTERM..."
      kill "$PID" || true
      for i in {1..5}; do
        if kill -0 "$PID" 2>/dev/null; then
          sleep 1
        else
          break
        fi
      done
      if kill -0 "$PID" 2>/dev/null; then
        echo "PID $PID still alive, sending SIGKILL..."
        kill -9 "$PID" || true
        sleep 1
      else
        echo "PID $PID stopped."
      fi
    else
      echo "PID $PID not running, removing pid file."
    fi
  else
    echo "No PID found in $PID_FILE, removing file.";
  fi
  rm -f "$PID_FILE"
else
  echo "No $PID_FILE file present.";
fi

# Ensure bin directory exists
mkdir -p bin

# Build the binary if missing
if [ -x "$BIN_FILE" ]; then
  echo "Found existing executable $BIN_FILE, skipping build."
else
  if [ -f "$BIN_FILE" ]; then
    echo "$BIN_FILE exists but is not executable. Making executable."
    chmod +x "$BIN_FILE" || true
  else
    echo "Building $BIN_FILE from ./cmd/api..."
    if go build -o "$BIN_FILE" ./cmd/api; then
      echo "Built $BIN_FILE"
    else
      echo "go build failed" >&2
      exit 1
    fi
  fi
fi

# Start server in background, capture stdout/stderr to server.out
echo "Starting server: nohup $BIN_FILE > $OUT_FILE 2>&1 &"
nohup "$BIN_FILE" > "$OUT_FILE" 2>&1 &
NEWPID=$!
# Give the process a moment to settle
sleep 1
# record PID
echo $NEWPID > "$PID_FILE"

echo "Started with PID $NEWPID"

# Show a short tail of the log
if [ -f "$OUT_FILE" ]; then
  echo "--- $OUT_FILE (last 40 lines) ---"
  tail -n 40 "$OUT_FILE" || true
else
  echo "No log file $OUT_FILE yet.";
fi

exit 0

