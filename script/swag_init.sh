#! /bin/sh
# Swagger need pgx types for generation
swag init --parseDependency --parseInternal -d cmd/server/,./internal/handlers/
