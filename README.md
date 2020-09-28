# snake-game-api

back-end project for snake.com, phaser game web app.

## Build Setup

``` bash
# install dependencies
go install

# serve at localhost:3060
go run app.go
```

## Database Setup

``` bash
# Start the cockroach node
cockroach start --insecure

# connect to CockroachDB SQL Client
cockroach sql --insecure
```

``` SQL
# create the Database and set it for current session
create database snake_game_db;
set database = snake_game_db;

# Run this script in the session with the CockroachDB SQL Client to create the records table
create table "tbl_record" (
    "id" SERIAL,
    "username" STRING(20),
    "best_score" INTEGER,
    "created_at" TIMESTAMPTZ,
    "updated_at" TIMESTAMPTZ,
    PRIMARY KEY ("id")
);
```
