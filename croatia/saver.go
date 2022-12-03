package croatia

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"missingPersons/common"
	"missingPersons/dataSource"
	"missingPersons/download"
	"missingPersons/types"
	"regexp"
	"strings"
)

func SaveCountry(people []common.RawPerson, country dataSource.Country) (types.Information, error) {
	fmt.Println("Croatia: Saving to database...")
	db := dataSource.DB()

	info := types.Information{
		UpdatedCount: 0,
		CreatedCount: 0,
		DeletedCount: 0,
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		for _, person := range people {
			id, err := createPersonId(person)

			if err != nil {
				return err
			}

			var dbPerson dataSource.Person
			if err := db.Where("custom_id = ?", id).First(&dbPerson).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if dbPerson.ID == "" {
				dbPerson = common.PersonFromRawPerson(id, country.ID, person)
			}

			if err != nil {
				return err
			}

			if person.ImageURL != "" {
				fileName := download.CreateImageName(person.ImageURL, id)
				ok, err := download.ImageExists(fileName)

				if !ok && err == nil {
					path, err := download.DownloadAndSaveImage(person.ImageURL, fileName)
					if err != nil {
						fmt.Println(fmt.Sprintf("Cannot download and save image: %s", err.Error()))
					} else {
						dbPerson.ImagePath = path
					}
				}
			}

			if dbPerson.ID == "" {
				if err := tx.Create(&dbPerson).Error; err != nil {
					return err
				}

				info.CreatedCount++
			} else {
				if err := tx.Save(&dbPerson).Error; err != nil {
					return err
				}

				info.UpdatedCount++
			}
		}

		return nil
	}); err != nil {
		return info, err
	}

	fmt.Println("Croatia: All records saved to database.")

	return info, nil
}

func createPersonId(person common.RawPerson) (string, error) {
	final := fmt.Sprintf("%s%s%s%s", person.Name, person.LastName, person.DOB, person.POB)
	final = strings.ToLower(final)

	re := regexp.MustCompile(`\s+`)
	final = re.ReplaceAllString(final, "")

	re = regexp.MustCompile(`\.`)
	final = re.ReplaceAllString(final, "")

	re = regexp.MustCompile(`,`)
	final = re.ReplaceAllString(final, "")

	return final, nil
}
