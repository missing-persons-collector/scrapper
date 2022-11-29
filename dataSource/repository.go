package dataSource

func SaveCountry(country *Country) error {
	result := db.Create(country)

	return result.Error
}

func FindCountry(name string, country *Country) error {
	result := db.Where("name = ?", name).First(country)

	return result.Error
}

func SavePerson(person *Person) error {
	result := db.Create(person)

	return result.Error
}
