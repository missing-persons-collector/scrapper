package dataSource

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        string `gorm:"primary_key;type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Country struct {
	Base

	Name    string   `gorm:"unique"`
	Persons []Person `gorm:"foreignKey:CountryID"`
}

func (u *Country) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}

func (*Country) TableName() string {
	return "countries"
}

func NewCountry(name string) Country {
	return Country{
		Name:    name,
		Persons: make([]Person, 0),
	}
}

type Person struct {
	Base

	CustomID  string
	CountryID string

	Name             string `gorm:"default:null"`
	LastName         string `gorm:"default:null"`
	MaidenName       string `gorm:"default:null"`
	Gender           string `gorm:"default:null"`
	DOB              string `gorm:"default:null"`
	POB              string `gorm:"default:null"`
	Citizenship      string `gorm:"default:null"`
	PrimaryAddress   string `gorm:"default:null"`
	SecondaryAddress string `gorm:"default:null"`
	Country          string `gorm:"default:null"`
	ImageID          string `gorm:"default:null"`

	Height   string `gorm:"default:null"`
	Hair     string `gorm:"default:null"`
	EyeColor string `gorm:"default:null"`
	Weight   string `gorm:"default:null"`

	// Date of disappearance
	DOD string `gorm:"default:null"`
	// place of disappearance
	POD string `gorm:"default:null"`

	Description string `gorm:"default:null"`
}

func (*Person) TableName() string {
	return "people"
}

func (u *Person) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()

	return nil
}

func NewPerson(id string) Person {
	return Person{
		CustomID: id,
	}
}
