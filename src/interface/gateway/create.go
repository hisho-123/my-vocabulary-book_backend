package gateway

import (
	"backend/src/domain"
	"backend/src/infra/db"
	"fmt"
	"log"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

// 単語帳の作成
func CreateBookByUserId(book domain.CreateBookInput) error {
	const MAX_BOOK_NAME int = 20
	if len(book.BookName) > MAX_BOOK_NAME {
		log.Printf("error: Book name too long.")
		return fmt.Errorf(domain.UnprocessableEntity)
	}

	db := db.OpenDB()
	defer db.Close()

	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		log.Println("error: ", err)
		return fmt.Errorf(domain.InternalServerError)
	}
	defer tx.Rollback() // エラー時は自動ロールバック

	queryCreateBook := "insert into books (user_id, book_name) values (?, ?);"
	result, err := tx.Exec(queryCreateBook, strconv.Itoa(book.UserId), book.BookName)
	if err != nil {
		log.Println("error: ", err)
		return fmt.Errorf(domain.InternalServerError)
	}

	// LastInsertId()で今挿入したbook_idを取得
	bookId, err := result.LastInsertId()
	if err != nil {
		log.Println("error: ", err)
		return fmt.Errorf(domain.InternalServerError)
	}

	// wordの挿入
	queryCreateWord := "insert into words (book_id, word, translated_word) values (?, ?, ?)"
	for _, v := range book.Words {
		_, err := tx.Exec(queryCreateWord, bookId, v.Word, v.Translated)
		if err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok {
				if mysqlErr.Number == 1406 {
					log.Println("error: ", err)
					return fmt.Errorf(domain.UnprocessableEntity)
				}
			}
			log.Println("error: ", err)
			return fmt.Errorf(domain.InternalServerError)
		}
	}

	// コミット
	if err := tx.Commit(); err != nil {
		log.Println("error: ", err)
		return fmt.Errorf(domain.InternalServerError)
	}

	return nil
}

// 既存の単語帳に単語を登録
func CreateWordByBookId(bookId int, word string, translated string) error {
	db := db.OpenDB()
	defer db.Close()

	queryCreateWord := "insert into words (book_id, word, translated_word) values (?, ?, ?)"
	_, err := db.Exec(queryCreateWord, bookId, word, translated)
	if err != nil {
		log.Println("error: ", err)
		return fmt.Errorf(domain.InternalServerError)
	}
	return nil
}
