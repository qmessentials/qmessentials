# Manage user credentials

The API gateway owns credentials in the `auth` migration stream. Run
`just migrate-all` before starting it. Set `DB_HOST`, `DB_PORT`, `DB_NAME`,
`DB_USER`, and `DB_PASSWORD` for the gateway, or use the corresponding
`AUTH_DB_*` variables with Docker Compose.

A user can change their password by sending their user ID, current password,
and new password to `POST /api/change-password`. The new password must be 12 to
72 bytes. User IDs are case-insensitive and use 3 to 64 letters, digits, dots,
underscores, or hyphens.

`POST /api/login` checks a user ID and password, returns the canonical user ID
and expiry, and sets an HttpOnly session cookie. It rejects unknown users, wrong
passwords, and disabled users with the same response. Protected gateway routes
require that cookie. Sessions expire after 12 hours and do not renew
automatically. `POST /api/session/renew` rotates a valid session token and starts
a new 12-hour period. `GET /api/session` returns the current user ID.
`POST /api/logout` revokes the current session and clears its cookie. Expired or
revoked sessions cannot be renewed or used for protected requests.

Account setup, password reset, and account disabling do not yet have an
operator or administrator interface.
