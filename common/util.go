package common

import (
	"missingPersons/dataSource"
)

func BuildFieldMap() map[string]string {
	fieldMap := make(map[string]string)
	fieldMap["Ime"] = "Name"
	fieldMap["Prezime"] = "LastName"
	fieldMap["Djevojačko prezime"] = "MaidenName"
	fieldMap["Spol"] = "Gender"
	fieldMap["Datum rođenja"] = "DOB"
	fieldMap["Mjesto rođenja"] = "POB"
	fieldMap["Državljanstvo"] = "Citizenship"
	fieldMap["Prebivalište"] = "PrimaryAddress"
	fieldMap["Boravište"] = "SecondaryAddress"
	fieldMap["Država"] = "Country"
	fieldMap["Visina"] = "Height"
	fieldMap["Kosa"] = "Hair"
	fieldMap["Boja očiju"] = "EyeColor"
	fieldMap["Datum nestanka"] = "DOD"
	fieldMap["Mjesto nestanka"] = "POD"
	fieldMap["Okolnosti nestanka"] = "Description"

	return fieldMap
}

func BuildSerbiaFieldMap() map[string]string {
	fieldMap := make(map[string]string)
	fieldMap["Ime"] = "Name"
	fieldMap["Prezime"] = "LastName"
	fieldMap["Pol"] = "Gender"
	fieldMap["Datum rođenja"] = "DOB"
	fieldMap["Mjesto rođenja"] = "POB"
	fieldMap["Državljanstvo"] = "Citizenship"
	fieldMap["Prebivalište"] = "PrimaryAddress"
	fieldMap["Boravište"] = "SecondaryAddress"
	fieldMap["Težina"] = "Weight"
	fieldMap["Država"] = "Country"
	fieldMap["Visina"] = "Height"
	fieldMap["Boja kose"] = "Hair"
	fieldMap["Boja očiju"] = "EyeColor"
	fieldMap["Datum nestanka"] = "DOD"
	fieldMap["Mesto nestanka"] = "POD"
	fieldMap["Opis u trenutku nestanka"] = "Description"

	return fieldMap
}

func PersonFromRawPerson(id string, countryId string, raw RawPerson) dataSource.Person {
	person := dataSource.NewPerson(id)

	person.Name = raw.Name
	person.LastName = raw.LastName
	person.DOD = raw.DOD
	person.POB = raw.POB
	person.DOB = raw.DOB
	person.POD = raw.POD
	person.Country = raw.Country
	person.CountryID = countryId
	person.Citizenship = raw.Citizenship
	person.CustomID = id
	person.PrimaryAddress = raw.PrimaryAddress
	person.SecondaryAddress = raw.SecondaryAddress
	person.MaidenName = raw.Name
	person.Height = raw.Name
	person.EyeColor = raw.Name
	person.Gender = raw.Name
	person.Hair = raw.Name
	person.Weight = raw.Weight
	person.Description = raw.Name

	return person
}

func BatchRawPeople(items []RawPerson, batchNum int) [][]RawPerson {
	batch := make([][]RawPerson, 0)
	l := len(items) / batchNum
	r := len(items) % 5
	for i := 0; i < l; i++ {
		slc := items[i*batchNum : (i+1)*batchNum]
		batch = append(batch, slc)
	}

	batch = append(batch, items[len(items)-r:])

	return batch
}
