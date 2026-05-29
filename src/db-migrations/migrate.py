#!/usr/bin/env python3

import argparse
import re
import sys
from pathlib import Path

import psycopg


MIGRATION_FILE_PATTERN = re.compile(r"^(\d+).*\.sql$")


def build_connection_kwargs(args: argparse.Namespace) -> dict:
    kwargs = {
        "host": args.db_host,
        "port": args.db_port,
        "dbname": args.db_name,
        "user": args.db_user,
        "password": args.db_password,
    }

    return kwargs


def discover_migrations(migration_dir: Path) -> list[tuple[int, Path]]:
    if not migration_dir.exists():
        raise RuntimeError(f"Migration directory does not exist: {migration_dir}")

    if not migration_dir.is_dir():
        raise RuntimeError(f"Migration path is not a directory: {migration_dir}")

    migrations: list[tuple[int, Path]] = []

    for path in migration_dir.iterdir():
        if not path.is_file():
            continue

        match = MIGRATION_FILE_PATTERN.match(path.name)
        if not match:
            continue

        version = int(match.group(1))
        migrations.append((version, path))

    migrations.sort(key=lambda item: item[0])

    return migrations


def get_current_schema_version(conn: psycopg.Connection) -> int:
    with conn.cursor() as cur:
        cur.execute("select max(version_number) from schema_migrations")
        result = cur.fetchone()

    if result is None or result[0] is None:
        raise RuntimeError("No current schema version found in schema_migrations")

    return int(result[0])


def apply_migration(
    conn: psycopg.Connection,
    version: int,
    migration_path: Path,
) -> None:
    sql = migration_path.read_text(encoding="utf-8")

    print(f"Applying migration {version}: {migration_path.name}")

    with conn.transaction():
        with conn.cursor() as cur:
            cur.execute(sql)

            cur.execute(
                """
                insert into schema_migrations (version_number)
                values (%s)
                """,
                (version,),
            )


def run(args: argparse.Namespace) -> int:
    migrations = discover_migrations(args.migration_dir)
    connection_kwargs = build_connection_kwargs(args)

    with psycopg.connect(**connection_kwargs) as conn:
        current_version = get_current_schema_version(conn)

        pending = [
            (version, path)
            for version, path in migrations
            if version > current_version
        ]

        print(f"Current schema version: {current_version}")

        if not pending:
            print("No pending migrations.")
            return 0

        print("Pending migrations:")
        for version, path in pending:
            print(f"  {version}: {path.name}")

        for version, path in pending:
            apply_migration(conn, version, path)

        print("Migrations complete.")

    return 0


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Apply PostgreSQL schema migrations."
    )

    parser.add_argument(
        "migration_dir",
        type=Path,
        help="Path to the directory containing SQL migration scripts.",
    )

    parser.add_argument(
        "--db-host",
        required=True,
        help="PostgreSQL host.",
    )

    parser.add_argument(
        "--db-port",
        type=int,
        default=5432,
        help="PostgreSQL port.",
    )

    parser.add_argument(
        "--db-name",
        required=True,
        help="PostgreSQL database name.",
    )

    parser.add_argument(
        "--db-user",
        required=True,
        help="PostgreSQL user.",
    )

    parser.add_argument(
        "--db-password",
        required=True,
        help="PostgreSQL password.",
    )

    parser.add_argument(
        "--db-sslmode",
        default=None,
        help="Optional PostgreSQL sslmode value.",
    )

    return parser.parse_args()


def main() -> int:
    args = parse_args()

    try:
        return run(args)
    except Exception as error:
        print(f"Migration failed: {error}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
