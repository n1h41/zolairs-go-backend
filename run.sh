#!/bin/bash

# Build and run the application

echo "Building Zolaris Backend..."
go build -o zolaris-backend

if [ $? -eq 0 ]; then
    echo "Build successful. Starting application..."
    ./zolaris-backend
else
    echo "Build failed."
    exit 1
fi

