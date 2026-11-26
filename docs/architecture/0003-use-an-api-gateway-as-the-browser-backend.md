# 3. Use an API gateway as the browser backend

Date: 2026-09-02

## Status

Accepted

## Context

The browser needs a stable backend boundary while domain capabilities are
split among several services. Exposing every service directly would distribute
browser authentication, routing, and cross-origin policy across those
services.

## Decision

The browser will call a dedicated API gateway. The gateway will own the
browser-facing authentication boundary and proxy requests to the appropriate
backend service.

## Consequences

The browser has one backend origin and individual services do not need to
implement browser-specific authentication and CORS independently. Internal
service boundaries remain hidden from the UI.

The gateway is an additional deployed component and can become a bottleneck or
coupling point. The current cookie login and shared internal secret are
development mechanisms and must be replaced or hardened for production.
