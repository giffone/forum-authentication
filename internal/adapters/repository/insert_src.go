package repository

import (
	"database/sql"
	"fmt"
	"forum/internal/constant"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func InsertSrc(db *sql.DB, q map[string]string) {
	src := make(map[string][]string)
	src[constant.TabLikes] = []string{"like", "dislike"}
	src[constant.TabCategories] = []string{"animal", "city", "body part", "vehicle", "actor", "18+",
		"clothing", "sport", "drink", "sea creature", "chemical element", "fruit", "country"}

	que := fmt.Sprintf(q[constant.QueInsert5], constant.TabUsers,
		constant.FieldLogin, constant.FieldPassword, constant.FieldEmail, constant.FieldRoot, constant.FieldCreated)
	pass, err := bcrypt.GenerateFromPassword([]byte("12345Aa"), bcrypt.MinCost)
	if err != nil {
		log.Printf("insert source: password for admin: %v\n", err)
	}
	_, err = db.Exec(que, "admin", string(pass), "admin@mail.ru", 1, time.Now()) // root=1 for admin
	if err != nil {
		log.Printf("insert source: admin did not created: %v\n", err)
	}

	for table, source := range src {
		numberLines := 0
		que = fmt.Sprintf(q[constant.QueInsert2], table, constant.FieldID, constant.FieldBody)
		for id, value := range source {
			_, err := db.Exec(que, id+1, value)
			if err != nil {
				log.Printf("insert source: \"%s\" did not inserted to \"%s\"\n", value, table)
				continue
			}
			numberLines++
		}
		log.Printf("insert source: %d lines added to \"%s\" table\n", numberLines, table)
	}
}
