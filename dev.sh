#!/bin/bash

# import environment variables
set -a; source .env; set +a;

go run main.go
