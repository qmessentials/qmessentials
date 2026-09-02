# Run QMEssentials locally

## Prerequisites

- Docker with Docker Compose
- Go, for running or testing services outside containers
- Node.js and npm, for the web application
- Python 3, for SQL linting
- [`just`](https://github.com/casey/just), for repository commands

## Configure the environment

The Compose file supplies development defaults for most service settings.
Create or update the root `.env` for values specific to your machine.
The web application reads its development settings from `src/web/.env.development`.

Do not commit credentials or production secrets.

## Start backend services

From the repository root, run:

```sh
just run-services
```

This builds and starts PostgreSQL, NATS, database migrations, the API gateway,
and the backend services defined in `docker-compose.yml`.

The API gateway listens on `http://localhost:8080`.
PostgreSQL is exposed to the host on port `5433`,
and NATS exposes its client and monitoring ports on `4222` and `8222` respectively.

## Start the web application

Install its dependencies once:

```sh
cd src/web
npm install
cd ../..
```

Then run:

```sh
just run-web
```

Vite prints the browser URL when it starts.

## Run database migrations manually

With the corresponding database environment variables set,
run one migration stream or all streams:

```sh
just migrate-configuration
just migrate-intake
just migrate-subscription
just migrate-all
```

Compose normally runs migrations through its `migrate` service.

## Lint SQL

Create a virtual environment and install the development dependency:

```sh
python3 -m venv .venv
.venv/bin/pip install -r requirements-dev.txt
```

Then run:

```sh
just lint-sql
```

`just fix-sql` applies SQLFluff's safe automatic formatting fixes.
