# go-tinker

Go learning project - experiments with Go fundamentals, web servers, and GORM ORM.

## Structure

```
.
├── cmd/
│   ├── basics/                    # Basics demo entry point
│   ├── lessons/                   # ORM lessons (kebab-case naming)
│   │   ├── orm-app/
│   │   ├── orm-cond-queries/
│   │   ├── orm-relations/
│   │   ├── orm-relations-download/
│   │   ├── orm-hooks/
│   │   ├── orm-transactions/
│   │   ├── orm-raw-exec/
│   │   ├── orm-errors-validation/
│   │   ├── orm-test/
│   │   └── orm-optimization/
│   └── movies-crud-cli/           # Movies CRUD CLI app
├── internal/
│   ├── greeting/                  # Greeting utility package
│   ├── even/                      # Even number filter utility
│   ├── models/                    # GORM models (domain entities)
│   ├── store/                     # In-memory data stores
│   │   ├── entities/              # Business entities with handlers
│   │   └── schemas/               # Request/Response DTOs
│   └── web/
│       ├── fiber/                 # Fiber web server
│       ├── stdlib/                # net/http web server
│       └── templates/             # HTML templates
├── basics/                        # Go fundamentals exercises
├── .env.example                   # Environment template
├── Makefile                       # Build/run commands
└── go.mod
```

## Quick Start

```bash
# Build all binaries
make build-all

# Run basics demo
make run-basics

# Run a specific lesson
make run-lesson-orm-app

# List all available lessons
make list-lessons

# Clean build artifacts
make clean
```

## Configuration

Copy `.env.example` to `.env` and update the database URL:

```bash
cp .env.example .env
```

Each lesson in `cmd/lessons/` has its own `.env` and `.env_default` files.

## Requirements

- Go 1.26+
- PostgreSQL (for ORM lessons)
