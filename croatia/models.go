package croatia

type rawPerson struct {
	Name             string
	LastName         string
	MaidenName       string
	Gender           string
	DOB              string
	POB              string
	Citizenship      string
	PrimaryAddress   string
	SecondaryAddress string
	Country          string

	Height   string
	Hair     string
	EyeColor string

	// Date of disappearance
	DOD string
	// place of disappearance
	POD string

	Description string
}

func newRawPerson() rawPerson {
	return rawPerson{}
}
