package mother

import (
	"log"
	"neatly/internal/model"
	"neatly/pkg/jwt"
	"time"
)

func TokenMother() string {
	a := AccountMother()
	token, err := jwt.GenerateAccessToken(a.ID)
	if err != nil {
		log.Fatal("can't create test token")
	}

	return token
}

func AccountMother() model.Account {

	testHash, _ := model.GeneratePasswordHash("testtest")

	return model.Account{
		ID:           0,
		Name:         "Test",
		Username:     "TestTest",
		Email:        "test",
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
