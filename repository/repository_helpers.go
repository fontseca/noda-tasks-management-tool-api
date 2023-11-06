package repository

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func isNonexistentUserError(err *pq.Error) bool {
	return err.Code == "P0001" &&
		strings.Contains(err.Message, "nonexistent user with ID")
}

func isNonexistentGroupError(err *pq.Error) bool {
	return err.Code == "P0001" &&
		strings.Contains(err.Message, "nonexistent group with ID")
}

func isNonexistentListError(err *pq.Error) bool {
	return err.Code == "P0001" &&
		strings.Contains(err.Message, "nonexistent list with ID")
}

func isContextDeadlineError(err error) bool {
	return strings.Compare(err.Error(), "context deadline exceeded") == 0
}

func isNonexistentPredefinedUserSettingError(err *pq.Error) bool {
	return err.Code == "P0001" &&
		strings.Contains(err.Message, "nonexistent predefined user setting key")
}

func isNotFoundEmailError(err *pq.Error) bool {
	return err.Code == "P0001" &&
		strings.Contains(err.Message, "nonexistent user email")
}

func isDuplicatedEmailError(err *pq.Error) bool {
	return err.Code == "23505" &&
		strings.Contains(err.Message, "duplicate key value violates unique constraint \"user_email_key\"")
}

func beQuiet() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	log.SetOutput(null)
	return func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}
