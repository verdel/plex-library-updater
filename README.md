# Plex Library Refresher

A lightweight Go service that periodically refreshes all libraries on a Plex server.

The tool connects to the Plex API and triggers a refresh for every library on a configurable interval. It is designed to run as a small background service or container.

## Features

- Periodic refresh of all Plex libraries
- Parallel library scanning
- Retry with exponential backoff
- Graceful shutdown
- Structured JSON logging using zerolog
- Simple configuration via environment variables

## Configuration

The application is configured using environment variables.

| Variable          | Description                                                | Required |
| ----------------- | ---------------------------------------------------------- | -------- |
| `PLEX_URL`        | Plex server URL (e.g. `http://localhost:32400`)            | yes      |
| `PLEX_TOKEN`      | Plex API token                                             | yes      |
| `UPDATE_INTERVAL` | Refresh interval (default "1m")                            | no       |
| `MAX_RETRIES`     | Number of retry attempts for Plex API calls (default: `3`) | no       |
| `LOG_FORMAT`      | Log format (`pretty`/`json`) (default: `pretty`)           | no       |

## Running locally

```bash
export PLEX_URL=http://localhost:32400
export PLEX_TOKEN=your_token
export UPDATE_INTERVAL=30m
export MAX_RETRIES=3
export LOG_FORMAT=json

go run .
```

## Example log output

### Pretty console format (default)

Example:

```text
2026-03-13T16:12:01Z INF plex refresher started interval=30m0s retries=3
2026-03-13T16:12:01Z INF starting library refresh
2026-03-13T16:12:02Z INF refresh triggered library=Movies library_id=2 status=200
2026-03-13T16:12:02Z INF refresh triggered library=TV Shows library_id=3 status=200
```

### JSON format

Example:

```json
{
  "level": "info",
  "interval": "30m0s",
  "message": "plex refresher started"
}
```

```json
{
  "level": "info",
  "library": "Movies",
  "library_id": "2",
  "status": 200,
  "message": "refresh triggered"
}
```

## How it works

1. The service periodically requests the list of libraries from the Plex API.
2. Each library refresh is triggered in parallel.
3. Failed API calls are retried with exponential backoff.
4. The service shuts down gracefully when receiving `SIGINT` or `SIGTERM`.
