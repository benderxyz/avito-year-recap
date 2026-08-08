# Infra

Infrastructure-related files: Docker notes, compose extensions, and deployment configuration.

For now, service Dockerfiles stay inside each service directory because `docker-compose.yml` builds services from those contexts.

`postgres/init/` creates databases `cards` and `users` on first Postgres start. Both services share one Postgres container and keep data in separate databases.
