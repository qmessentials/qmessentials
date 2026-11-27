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
a new 12-hour period. `GET /api/session` returns the current user ID and roles.
`POST /api/logout` revokes the current session and clears its cookie. Expired or
revoked sessions cannot be renewed or used for protected requests.

To establish the first Administrator, set `BOOTSTRAP_ADMIN_USER_ID` and
`BOOTSTRAP_ADMIN_PASSWORD` before starting the gateway. The password must be
12 to 72 bytes. The gateway creates this account only when no Administrator
exists. Remove the variables after the first successful start.

An Administrator can use the Users screen to create users, assign the
Administrator, Quality Manager, and Technician roles, and disable or re-enable
accounts. The gateway also exposes `GET /api/admin/users`,
`POST /api/admin/users`, `PUT /api/admin/users/:userId/roles`, and
`PUT /api/admin/users/:userId/disabled` for this screen. These endpoints require
an active Administrator session. The last enabled Administrator cannot be
disabled or lose the Administrator role.
Role changes are checked on each subsequent management request. Disabled users
cannot log in or use existing sessions.
