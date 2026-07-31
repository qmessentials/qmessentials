# QMEssentials: Easy Quality Management for SMBs Who Make Things

QMEssentials contains a set of tools
that enable small and medium-sized manufacturers
to track the quality of their products.

## SQL linting

Install the development dependency in a virtual environment:

```sh
python3 -m venv .venv
.venv/bin/pip install -r requirements-dev.txt
```

Lint PostgreSQL migrations and test data with `just lint-sql`. Run
`just fix-sql` to apply SQLFluff's safe automatic formatting fixes.
