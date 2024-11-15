package utils

import (
	"time"

	"gorm.io/gorm"
)

// GORM will by default sees ID variable name as primary key
/* gorm.Model provides a predefined struct:
type Model struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}*/
// GORM convention: first letter of struct variable needed to be uppercase

type Task struct {
	//ID          int       `json:"id"`
	gorm.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int
	User        User
}

// Now GORM knows UserID is foreign key by struct and stuctID!
// can use customized foreign key by tagging `gorm:"foreignKey:<attribute_name>"` at User struct line

type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required,fullname"`
  //Age      int    `json:"age" validate:"required,numeric,min=1"`
}

/* example
type Book struct {
  gorm.Model
  Name        string `json:"name"`
  Author      string `json:"author"`
  Description string `json:"description"`
  PublisherID uint
  Publisher   Publisher
  Authors     []Author `gorm:"many2many:author_books;"`
}

type Publisher struct {
  gorm.Model
  Details string
  Name    string
}

type Author struct {
  gorm.Model
  Name  string
  Books []Book `gorm:"many2many:author_books;"`
}

type AuthorBook struct {
  AuthorID uint
  Author   Author
  BookID   uint
  Book     Book
}

func createPublisher(db *gorm.DB, publisher *Publisher) error {
  result := db.Create(publisher)
  if result.Error != nil {
    return result.Error
  }
  return nil
}

func createAuthor(db *gorm.DB, author *Author) error {
  result := db.Create(author)
  if result.Error != nil {
    return result.Error
  }
  return nil
}

func createBookWithAuthor(db *gorm.DB, book *Book, authorIDs []uint) error {
  // First, create the book
  if err := db.Create(book).Error; err != nil {
    return err
  }

  return nil
}

func getBookWithPublisher(db *gorm.DB, bookID uint) (*Book, error) {
  var book Book
  result := db.Preload("Publisher").First(&book, bookID)
  if result.Error != nil {
    return nil, result.Error
  }
  return &book, nil
}

func getBookWithAuthors(db *gorm.DB, bookID uint) (*Book, error) {
  var book Book
  result := db.Preload("Authors").First(&book, bookID)
  if result.Error != nil {
    return nil, result.Error
  }
  return &book, nil
}

func listBooksOfAuthor(db *gorm.DB, authorID uint) ([]Book, error) {
  var books []Book
  result := db.Joins("JOIN author_books on author_books.book_id = books.id").
    Where("author_books.author_id = ?", authorID).
    Find(&books)
  if result.Error != nil {
    return nil, result.Error
  }
  return books, nil
}
*/
