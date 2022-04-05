package schema

import "forum/internal/constant"

func Query() map[string]string {
	return map[string]string{
		constant.QueAttach: `ATTACH DATABASE ? AS ?;`,

		constant.QueDetach: `DETACH DATABASE ?;`,

		constant.QueRestore: `INSERT INTO %s SELECT * FROM %s.%s;`,

		constant.QueInsert2: `INSERT INTO %s (%s, %s)  
		VALUES (?,?);`,

		constant.QueInsert3: `INSERT INTO %s (%s, %s, %s)
		VALUES (?,?,?);`,

		constant.QueInsert4: `INSERT INTO %s (%s, %s, %s, %s)
		VALUES (?,?,?,?);`,

		constant.QueInsert5: `INSERT INTO %s (%s, %s, %s, %s, %s)
		VALUES (?,?,?,?,?);`,

		constant.TabUsers: `CREATE TABLE IF NOT EXISTS %s (
		"id"		INTEGER NOT NULL,
		"login"		TEXT NOT NULL UNIQUE,
		"password"	TEXT NOT NULL,
		"email"		TEXT NOT NULL UNIQUE,
		"root"		INTEGER NOT NULL DEFAULT 0,
		"created"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY("id" AUTOINCREMENT));`,

		constant.TabCategories: `CREATE TABLE IF NOT EXISTS %s (
		"id"	INTEGER NOT NULL,
		"body"	TEXT NOT NULL UNIQUE,
		PRIMARY KEY("id"));`,

		constant.TabLikes: `CREATE TABLE IF NOT EXISTS %s (
		"id"	INTEGER NOT NULL,
		"body"	TEXT NOT NULL UNIQUE,
		PRIMARY KEY("id"));`,

		constant.TabPosts: `CREATE TABLE IF NOT EXISTS %s (
		"id"		INTEGER NOT NULL,
		"user"		INTEGER NOT NULL,
		"title"		TEXT NOT NULL,
		"body"		TEXT NOT NULL,
		"created"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY("user") REFERENCES "src_users"("id"),
		PRIMARY KEY("id" AUTOINCREMENT));`,

		constant.TabPostsLikes: `CREATE TABLE IF NOT EXISTS %s (
		"id"		INTEGER NOT NULL,
		"user"		INTEGER NOT NULL,
		"post"		INTEGER NOT NULL,
		"like"		INTEGER NOT NULL,
		"created"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY("user") REFERENCES "src_users"("id"),
		FOREIGN KEY("post") REFERENCES "posts"("id"),
		FOREIGN KEY("like") REFERENCES "src_likes"("id"),
		PRIMARY KEY("id" AUTOINCREMENT));`,

		constant.TabPostsCategories: `CREATE TABLE IF NOT EXISTS %s (
		"id"		INTEGER NOT NULL UNIQUE,
		"post"		INTEGER NOT NULL,
		"category"	INTEGER NOT NULL,
		FOREIGN KEY("post") REFERENCES "posts"("id"),
		FOREIGN KEY("category") REFERENCES "src_categories"("id"),
		PRIMARY KEY("id" AUTOINCREMENT));`,

		constant.TabComments: `CREATE TABLE  IF NOT EXISTS %s (
		"id"		INTEGER NOT NULL,
		"user"		INTEGER NOT NULL,
		"post"		INTEGER NOT NULL,
		"body"		TEXT NOT NULL,
		"created"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY("user") REFERENCES "src_users"("id"),
		FOREIGN KEY("post") REFERENCES "posts"("id"),
		PRIMARY KEY("id" AUTOINCREMENT));`,

		constant.TabCommentsLikes: `CREATE TABLE IF NOT EXISTS %s (
		"id"		INTEGER NOT NULL,
		"user"		INTEGER NOT NULL,
		"comment"	INTEGER NOT NULL,
		"like"		INTEGER NOT NULL,
		"created"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY("user") REFERENCES "src_users"("id"),
		FOREIGN KEY("comment") REFERENCES "comments"("id"),
		FOREIGN KEY("like") REFERENCES "src_likes"("id"),
		PRIMARY KEY("id" AUTOINCREMENT));`,

		constant.TabSessions: `CREATE TABLE IF NOT EXISTS %s (
		"user"		TEXT NOT NULL,
		"uuid"		TEXT NOT NULL UNIQUE,
		"expire"	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP);`,

		constant.QueSelect: `SELECT id
		FROM		%s
		WHERE		%s.%s = ?;`,

		constant.QueSelectPosts: `SELECT posts.id, posts.title, posts.body, posts.user, src_users.login, posts.created
		FROM		posts
		LEFT JOIN	src_users ON posts.user = src_users.id
		ORDER BY	posts.id DESC;`,

		constant.QueSelectUsers: `SELECT src_users.id, src_users.login, src_users.password, src_users.email, src_users.root, src_users.created
		FROM	src_users;`,

		constant.QueSelectCategories: `SELECT src_categories.id, src_categories.body
		FROM		src_categories
		ORDER BY	src_categories.id ASC;`,

		constant.QueSelectUserBy: `SELECT src_users.id, src_users.login, src_users.password, src_users.email, src_users.root, src_users.created
		FROM	src_users
		WHERE	src_users.%s = ?;`,

		constant.QueSelectPostsBy: `SELECT posts.id, posts.title, posts.body, posts.user, src_users.login, posts.created
		FROM		posts
		LEFT JOIN	src_users ON src_users.id = posts.user
		WHERE		%s.%s = ?
		ORDER BY	posts.id DESC;`,

		constant.QueSelectCommentsBy: `SELECT comments.id, src_users.login, comments.body, comments.created, comments.post
		FROM		comments
		LEFT JOIN	src_users ON comments.user = src_users.id
		WHERE		%s.%s = ?
		ORDER BY	comments.id DESC;`,

		constant.QueSelectSessionBy: `SELECT sessions.user, sessions.uuid, sessions.expire
		FROM	sessions
		WHERE	sessions.%s = ?
		AND 	sessions.%s = ?;`,

		constant.QueSelectCategoryBy: `SELECT src_categories.id, src_categories.body
		FROM		posts_categories
		LEFT JOIN	src_categories ON src_categories.id = posts_categories.category
		WHERE 		posts_categories.%s = ?
		ORDER BY	src_categories.body;`,

		constant.QueSelectLikeBy: `SELECT %s.id, %s.like, src_likes.body, %s.created
		FROM		%s
		LEFT JOIN	src_likes ON src_likes.id = %s.like
		WHERE		%s.%s = ?
		AND %s.%s = ?;`,

		constant.QueSelectPostsRatedBy: `SELECT posts.id, posts.title, posts.body, posts.user, src_users.login, posts.created, src_likes.body
		FROM		posts
		LEFT JOIN	src_users ON src_users.id = posts.user
		INNER JOIN	posts_likes ON posts_likes.post = posts.id
		LEFT JOIN	src_likes ON src_likes.id = posts_likes.like
		WHERE		%s.%s = ?
		ORDER BY	posts.id DESC;`,

		constant.QueSelectCommentsRatedBy: `SELECT comments.id, src_users.login, comments.body, comments.created, comments.post, src_likes.body
		FROM		comments
		LEFT JOIN	src_users ON src_users.id = comments.user
		INNER JOIN	comments_likes ON comments_likes.comment = comments.id
		LEFT JOIN	src_likes ON src_likes.id = comments_likes.like
		WHERE		%s.%s = ?
		ORDER BY	comments.id DESC`,

		constant.QueSelectPostsAndCategoryBy: `SELECT posts.id, posts.title, posts.body, posts.user, src_users.login, posts.created
		FROM		posts
		LEFT JOIN	src_users ON src_users.id = posts.user
		LEFT JOIN	posts_categories ON posts.id = posts_categories.post
		WHERE		%s.%s = ?
		ORDER BY	posts.id DESC;`,

		constant.QueSelectLikeCountBy: `SELECT src_likes.id, src_likes.body, COUNT(*) as N
		FROM		%s
		LEFT JOIN	src_likes ON src_likes.id = %s.like
		WHERE 		%s.%s = ?
		GROUP BY	src_likes.body;`,

		constant.QueSelectLikedOrNot: `SELECT src_likes.body
		FROM		%s
		LEFT JOIN	src_likes ON src_likes.id = %s.like
		WHERE 		%s.%s = ?
		AND 		%s.%s = ?;`,

		constant.QueDeleteBy: `DELETE
		FROM	%s
		WHERE	%s.%s = ?;`,

		constant.QueSelectCount: `SELECT COUNT(*) as N
		FROM	%s;`,
	}
}
