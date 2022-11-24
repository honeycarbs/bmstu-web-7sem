package mother

import (
	"log"
	"neatly/internal/model"
	"neatly/pkg/jwt"
	"neatly/pkg/logging"
	"os"
	"time"
)

func TokenMother() string {
	a := AccountMother()

	logging.Init()
	os.Setenv("CONF_FILE", "../../../etc/config/local.yml")

	token, err := jwt.GenerateAccessToken(a.ID)
	if err != nil {
		log.Fatal("can't create sqlite token")
	}

	return token
}

func AccountMother() model.Account {

	testHash, _ := model.GeneratePasswordHash("testtest")

	return model.Account{
		ID:           0,
		Name:         "Test",
		Username:     "TestTest",
		Email:        "sqlite",
		Password:     "testtest",
		PasswordHash: testHash,
	}
}

func NoteMother() model.Note {
	return model.Note{
		ID:        0,
		Header:    "",
		Body:      "",
		ShortBody: "",
		Tags:      nil,
		Color:     "",
		Edited:    time.Time{},
	}
}

func TagMother() model.Tag {
	return model.Tag{
		ID:    0,
		Label: "",
		Color: "",
	}
}
