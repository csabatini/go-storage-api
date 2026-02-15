# Architectural Decisions

## Decision Log

### ADR-001: Interface-Based Storage Abstraction

- **Date:** 2026-02-15
- **Status:** Accepted
- **Context:** The API must support multiple file protocols (local filesystem, SMB, FTP) and allow swapping backends without changing HTTP handler code. We need a clean abstraction that decouples protocol-specific logic from the API layer.
- **Decision:** Define a single `storage.Storage` Go interface with five methods (`List`, `Read`, `Write`, `Delete`, `Stat`). Each backend implements this interface in its own package. HTTP handlers accept the interface via dependency injection.
- **Consequences:**
  - Adding a new backend (e.g. S3, SFTP) requires only implementing the interface in a new package and adding a case to the startup switch — no handler changes.
  - Testing becomes trivial: mock the interface to test handlers without a real filesystem.
  - Each backend is isolated; SMB dependencies don't affect FTP code.
  - Tradeoff: protocol-specific features (e.g. SMB file locking) cannot be exposed through the generic interface without extending it.

### ADR-002: Streaming I/O via io.Reader / io.ReadCloser

- **Date:** 2026-02-15
- **Status:** Accepted
- **Context:** The service will handle files of arbitrary size. Loading entire files into memory (e.g. `[]byte`) would cause out-of-memory conditions for large files and increase latency.
- **Decision:** The `Storage.Read` method returns `io.ReadCloser` and `Storage.Write` accepts `io.Reader`. File content is streamed from source to destination without full buffering.
- **Consequences:**
  - Memory usage stays constant regardless of file size.
  - Large file transfers (multi-GB) are supported without special handling.
  - Callers must remember to close the `ReadCloser` to avoid resource leaks.
  - Error handling during streaming is more nuanced — partial writes are possible if the stream fails mid-transfer.

### ADR-003: Backend-Per-Package Structure

- **Date:** 2026-02-15
- **Status:** Accepted
- **Context:** Each file protocol (local, SMB, FTP) has different dependencies, connection semantics, and configuration requirements. Mixing them in a single package would create tight coupling and import bloat.
- **Decision:** Each backend lives in its own sub-package under `internal/storage/` (e.g. `internal/storage/local/`, `internal/storage/smb/`, `internal/storage/ftp/`). Each package only imports the libraries it needs.
- **Consequences:**
  - Clear separation of concerns — changes to the FTP backend cannot break the SMB backend.
  - Build dependencies are scoped: if you only use the local backend, SMB/FTP libraries are not compiled in (assuming build tags or selective imports).
  - More packages to navigate, but each is small and focused.

### ADR-004: Configuration via Environment Variables

- **Date:** 2026-02-15
- **Status:** Accepted
- **Context:** The service needs different configuration per environment (development, staging, production) and per backend (local root path vs. SMB host/share vs. FTP credentials). We need a configuration approach that works across container orchestrators, CI/CD, and local development.
- **Decision:** All configuration is loaded from environment variables via `internal/config/`. A `.env` file is supported for local development (never committed). The `STORAGE_BACKEND` variable selects the active backend; backend-specific variables (e.g. `SMB_HOST`, `FTP_PORT`) configure that backend.
- **Consequences:**
  - Follows 12-factor app methodology. Works naturally with Docker, Kubernetes, and CI/CD.
  - No config files to manage or keep in sync across environments.
  - Credentials are never hardcoded or committed to version control.
  - Tradeoff: complex nested configuration is harder to express in flat env vars compared to YAML/TOML.

### ADR-005: Path as Query Parameter

- **Date:** 2026-02-15
- **Status:** Accepted
- **Context:** File paths can contain special characters, deeply nested directories, and characters that conflict with URL path segments (e.g. `/`, `.`, `%`). Encoding file paths as part of the URL path creates ambiguity and routing issues.
- **Decision:** File paths are passed as a `path` query parameter (e.g. `GET /api/v1/files?path=/docs/report.pdf`) rather than embedded in the URL path.
- **Consequences:**
  - No ambiguity between route segments and file path segments.
  - Paths with special characters are handled naturally by standard query parameter encoding.
  - All file endpoints share a consistent parameter convention.
  - Tradeoff: slightly less "RESTful" than path-based resource identification, but more practical for arbitrary filesystem paths.

### ADR-006: Path Traversal Prevention via Middleware

- **Date:** 2026-02-15
- **Status:** Accepted
- **Context:** File path manipulation is the primary attack vector for a file service. Path traversal attacks (e.g. `../../etc/passwd`) could allow access to files outside the intended scope.
- **Decision:** A `pathguard` middleware normalizes all incoming file paths and rejects any path containing `..` or absolute path escapes before the request reaches a handler. Each backend additionally scopes operations to its configured root directory or share.
- **Consequences:**
  - Defense in depth: two layers of protection (middleware + backend scoping).
  - Centralized validation — no need to repeat path checks in every handler.
  - Overly strict normalization could reject legitimate paths in edge cases, but this is a safer default.

### ADR-007: Use internal/ Package Convention

- **Date:** 2026-02-15
- **Status:** Accepted
- **Context:** Go's `internal/` directory convention prevents external packages from importing internal code. Since this is a standalone service (not a library), all application code should be private.
- **Decision:** All application packages live under `internal/`. Only `cmd/server/main.go` sits outside as the entry point.
- **Consequences:**
  - External consumers cannot import our handlers, storage implementations, or config — reducing the API surface we need to maintain.
  - Follows standard Go project layout conventions.
  - If we later need to expose a client SDK, we would create a separate `pkg/` directory for public types.
