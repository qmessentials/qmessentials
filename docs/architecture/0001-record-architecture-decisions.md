# 1. Record architecture decisions

Date: 2026-09-02

## Status

Accepted

## Context

QMEssentials has consequential technical choices embedded in code and mixed
into general project documentation. Without a durable decision log, it is hard
to distinguish an accepted constraint from a current implementation detail or
an uncommitted idea.

## Decision

We will record accepted architectural decisions as Architecture Decision
Records under `docs/architecture`. The repository will use `adr-tools`, with
the ADR directory configured by the root `.adr-dir` file.

## Consequences

ADRs preserve the context and tradeoffs behind decisions. Current-system
documentation may link to them instead of repeating their rationale.

Writing and maintaining ADRs adds a small documentation cost. Proposed work
and open design questions must remain outside the ADR log until a decision is
accepted.
