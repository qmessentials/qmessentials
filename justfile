set dotenv-load := true

default:
    @just --list

[working-directory("src/web")]
run:
    npm run dev

migrate-configuration:
    python3 src/db-migrations/migrate.py \
        --db-host "$CONFIGURATION_DB_HOST" \
        --db-port "$CONFIGURATION_DB_PORT" \
        --db-name "$CONFIGURATION_DB_NAME" \
        --db-user "$CONFIGURATION_DB_USER" \
        --db-password "$CONFIGURATION_DB_PASSWORD" \
        src/db-migrations/configuration
