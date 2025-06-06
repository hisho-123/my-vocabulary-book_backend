package gateway

import (
	"backend/src/domain"
	"backend/src/infra/db"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func GetBookListByUserId(userId int) (books []domain.Book, err error) {
	db := db.OpenDB()
	defer db.Close()

	queryGetBookList := "select book_id, book_name, first_review from books where user_id = ?"

	var bookId int
	var bookName string
	var firstReview sql.NullString

	rows, err := db.Query(queryGetBookList, strconv.Itoa(userId))
	if err != nil {
		log.Println("error: ", err)
		return nil, fmt.Errorf(domain.InternalServerError)
	}

	for rows.Next() {
		if err := rows.Scan(&bookId, &bookName, &firstReview); err != nil {
			log.Println("error: ", err)
			return nil, fmt.Errorf(domain.InternalServerError)
		}

		books = append(books, domain.Book{
			Id:          bookId,
			UserId:      userId,
			Name:        bookName,
			FirstReview: firstReview.String,
		})
	}

	if err := rows.Err(); err != nil {
		log.Println("error: ", err)
		return nil, fmt.Errorf(domain.InternalServerError)
	}

	return books, nil
}

func GetBookByBookId(bookId int) (book *domain.GetBookOutput, err error) {
	db := db.OpenDB()
	defer db.Close()

	// 構造体の初期化
	book = &domain.GetBookOutput{}

	// 単語帳の名前を取得
	queryGetBook := "select user_id, book_name from books where book_id = ?"
	err = db.QueryRow(queryGetBook, bookId).Scan(&book.UserId, &book.BookName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("error: ", err)
			return nil, fmt.Errorf(domain.InternalServerError)
		}
		log.Println("error: ", err)
		return nil, fmt.Errorf(domain.InternalServerError)
	}

	// 単語取得
	queryRowWords := "select word_id, word, translated_word from words where book_id = ?"
	wordsRows, err := db.Query(queryRowWords, bookId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("error: ", err)
			return nil, fmt.Errorf(domain.InternalServerError)
		}
		log.Println("error: ", err)
		return nil, fmt.Errorf(domain.InternalServerError)
	}

	var wordId int
	var word string
	var translated string

	for wordsRows.Next() {
		if err := wordsRows.Scan(&wordId, &word, &translated); err != nil {
			log.Println("error: ", err)
			return nil, fmt.Errorf(domain.InternalServerError)
		}

		book.Words = append(book.Words, domain.Word{
			Id:         wordId,
			Word:       word,
			Translated: translated,
		})
	}

	if err := wordsRows.Err(); err != nil {
		log.Println("error: ", err)
		return nil, fmt.Errorf(domain.InternalServerError)
	}

	return book, nil
}
