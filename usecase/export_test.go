package usecase

var GenerateRandomString = generateRandomString

func (u *usecase) SetGenerateShortURL(genURLFunc func(n int) (string, error)) {
	u.generateShortURL = genURLFunc
}
