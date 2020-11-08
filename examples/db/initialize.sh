#!/bin/sh

cd /usr/pg-mig
./pg-mig init -name="$POSTGRES_DB" -credentials="$POSTGRES_USER":"$POSTGRES_PASSWORD" -path=./workspace -db=""