#!/bin/bash
set -e # Skrip akan berhenti jika ada error

# 'psql' adalah command-line tool untuk Postgres
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    
    -- Membuat database untuk setiap service
    CREATE DATABASE auth_db;
    CREATE DATABASE product_db;
    CREATE DATABASE supplier_db;
    CREATE DATABASE inventory_db;
    CREATE DATABASE po_db;
    
    -- (Opsional) Beri hak akses ke user-mu
    GRANT ALL PRIVILEGES ON DATABASE auth_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE product_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE supplier_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE inventory_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE po_db TO $POSTGRES_USER;

EOSQL

echo "Database berhasil dibuat."