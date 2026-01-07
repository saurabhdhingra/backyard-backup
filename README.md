# Backyard Backup CLI

A versatile command-line utility for backing up databases (PostgreSQL, MySQL, MongoDB, SQLite) to various storage backends (Local, AWS S3) with support for compression, scheduling, and notifications.

## Features

-   **Databases**: PostgreSQL, MySQL, MongoDB, SQLite.
-   **Storage**: Local Filesystem, AWS S3 (with static credentials).
-   **Compression**: Automatic Gzip compression.
-   **Scheduling**: Cron-based scheduling for recurring backups.
-   **Notifications**: Slack webhook integration.
-   **Config**: Simple YAML-based configuration.

## Prerequisites

Ensure you have the command-line tools installed for your database(s):

-   **PostgreSQL**: `pg_dump`, `psql` (Install: `brew install libpq`)
-   **MySQL**: `mysqldump`, `mysql`
-   **MongoDB**: `mongodump`, `mongorestore` (Install: `brew install mongodb-database-tools`)
-   **SQLite**: `sqlite3`

## Installation

Clone the repository and build:

```bash
git clone https://github.com/saurabhdhingra/backyard-backup.git
cd backyard-backup
go mod tidy
go build -o dbbackup
```

## Configuration

Copy the example configuration:

```bash
cp config.yaml.example config.yaml
```

Edit `config.yaml` with your credentials. You can use standard fields (`host`, `user`, etc.) or a raw connection string (`dsn`).

### Example (PostgreSQL + S3)

```yaml
database:
  type: postgres
  dsn: "postgresql://user:pass@host:5432/db?sslmode=require"

storage:
  type: s3
  bucket: "my-backups"
  region: "us-east-1"
  access_key: "AWS_ACCESS_KEY"
  secret_key: "AWS_SECRET_KEY"

backup:
  compression: true
  schedule: "@daily"

notify:
  enabled: true
  slack_webhook: "https://hooks.slack.com/..."
```

## Usage

### Backup
Run an immediate backup:
```bash
./dbbackup backup
```

### Restore
Restore from a specific backup file in your storage:
```bash
./dbbackup restore --file db_20240101.sql.gz
```
*Note: If using local storage, provide the filename relative to the backup directory configured.*

### Schedule
Start the scheduler process:
```bash
./dbbackup schedule
```

## Acknowledgement
https://roadmap.sh/projects/database-backup-utility

## License
MIT
