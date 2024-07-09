# Running Tests for Heisen-Flow

## Prerequisites

- Docker and Golang installed on your system.

## Steps to Run Tests

### 1. Pull the PostgreSQL Docker Image

```sh
docker pull postgres:16.2-bookworm
```

### 2. Run the PostgreSQL Docker Container

```sh
docker run --name postgres_test -e POSTGRES_PASSWORD=123456 -p 5432:5432 -d postgres:16.2-bookworm
```

### 3. Set Up the Database

#### Connect to the Database

Connect to the PostgreSQL container using `psql`:

```sh
docker exec -it postgres_test psql -U postgres
```

#### Create the `root` User and `heisenflow_test` Database

Once connected to the PostgreSQL prompt, run the following SQL commands:

```sql
CREATE USER root WITH PASSWORD '123456';
CREATE DATABASE heisenflow_test;
GRANT ALL PRIVILEGES ON DATABASE heisenflow_test TO root;

\c heisenflow_test

GRANT ALL PRIVILEGES ON SCHEMA public TO root;
```
Ensure to execute these commands seperately

### 4. Run the Tests

Navigate to your Golang application's root directory and run the tests using the `go test` command:

```sh
go test ./test/... -v
```

This command will execute all the tests located in the `test` folder.

## Configuration Summary

Here is the configuration used for the database:

```yaml
db:
  user: "root"
  pass: "123456"
  host: "localhost"
  port: 5432
  db_name: "heisenflow_test"
```

## Notes

- Make sure Docker is running before executing the commands.
- Ensure that the port `5432` is not being used by another PostgreSQL instance on your machine.
- The `go test` command assumes that your tests are located inside the `test` folder.

