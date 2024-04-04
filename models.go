package main

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title    string
	Length   int
	Language string
	Authors  []Author `gorm:"many2many:author_books;"`
}

type Author struct {
	gorm.Model
	Name  string
	Books []Book `gorm:"many2many:author_books;"`
}
