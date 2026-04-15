# Daedalus — Clean Architecture

## Overview

The Daedalus backend follows Clean Architecture (Hexagonal Architecture), inspired by the Kliops reference project.

## Layer Diagram

```
┌──────────────────────────────────────────────┐
│                 HTTP Request                  │
└────────────────────┬─────────────────────────┘
                     ▼
┌──────────────────────────────────────────────┐
│          adapters/http/ (Routes)              │
│  Flask Blueprint, Marshmallow Schemas,        │
│  Error Handlers, Middleware                   │
└────────────────────┬─────────────────────────┘
                     ▼
┌──────────────────────────────────────────────┐
│        core/services/ (Use Cases)             │
│  ProjectService: create, list, update,        │
│  autosave, archive, restore, delete           │
│  Business validation, logging                 │
└────────────────────┬─────────────────────────┘
                     ▼
┌──────────────────────────────────────────────┐
│         core/ports/ (Interfaces)              │
│  ProjectRepository ABC                        │
│  Defines contracts, NO implementation         │
└────────────────────┬─────────────────────────┘
                     ▼
┌──────────────────────────────────────────────┐
│    adapters/persistence/ (Repository)         │
│  SQLAlchemyProjectRepository                  │
│  ORM Model ↔ Domain Entity mapping            │
└──────────────────────────────────────────────┘
                     ▼
┌──────────────────────────────────────────────┐
│              PostgreSQL 16                     │
└──────────────────────────────────────────────┘
```

## Dependency Rule

Dependencies point INWARD:
- **adapters** depend on **core** (never the reverse)
- **core/services** depend on **core/ports** and **core/domain**
- **core/domain** has ZERO external dependencies

## Key Patterns (from Kliops)

1. **Port/Adapter** — `ProjectRepository` ABC defines the contract, `SQLAlchemyProjectRepository` implements it
2. **Domain Entities** — Pure Python dataclasses, no framework coupling
3. **Error Wrapping** — Domain exceptions mapped to HTTP codes in error handlers
4. **Dependency Injection** — Service receives repository via constructor in `app/__init__.py`
5. **Structured Logging** — JSON format (Winston-style)
