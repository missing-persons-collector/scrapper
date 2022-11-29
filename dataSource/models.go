package dataSource

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Country struct {
	gorm.Model

	UUID string
	Name string `gorm:"unique"`

	Persons []Person `gorm:"foreignKey:CountryID"`
}

type Person struct {
	gorm.Model

	UUID      string
	CountryID uint

	Entries []Entry `gorm:"foreignKey:PersonID"`
}

type Entry struct {
	gorm.Model

	PersonID uint
	UUID     string

	Key   string
	Value string
}

func (Country) TableName() string {
	return "countries"
}

func (Person) TableName() string {
	return "people"
}

func (Entry) TableName() string {
	return "entries"
}

func NewPerson() Person {
	return Person{
		UUID:    uuid.New().String(),
		Entries: make([]Entry, 0),
	}
}

func NewCountry(name string) Country {
	return Country{
		UUID:    uuid.New().String(),
		Name:    name,
		Persons: make([]Person, 0),
	}
}

func NewEntry(key string, value string) Entry {
	return Entry{
		UUID:  uuid.New().String(),
		Key:   key,
		Value: value,
	}
}
