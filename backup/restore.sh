#!/bin/bash

# Nama container dan database
CONTAINER_NAME=cash-flow-go-db-1
DB_NAME=cashflow
DB_USER=postgres

# File yang ingin direstore
BACKUP_FILE=$1

if [ -z "$BACKUP_FILE" ]; then
  echo "Gunakan: ./restore.sh nama_file_backup.sql"
  exit 1
fi

# Restore
cat "./backup/$BACKUP_FILE" | docker exec -i $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME

echo "Restore selesai dari $BACKUP_FILE"
