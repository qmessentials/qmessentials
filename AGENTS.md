# Agent instructions

## Start here

Read the root `README.md`, then consult only the documentation relevant to the
task:

- `docs/explanation/architecture.md` describes the intended system shape. It
  is direction, not evidence that a capability is implemented.
- `docs/reference/product-requirements.md` describes intended product behavior.
- `docs/architecture/` records accepted architectural decisions.
- `docs/work-tracker-inbox.md` contains unprioritized follow-up work and open
  questions; it is not an implementation specification.

Unless a document explicitly says otherwise, project documentation describes
the desired future state. Treat code, tests, migrations, and deployment
configuration as the source of truth for implementation status. Do not infer
that a capability exists merely because documentation describes it.

## Repository boundaries

- The browser communicates with backend services through `src/apigw`.
- Configuration, intake, and subscription services own separate PostgreSQL
  schemas and migration streams. Do not introduce cross-schema access.
- Add schema changes as new, ordered SQL files under the owning service's
  directory in `src/db-migrations`; do not edit an applied migration casually.
- Keep changes scoped to the requested work and preserve unrelated working-tree
  changes.

## Verification

- Run Go tests from each changed Go module with `go test ./...`.
- For changes under `src/web`, run `npm run lint` and `npm run build` there.
- For SQL changes, run `just lint-sql`.
- Use `git diff --check` before handing work back.

If a relevant command cannot be run, state that explicitly in the handoff.

## Documentation lifecycle

- Keep explanation documents focused on stable concepts and intended system
  direction rather than maintaining a feature-by-feature implementation log.
- Record an accepted, consequential architecture choice as an ADR under
  `docs/architecture`. Use `adr new "Decision title"` when available.
- Put newly discovered but unprioritized work in
  `docs/work-tracker-inbox.md`. Do not treat that file as the authoritative
  backlog, and remove entries after transferring them to the work tracker.
- Avoid copying the same information into multiple documents; link to its
  authoritative location instead.
