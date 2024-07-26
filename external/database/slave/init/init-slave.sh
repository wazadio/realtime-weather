#!/bin/bash
set -e

# Stop PostgreSQL service
pg_ctl -D "$PGDATA" -m fast -w stop

# Clear existing data
rm -rf "$PGDATA"/*

# Base backup from master
PGPASSWORD=replica_password pg_basebackup -h $REPLICATE_FROM -D "$PGDATA" -U replica_user -v -P --wal-method=stream

# Create the standby signal file
touch "$PGDATA/standby.signal"

# Set primary connection info
cat >> "$PGDATA/postgresql.conf" <<EOF
primary_conninfo = 'host=$REPLICATE_FROM port=5432 user=replica_user password=replica_password'
EOF

# Set appropriate permissions
chown -R postgres:postgres "$PGDATA"
chmod 700 "$PGDATA"

# Start PostgreSQL service
pg_ctl -D "$PGDATA" -w start
