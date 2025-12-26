MIG_PATH ?= 'cmd/migration/scripts'
DB_POSTGRES_URL ?= localhost
DB_POSTGRES_PORT ?= 5432
DB_POSTGRES_USER ?= postgres
DB_POSTGRES_PASSWORD ?= postgres
DB_POSTGRES_DATABASE ?= otp
DB_POSTGRES_SSL_MODE ?= disable
DB_POSTGRES_MULTI_STATEMENT = false

DB_PG_URL ?= postgres://$(DB_POSTGRES_USER):$(DB_POSTGRES_PASSWORD)@$(DB_POSTGRES_URL)$(DB_POSTGRES_HOST):$(DB_POSTGRES_PORT)/$(DB_POSTGRES_DATABASE)?sslmode=$(DB_POSTGRES_SSL_MODE)\&x-multi-statement=$(DB_POSTGRES_MULTI_STATEMENT)

mig-install:
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Use it like make mig-goto ARRGS=12
mig-goto:
	@migrate -database $(DB_PG_URL) -path $(MIG_PATH) goto $(ARGS)

# Use it like
# Fully migrate:
# make mig-up
# Migrate from a vertion:
# make mig-up ARRGS=12
mig-up:
	migrate -database $(DB_PG_URL) -path $(MIG_PATH) up $(ARGS)

# Use it like
# Fully migrate:
# make mig-down
# Migrate from a vertion:
# make mig-down ARRGS=12
mig-down:
	migrate -database $(DB_PG_URL) -path $(MIG_PATH) down $(ARGS)

# Make a migration:
# make mk-mig ARRGS=[verb]-[entity]-[column/table/index etc.]
# make mk-mig ARRGS=add-vote-table
#
# Versioning options:
# Use -seq option to generate sequential up/down migration with N digits.
# Use -format option to specify a Go time format string. Note: migration with the same time cause "duplicate migration version" error.
# Use -tz option to specify the timezone that will be used when generating non-sequential migration (defaults: UTC).
# @migrate create -ext sql -dir $(MIG_PATH) -tz 'Asia/Tehran' $(ARGS)
mk-mig:
	@migrate create -ext sql -dir $(MIG_PATH) -seq $(ARGS)

# Force to retry a version which failed to run and is dirty
# Migrate force from a version:
# make mig-up ARGS=12
mig-force:
	@migrate -database $(DB_PG_URL) -path $(MIG_PATH) force $(ARGS)
