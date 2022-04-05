## forum

### Objectives

This project consists in creating a web forum that allows:

- communication between users.
- associating categories to posts.
- liking and disliking posts and comments.
- filtering posts.

#### SQLite

In this project we use sqlite db.

#### Run

To run project please type in command line `go run ./cmd/forumsqlite/` or `make run`

#### Randomizer
Delete current db `forum/db/database-sqlite3.db`

To use random `user`, `categories` and `post` for testing, need to uncomment this lines in directory : `forum/internal/app/app.go`:

- `_, _, schema := repo.ExportSettings()`
- `repository.NewLoremIpsum().Run(db, schema)`

You can then return comment.

For example random user to login:

login: `blackbeard`

password: `12345Aa`

#### Docker

Run command `make docker` in command line.

#### Thanks

Thanks for support @kazykenov1, @Jambul, @damirkap89, @alseiitov, @devstackq