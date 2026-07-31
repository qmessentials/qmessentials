set dotenv-load := true

default:
    @just --list

run-services:
    docker compose up --build

[working-directory("src/web")]
run-web:
    npm run dev

[working-directory("src/utils")]
migrate-configuration:
    go run . migrate-db \
        --host "$CONFIGURATION_DB_HOST" \
        --port "$CONFIGURATION_DB_PORT" \
        --database "$CONFIGURATION_DB_NAME" \
        --user "$CONFIGURATION_DB_USER" \
        --password "$CONFIGURATION_DB_PASSWORD" \
        --path ../db-migrations/configuration

[working-directory("src/utils")]
migrate-intake:
    go run . migrate-db \
        --host "$INTAKE_DB_HOST" \
        --port "$INTAKE_DB_PORT" \
        --database "$INTAKE_DB_NAME" \
        --user "$INTAKE_DB_USER" \
        --password "$INTAKE_DB_PASSWORD" \
        --path ../db-migrations/intake

[working-directory("src/utils")]
migrate-subscription:
    go run . migrate-db \
        --host "$SUBSCRIPTION_DB_HOST" \
        --port "$SUBSCRIPTION_DB_PORT" \
        --database "$SUBSCRIPTION_DB_NAME" \
        --user "$SUBSCRIPTION_DB_USER" \
        --password "$SUBSCRIPTION_DB_PASSWORD" \
        --path ../db-migrations/subscription

migrate-all:
    just migrate-configuration
    just migrate-intake
    just migrate-subscription

lint-sql:
    sqlfluff lint $(git ls-files '*.sql')

fix-sql:
    sqlfluff fix $(git ls-files '*.sql')
