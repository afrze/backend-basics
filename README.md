# Backend Basics

This repo contains Go, Docker, Kubernetes, gRPC, AWS setup using a simple use case. This repository is tinkering experiment based on the Udemy course [Backend Master Class](https://www.udemy.com/course/backend-master-class-golang-postgresql-kubernetes/?kw=backend&src=sac&couponCode=OF83024F)

## Database design

![db schema diagram](https://i.ibb.co/j8R9Rt5/db-schema-image.png)

Reference Course: [Backend Master Class](https://www.udemy.com/course/backend-master-class-golang-postgresql-kubernetes/?couponCode=24T4MT92724B)

## Docker

Commands

- `docker ps` to list all the running containers
- `docker images` to list all the existing images
- To pull an image from the `docker hub`

```bash
docker pull <image-name>:<tag-name>
```

- To see vulnerabilities and recommendations of a pulled image

```bash
docker scout quickview <image-name>:<tag-name>
```

- To remove images `docker rmi <image_id>`
- To remove containers `docker rm <container_id>`
- To remove containers and volumes `docker rm --volumes <volume_name>`
- Run commands in the Container.

```bash
docker exec -it <container_name_or_id> <command> [args]
```  

- View container logs `docker logs <container_name_or_id>`
- To stop a container `docker stop <container_name>`

### Postgres Docker commands

- To run a container from a pulled image

```bash
docker run --name <container_name> -e <environment_variable> -d <image_name>
# Example command
docker run --name postgres17 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine
```

- Same command with port mapping

```bash
docker run --name <container_name> -e <environment_variable> -p <host_ports:container_ports> -d <image_name>
# Example command
docker run --name postgres17 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:alpine
```

- To create a database inside a docker

```bash
docker exec -it <container_name> createdb --username=<user_name> --owner=<owner_name> <databse_name>
# For Example
docker exec -it postgres createdb --username=root --owner=root simple_bank
```

- To access the database directly using the `psql` command

```bash
docker exec -it <container_name> psql -U <user_name> <database_name>
        
docker exec -it postgres psql -U root simple_bank
```

## Migrate - Database Migrations

When using Go, we have a handy migrations scripts generator to handle database migrations. Please refer: [golang-migrate/migrate](https://github.com/golang-migrate/migrate).

These create `up` and `down` scripts for our database.

Commands

- To create a migration script

```bash
migrate create -ext <file_extension> -dir <target_directory> -seq <sequence_name>
# Example
migrate create -ext sql -dir db/migration -seq init_schema
```

- To migrate: Up and Down sequence.

```bash
migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
    
migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
```
## SQLC

Alternate way of handling Database operations

- To initialize `sqlc init` - This generates a `sqlc.yaml` file
- Settings for `sqlc.yaml` 
    
```yaml
  version: "2"
  sql:
    - engine: "postgresql"
      queries: "./db/query"
      schema: "./db/migration"
      gen:
        go:
          package: "db"
          out: "./db/sqlc"
          sql_package: "pgx/v5"
```
    
- Refer this documentation: [SQLC Documentation](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#setting-up)
- Example query for `sqlc` 
    
```sql
  -- name: CreateAccount :one
  INSERT INTO accounts (
    owner,
    balance,
    currency
  ) VALUES (
    $1, $2, $3
  ) RETURNING *;
    
  -- name: GetAccount :one
  SELECT * FROM accounts 
  WHERE id = $1 LIMIT 1;
    
  -- name: ListAccounts :many
  SELECT * FROM accounts
  ORDER BY id 
  LIMIT $1
  OFFSET $2;
    
  -- name: UpdateAccount :one
  UPDATE accounts
  SET balance = $2
  WHERE id = $1
  RETURNING *;
    
  -- name: DeleteAccount :exec
  DELETE FROM accounts
  WHERE id = $1;
```
- To generate files `sqlc generate`
