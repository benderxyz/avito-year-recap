package profiles

type Profile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var testProfiles = []Profile{
	{ID: "1", Name: "Анна", Description: "Много смотрит недвижимость"},
	{ID: "2", Name: "Игорь", Description: "Активно продаёт вещи"},
	{ID: "3", Name: "Марина", Description: "Почти не пользовалась Авито"},
}

func List() []Profile {
	return testProfiles
}

func GetByID(id string) (Profile, bool) {
	for _, profile := range testProfiles {
		if profile.ID == id {
			return profile, true
		}
	}
	return Profile{}, false
}
