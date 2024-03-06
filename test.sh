#!/bin/bash

export APP_DEBUG="true"
export DB_USER="root"
export DB_PASSWORD="root"
export DB_HOST="localhost"
export DB_PORT="3306"
export DB_DATABASE="test_db"

go mod tidy
go test -v $1