#!/usr/bin/env bash

set -e

echo "🔧 Initializing PostgreSQL database and user..."

psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" <<-EOSQL
    CREATE USER ${APP_DB_USER} WITH PASSWORD '${APP_DB_PASSWORD}';
    CREATE DATABASE ${APP_DB_NAME};
    \c ${APP_DB_NAME}
    GRANT ALL PRIVILEGES ON DATABASE ${APP_DB_NAME} TO ${APP_DB_USER};
    GRANT ALL ON SCHEMA public TO ${APP_DB_USER};
EOSQL

echo "✅ PostgreSQL initialization complete!"
echo "   Database: ${APP_DB_NAME}"
echo "   User: ${APP_DB_USER}"
