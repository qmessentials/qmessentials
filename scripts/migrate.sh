#!/bin/bash
set -e

# Load environment and run migrations
cd /app
just migrate-all
