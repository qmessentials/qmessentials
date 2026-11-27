# Manage user credentials

The API gateway owns credentials in the `auth` migration stream. Run
`just migrate-all` before starting it. Set `DB_HOST`, `DB_PORT`, `DB_NAME`,
`DB_USER`, and `DB_PASSWORD` for the gateway, or use the corresponding
`AUTH_DB_*` variables with Docker Compose.

A user can change their password by sending their user ID, current password,
and new password to `POST /api/change-password`. The new password must be 12 to
72 bytes. User IDs are case-insensitive and use 3 to 64 letters, digits, dots,
underscores, or hyphens.

`POST /api/login` checks a user ID and password and returns the canonical user
ID on success. It rejects unknown users, wrong passwords, and disabled users
with the same response. The endpoint does not create a session. Session issuance
and access to protected routes are part of the next authentication increment.

Account setup, password reset, and account disabling do not yet have an
operator or administrator interface.
