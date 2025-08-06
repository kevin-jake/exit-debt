# PostgreSQL Docker Setup

This directory contains Docker Compose configuration for running PostgreSQL locally for the Go Debt Tracker application.

## Quick Start

### Start PostgreSQL

```bash
docker-compose up -d
```

### Stop PostgreSQL

```bash
docker-compose down
```

### Stop and remove all data

```bash
docker-compose down -v
```

## Configuration

### Database Details

- **Host**: localhost
- **Port**: 5432
- **Database**: debt_tracker
- **Username**: postgres
- **Password**: postgres

### Connection String

```
postgresql://postgres:postgres@localhost:5432/debt_tracker
```

## Features

- **Persistent Data**: Data is stored in a Docker volume (`postgres_data`)
- **Initialization Script**: The `init.sql` file runs automatically on first startup
- **Network Isolation**: Uses a dedicated Docker network
- **Auto-restart**: Container restarts automatically unless manually stopped

## Customization

### Environment Variables

You can modify the database configuration by editing the `environment` section in `docker-compose.yml`:

```yaml
environment:
  POSTGRES_DB: your_database_name
  POSTGRES_USER: your_username
  POSTGRES_PASSWORD: your_password
```

### Initialization Script

Edit `init.sql` to add:

- Database extensions
- Initial schemas
- Sample data
- Custom functions

## Troubleshooting

### Check container status

```bash
docker-compose ps
```

### View logs

```bash
docker-compose logs postgres
```

### Connect to database

```bash
docker-compose exec postgres psql -U postgres -d debt_tracker
```

### Reset database

```bash
docker-compose down -v
docker-compose up -d
```
