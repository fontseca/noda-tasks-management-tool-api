package repository

import (
	"database/sql"
	"errors"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/lib/pq"
	"log"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
)

type UserRepository interface {
	Save(creation *transfer.UserCreation) (insertedID string, err error)
	FetchByID(id string) (user *model.User, err error)
	FetchShallowUserByID(id string) (user *transfer.User, err error)
	FetchByEmail(email string) (user *model.User, err error)
	FetchShallowUserByEmail(email string) (user *transfer.User, err error)
	Fetch(page, rpp int64, needle, sortExpr string) (users []*transfer.User, err error)
	FetchBlocked(page, rpp int64, needle, sortExpr string) (users []*transfer.User, err error)
	FetchSettings(userID string, page, rpp int64, needle, sortExpr string) (settings []*transfer.UserSetting, err error)
	FetchOneSetting(userID string, settingKey string) (setting *transfer.UserSetting, err error)
	Search(page, rpp int64, needle, sortExpr string) (users []*transfer.User, err error)
	Update(id string, update *transfer.UserUpdate) (ok bool, err error)
	UpdateUserSetting(userID, settingKey, newValue string) (ok bool, err error)
	Block(id string) (ok bool, err error)
	Unblock(id string) (ok bool, err error)
	PromoteToAdmin(id string) (ok bool, err error)
	DegradeToUser(id string) (ok bool, err error)
	RemoveHardly(id string) error
	RemoveSoftly(id string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r userRepository) Save(next *transfer.UserCreation) (string, error) {
	row := r.db.QueryRow("SELECT make_user ($1, $2, $3, $4, $5, $6);",
		next.FirstName, next.MiddleName, next.LastName, next.Surname, next.Email, next.Password)
	var insertedID string
	if err := row.Scan(&insertedID); err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isDuplicatedEmailError(pqerr) {
				return "", noda.ErrSameEmail
			}
			log.Println(noda.PQErrorToString(pqerr))
		}
		return "", err
	}
	return insertedID, nil
}

func (r userRepository) Update(userID string, up *transfer.UserUpdate) (bool, error) {
	row := r.db.QueryRow("SELECT update_user ($1, $2, $3, $4, $5, NULL, NULL, NULL);",
		userID, up.FirstName, up.MiddleName, up.LastName, up.Surname)
	var wasUpdated bool
	if err := row.Scan(&wasUpdated); err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return false, noda.ErrUserNotFound
			}
			log.Println(noda.PQErrorToString(pqerr))
		}
		return false, err
	}
	return wasUpdated, nil
}

func (r userRepository) PromoteToAdmin(userID string) (bool, error) {
	row := r.db.QueryRow("SELECT promote_user_to_admin ($1);", userID)
	var wasPromoted bool
	if err := row.Scan(&wasPromoted); err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return false, noda.ErrUserNotFound
			}
			log.Println(noda.PQErrorToString(pqerr))
		}
		return false, err
	}
	return wasPromoted, nil
}

func (r userRepository) DegradeToUser(userID string) (bool, error) {
	row := r.db.QueryRow("SELECT degrade_admin_to_user ($1);", userID)
	var wasDegraded bool
	if err := row.Scan(&wasDegraded); err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return false, noda.ErrUserNotFound
			}
			log.Println(noda.PQErrorToString(pqerr))
		}
		return false, err
	}
	return wasDegraded, nil
}

func (r userRepository) Block(userID string) (bool, error) {
	row := r.db.QueryRow("SELECT block_user ($1);", userID)
	var wasBlocked bool
	if err := row.Scan(&wasBlocked); err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return false, noda.ErrUserNotFound
			}
			log.Println(noda.PQErrorToString(pqerr))
		}
		return false, err
	}
	return wasBlocked, nil
}

func (r userRepository) Unblock(userID string) (bool, error) {
	row := r.db.QueryRow("SELECT unblock_user ($1);", userID)
	var wasUnblocked bool
	if err := row.Scan(&wasUnblocked); err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return false, noda.ErrUserNotFound
			}
			log.Println(noda.PQErrorToString(pqerr))
		}
		return false, err
	}
	return wasUnblocked, nil
}

func (r userRepository) Fetch(page, rpp int64, needle, sortExpr string) ([]*transfer.User, error) {
	query := `
	SELECT "user_id" AS "id",
	       "role_id" AS "role",
	       "first_name",
	       "middle_name",
	       "last_name",
	       "surname",
	       "picture_url",
	       "email",
	       "created_at",
	       "updated_at"
	  FROM fetch_users ($1, $2, $3, $4);`
	rows, err := r.db.Query(query, page, rpp, needle, sortExpr)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(noda.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer rows.Close()
	var users = make([]*transfer.User, 0)
	if err = sqlscan.ScanAll(&users, rows); err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}

func (r userRepository) Search(page, rpp int64, needle, sortExpr string) ([]*transfer.User, error) {
	query := `
	SELECT "user_id" AS "id",
	       "role_id" AS "role",
	       "first_name",
	       "middle_name",
	       "last_name",
	       "surname",
	       "picture_url",
	       "email",
	       "created_at",
	       "updated_at"
	  FROM fetch_users ($1, $2, $3, $4);`
	rows, err := r.db.Query(query, page, rpp, needle, sortExpr)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(noda.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer rows.Close()
	var users = make([]*transfer.User, 0)
	if err = sqlscan.ScanAll(&users, rows); err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}

func (r userRepository) FetchSettings(userID string, page, rpp int64, needle, sortExpr string) ([]*transfer.UserSetting, error) {
	rows, err := r.db.Query("SELECT * FROM fetch_user_settings ($1, $2, $3, $4, $5);",
		userID, page, rpp, needle, sortExpr)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return nil, noda.ErrUserNotFound
			}
			log.Println(noda.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer rows.Close()
	var settings = make([]*transfer.UserSetting, 0)
	if err = sqlscan.ScanAll(&settings, rows); err != nil {
		log.Println(err)
		return nil, err
	}
	return settings, nil
}

func (r userRepository) FetchOneSetting(userID, settingKey string) (*transfer.UserSetting, error) {
	result, err := r.db.Query("SELECT * FROM fetch_one_user_setting ($1, $2);",
		userID, settingKey)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			switch {
			default:
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return nil, noda.ErrUserNotFound
			case isNonexistentPredefinedUserSettingError(pqerr):
				return nil, noda.ErrSettingNotFound
			}
		}
		return nil, err
	}
	defer result.Close()
	setting := transfer.UserSetting{}
	if err = sqlscan.ScanOne(&setting, result); err != nil {
		if sqlscan.NotFound(err) {
			return nil, noda.ErrSettingNotFound
		}
		log.Println(err)
		return nil, err
	}
	return &setting, nil
}

