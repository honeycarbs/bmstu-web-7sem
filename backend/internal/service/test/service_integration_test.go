package test

//var (
//	insertAccountQuery = "INSERT INTO users (name, username, email, password_hash) VALUES ('%v', '%v', '%v', '%v') RETURNING id"
//)
//
//func TestAccountCreate(t *testing.T) {
//	testSuites := []struct {
//		testName      string
//		inAccount     model.Account
//		outAccount    model.Account
//		ExpectedError error
//	}{
//		{
//			testName:      "ValidAccountRegistration",
//			inAccount:     mother.AccountMother(),
//			outAccount:    mother.AccountMother(),
//			ExpectedError: nil,
//		},
//		{
//			testName:      "UserAlreadyExists",
//			inAccount:     mother.AccountMother(),
//			outAccount:    mother.AccountMother(),
//			ExpectedError: e.ClientAccountError,
//		},
//	}
//
//	logging.Init()
//	logger := logging.GetLogger()
//
//	for _, testSuite := range testSuites {
//		t.Run(testSuite.testName, func(t *testing.T) {
//
//			dbclient, err := integration.GetTestResource()
//			if err != nil {
//				t.Fatal(err)
//			}
//
//			repo := repository.NewAccountRepositoryImpl(dbclient, logger)
//			//if len(testSuite.prepDBActions) != 0 {
//			//
//			//}
//			service := account.NewService(repo, logger)
//
//			err = service.CreateAccount(&testSuite.inAccount)
//
//			assert.Equal(t, testSuite.ExpectedError, err)
//		})
//	}
//}
