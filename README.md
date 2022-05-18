## forum: version authentication

forum based on: "basic forum" + "authentication"

In this version:

- authentication third-part by `Google`,`Github` or `Facebook`.
- `clientID` and `clientSecret` must be written in `/.env` file.

### Objectives

This project consists in creating a web forum that allows:

- communication between users.
- associating categories to posts.
- liking and disliking posts and comments.
- filtering posts.

### SQLite

In this project we use sqlite db.

### Run

To run project please type in command line `go run ./cmd/forumsqlite/` or `make run`

### Randomizer
Delete current db `forum/db/database-sqlite3.db`

To use random `user`, `categories` and `post` for testing, need to uncomment this lines in directory : `forum/internal/app/app.go`:

- `_, _, schema := repo.ExportSettings()`
- `repository.NewLoremIpsum().Run(db, schema)`

You can then return comment.

For example random user to login:

login: `blackbeard`

password: `12345Aa`

### Docker

Run command `make docker` in command line.
