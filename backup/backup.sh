#!/bin/bash

# Variabel
CONTAINER_NAME=cash-flow-go-db-1
DB_NAME=cashflow
DB_USER=postgres
BACKUP_DIR=./backup
RETENTION_DAYS=7

# Timestamp & nama file
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="db_backup_$TIMESTAMP.sql"

# Buat folder backup kalau belum ada
mkdir -p "$BACKUP_DIR"

# Eksekusi backup
docker exec -t $CONTAINER_NAME pg_dump -U $DB_USER -d $DB_NAME > "$BACKUP_DIR/$BACKUP_FILE"

# Cek apakah berhasil
if [ $? -eq 0 ]; then
  echo "✅ Backup berhasil: $BACKUP_FILE"
else
  echo "❌ Gagal backup!"
  exit 1
fi

# Auto hapus file backup lebih dari RETENTION_DAYS
find "$BACKUP_DIR" -name "db_backup_*.sql" -type f -mtime +$RETENTION_DAYS -exec rm -f {} \;

echo "🧹 Backup lama lebih dari $RETENTION_DAYS hari telah dihapus"
