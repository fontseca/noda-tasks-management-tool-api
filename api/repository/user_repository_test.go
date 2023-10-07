package repository

import (
	"database/sql"
	"log"
	"noda/api/data/model"
	"noda/api/data/transfer"
	"noda/api/data/types"
	"noda/failure"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var next *transfer.UserCreation = &transfer.UserCreation{
	FirstName:  "Jeremy",
	MiddleName: "Alexander",
	LastName:   "Fonseca",
	Surname:    "Blanco",
	Email:      "fonseca@mail.com",
	Password:   "passowrd1234",
}

var raw *model.User = &model.User{
	ID:         uuid.New(),
	FirstName:  next.FirstName,
	MiddleName: next.MiddleName,
	LastName:   next.LastName,
	Surname:    next.Surname,
	Email:      next.Email,
	Password:   next.Password,
	Role:       types.RoleUser,
	PictureUrl: nil,
	IsBlocked:  false,
	CreatedAt:  time.Now(),
	UpdatedAt:  time.Now(),
}

var shallow *transfer.User = &transfer.User{
	ID:         raw.ID,
	FirstName:  raw.FirstName,
	MiddleName: raw.MiddleName,
	LastName:   raw.LastName,
	Surname:    raw.Surname,
	Email:      raw.Email,
	Role:       raw.Role,
	PictureUrl: raw.PictureUrl,
	IsBlocked:  raw.IsBlocked,
	CreatedAt:  raw.CreatedAt,
	UpdatedAt:  raw.UpdatedAt,
}

func quiet() func() {
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

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestInsertUser_ValidInput(t *testing.T) {
	db, mock := NewMock()
	r := NewUserRepository(db)
	query := `SELECT make_user ($1, $2, $3, $4, $5, $6);`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(next.FirstName, next.MiddleName, next.LastName, next.Surname, next.Email, next.Password).
		WillReturnRows(sqlmock.
			NewRows([]string{"make_user"}).
			AddRow(raw.ID))
	res, err := r.InsertUser(next)
	assert.NoError(t, err)
	assert.Equal(t, res, shallow.ID.String())
}

func TestInsertUser_InvalidEmail(t *testing.T) {
	defer quiet()()
	n := transfer.UserCreation{Email: "invalid-email"}
	db, mock := NewMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := `SELECT make_user ($1, $2, $3, $4, $5, $6);`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(n.FirstName, n.MiddleName, n.LastName, n.Surname, n.Email, n.Password).
		WillReturnError(&pq.Error{Code: "23514", Message: "value for domain email_t violates check constraint \"email_t_check\""})
	res, err := r.InsertUser(&n)
	assert.Error(t, err)
	assert.Equal(t, "", res)
}

func TestInsertUser_DuplicatedEmail(t *testing.T) {
	defer quiet()()
	n := transfer.UserCreation{Email: "mail@mail.com"}
	db, mock := NewMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := `SELECT make_user ($1, $2, $3, $4, $5, $6);`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(n.FirstName, n.MiddleName, n.LastName, n.Surname, n.Email, n.Password).
		WillReturnError(&pq.Error{Code: "23505", Message: "duplicate key value violates unique constraint \"user_email_key\""})
	res, err := r.InsertUser(&n)
	assert.ErrorIs(t, err, failure.ErrSameEmail)
	assert.Equal(t, "", res)
}

func TestInsertUser_UnexpectedDatabaseError(t *testing.T) {
	defer quiet()()
	n := transfer.UserCreation{}
	db, mock := NewMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := `SELECT make_user ($1, $2, $3, $4, $5, $6);`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(n.FirstName, n.MiddleName, n.LastName, n.Surname, n.Email, n.Password).
		WillReturnError(&pq.Error{})
	res, err := r.InsertUser(&n)
	assert.Error(t, err)
	assert.Equal(t, "", res)
}

func TestUpdateUser_ValidInput(t *testing.T) {
	up := transfer.UserUpdate{}
	db, mock := NewMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := "SELECT update_user ($1, $2, $3, $4, $5, NULL, NULL, NULL);"
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(raw.ID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
		WillReturnRows(sqlmock.
			NewRows([]string{"update_user"}).
			AddRow(true))
	res, err := r.UpdateUser(raw.ID.String(), &up)
	assert.NoError(t, err)
	assert.Equal(t, res, true)
}

func TestUpdateUser_FailureWithoutError(t *testing.T) {
	up := transfer.UserUpdate{}
	db, mock := NewMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := "SELECT update_user ($1, $2, $3, $4, $5, NULL, NULL, NULL);"
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(raw.ID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
		WillReturnRows(sqlmock.
			NewRows([]string{"update_user"}).
			AddRow(false))
	res, err := r.UpdateUser(raw.ID.String(), &up)
	assert.NoError(t, err)
	assert.Equal(t, res, false)
}

func TestUpdateUser_NonexistentUserID(t *testing.T) {
	up := transfer.UserUpdate{}
	db, mock := NewMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := "SELECT update_user ($1, $2, $3, $4, $5, NULL, NULL, NULL);"
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(raw.ID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err := r.UpdateUser(raw.ID.String(), &up)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.Equal(t, res, false)
}

func TestUpdateUser_UnexpectedDatabaseError(t *testing.T) {
	defer quiet()()
	up := transfer.UserUpdate{}
	db, mock := NewMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := "SELECT update_user ($1, $2, $3, $4, $5, NULL, NULL, NULL);"
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(raw.ID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
		WillReturnError(&pq.Error{})
	res, err := r.UpdateUser(raw.ID.String(), &up)
	assert.Error(t, err)
	assert.Equal(t, res, false)
}

func TestPromoteUserToAdmin_ValidInput(t *testing.T) {
	userID := raw.ID
	db, mock := NewMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := "SELECT promote_user_to_admin ($1);"
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID).
		WillReturnRows(sqlmock.
			NewRows([]string{"promote_user_to_admin"}).
			AddRow(true))
	got, err := r.PromoteUserToAdmin(userID.String())
	assert.NoError(t, err)
	assert.Equal(t, got, true)
}