func (r userRepository) UpdateUserSetting(userID, settingKey string, value string) (bool, error) {
	query := `SELECT update_user_setting ($1, $2, $3);`
	row := r.db.QueryRow(query, userID, settingKey, value)
	var wasUpdated bool
	if err := row.Scan(&wasUpdated); err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			switch {
			case isNonexistentUserError(pqerr):
				return false, noda.ErrUserNotFound
			case isNonexistentPredefinedUserSettingError(pqerr):
				return false, noda.ErrSettingNotFound
			}
			log.Println(noda.PQErrorToString(pqerr))
		}
		return false, err
	}
	if wasUpdated {
		return true, nil
	}
	return false, nil
}

func (r userRepository) FetchBlocked(page, rpp int64, needle, sortExpr string) ([]*transfer.User, error) {
	query := `
	SELECT "user_id" AS "id",
	       "role_id" AS "role",
	       "first_name",
	       "middle_name",
	       "last_name",
	       "surname",
	       "picture_url",
	       "email",
	       "created_at",
	       "updated_at"
	  FROM fetch_blocked_users ($1, $2, $3, $4);`
	rows, err := r.db.Query(query, page, rpp, needle, sortExpr)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(noda.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer rows.Close()
	var users = make([]*transfer.User, 0)
	if err = sqlscan.ScanAll(&users, rows); err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}

func (r userRepository) FetchByID(userID string) (*model.User, error) {
	query := `
	SELECT "user_id" AS "id",
	       "role_id" AS "role",
	       "first_name",
	       "middle_name",
	       "last_name",
	       "surname",
	       "picture_url",
	       "email",
				 "password",
	       "created_at",
	       "updated_at"
	  FROM fetch_user_by_id ($1);`
	row, err := r.db.Query(query, userID)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return nil, noda.ErrUserNotFound
			}
			log.Println(noda.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer row.Close()

	user := model.User{}
	if err := sqlscan.ScanOne(&user, row); err != nil {
		switch {
		default:
			log.Println(err)
			return nil, err
		}
	}
	return &user, nil
}

func (r userRepository) FetchByEmail(email string) (*model.User, error) {
	query := `
	SELECT "user_id" AS "id",
	       "role_id" AS "role",
	       "first_name",
	       "middle_name",
	       "last_name",
	       "surname",
	       "picture_url",
	       "email",
				 "password",
	       "created_at",
	       "updated_at"
	  FROM fetch_user_by_email ($1);`
	row, err := r.db.Query(query, email)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNotFoundEmailError(pqerr) {
				return nil, noda.ErrUserNotFound
			}
			log.Println(noda.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer row.Close()

	user := model.User{}
	if err := sqlscan.ScanOne(&user, row); err != nil {
		switch {
		default:
			log.Println(err)
			return nil, err
		}
	}
	return &user, nil
}

func (r userRepository) FetchShallowUserByEmail(email string) (*transfer.User, error) {
	user, err := r.FetchByEmail(email)
	if err != nil {
		return nil, err
	}
	return &transfer.User{
		ID:         user.ID,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Surname:    user.Surname,
		PictureUrl: user.PictureUrl,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}, nil
}

func (r userRepository) FetchShallowUserByID(userID string) (*transfer.User, error) {
	query := `
	SELECT "user_id" AS "id",
	       "role_id" AS "role",
	       "first_name",
	       "middle_name",
	       "last_name",
	       "surname",
	       "picture_url",
	       "email",
	       "created_at",
	       "updated_at"
	  FROM fetch_user_by_id ($1);`
	row, err := r.db.Query(query, userID)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return nil, noda.ErrUserNotFound
			}
			log.Println(noda.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer row.Close()
	user := transfer.User{}
	if err := sqlscan.ScanOne(&user, row); err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}

func (r userRepository) RemoveHardly(userID string) error {
	err := r.db.
		QueryRow("SELECT delete_user_hardly ($1);", userID).
		Err()
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return noda.ErrUserNotFound
			}
			log.Println(noda.PQErrorToString(pqerr))
		}
		return err
	}
	return nil
}

func (r userRepository) RemoveSoftly(userID string) error {
	var (
		query = `
		DELETE FROM "user"
					WHERE "user_id" = $1;`
		row = r.db.QueryRow(query, userID)
	)
	var err = row.Err()
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
			return err
		case errors.As(err, &pqerr):
			log.Println(noda.PQErrorToString(pqerr))
			return err
		case errors.Is(err, sql.ErrNoRows):
			return noda.ErrUserNotFound
		}
	}
	return nil
}
