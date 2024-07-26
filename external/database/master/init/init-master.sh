#!/bin/bash
set -e

# Enable replication and configure necessary settings
cat >> /var/lib/postgresql/data/postgresql.conf <<EOF
wal_level = replica
max_wal_senders = 10
wal_keep_size = 16MB
EOF

# Allow replication user to connect
echo "host replication replica_user 0.0.0.0/0 md5" >> /var/lib/postgresql/data/pg_hba.conf

# Create replication user
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE ROLE replica_user REPLICATION LOGIN PASSWORD 'replica_password';
EOSQL
