package croatia

import "missingPersons/dataSource"

func buildFieldMap() map[string]string {
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

func personFromRawPerson(id string, countryId string, raw rawPerson) dataSource.Person {
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
	person.Description = raw.Name

	return person
}
