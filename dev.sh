#!/bin/sh

# Watch for file changes and rebuild container
while true; do
  echo "Starting application..."
  docker-compose down
  docker-compose build --no-cache
  docker-compose up -d
  
  echo "Watching for changes (Ctrl+C to stop)..."
  # Wait for any file changes (uses find to check timestamps)
  find . -type f -name "*.go" | entr -d echo "Change detected, rebuilding..."
done
