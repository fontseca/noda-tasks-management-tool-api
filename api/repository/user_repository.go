package repository

import (
	"database/sql"
	"errors"
	"log"
	"noda/api/data/model"
	"noda/api/data/transfer"
	"noda/failure"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

type ur = UserRepository

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *ur) InsertUser(next *transfer.UserCreation) (string, error) {
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
				return "", failure.ErrSameEmail
			}
			log.Println(failure.PQErrorToString(pqerr))
		}
		return "", err
	}
	return insertedID, nil
}

func (r *ur) UpdateUser(userID string, up *transfer.UserUpdate) (bool, error) {
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
				return false, failure.ErrNotFound
			}
			log.Println(failure.PQErrorToString(pqerr))
		}
		return false, err
	}
	return wasUpdated, nil
}

func (r *ur) PromoteUserToAdmin(userID string) (bool, error) {
	row := r.db.QueryRow("SELECT promote_user_to_admin ($1);", userID)
	var wasPromoted bool
	if err := row.Scan(&wasPromoted); err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return false, failure.ErrNotFound
			}
			log.Println(failure.PQErrorToString(pqerr))
		}
		return false, err
	}
	return wasPromoted, nil
}

func (r *ur) DegradeAdminToNormalUser(userID string) (bool, error) {
	row := r.db.QueryRow("SELECT degrade_admin_to_user ($1);", userID)
	var wasDegraded bool
	if err := row.Scan(&wasDegraded); err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return false, failure.ErrNotFound
			}
			log.Println(failure.PQErrorToString(pqerr))
		}
		return false, err
	}
	return wasDegraded, nil
}

func (r *ur) BlockUser(userID string) (bool, error) {
	row := r.db.QueryRow("SELECT block_user ($1);", userID)
	var wasBlocked bool
	if err := row.Scan(&wasBlocked); err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return false, failure.ErrNotFound
			}
			log.Println(failure.PQErrorToString(pqerr))
		}
		return false, err
	}
	return wasBlocked, nil
}

func (r *ur) UnblockUser(userID string) (bool, error) {
	row := r.db.QueryRow("SELECT unblock_user ($1);", userID)
	var wasUnblocked bool
	if err := row.Scan(&wasUnblocked); err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return false, failure.ErrNotFound
			}
			log.Println(failure.PQErrorToString(pqerr))
		}
		return false, err
	}
	return wasUnblocked, nil
}

func (r *ur) FetchUsers(page, rpp int64) ([]*transfer.User, error) {
	query := `
	SELECT "user_id" AS "id",
	       "role_id" AS "role",
	       "first_name",
	       "middle_name",
	       "last_name",
	       "surname",
	       "picture_url",
	       "email",
	       "is_blocked",
	       "created_at",
	       "updated_at"
	  FROM fetch_users ($1, $2);`
	rows, err := r.db.Query(query, page, rpp)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(failure.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer rows.Close()
	users := []*transfer.User{}
	if err = sqlscan.ScanAll(&users, rows); err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}

func (r *ur) SearchUsers(page, rpp int64, needle, sortExpr string) ([]*transfer.User, error) {
	query := `
	SELECT "user_id" AS "id",
	       "role_id" AS "role",
	       "first_name",
	       "middle_name",
	       "last_name",
	       "surname",
	       "picture_url",
	       "email",
	       "is_blocked",
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
			log.Println(failure.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer rows.Close()
	users := []*transfer.User{}
	if err = sqlscan.ScanAll(&users, rows); err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}

func (r *ur) FetchUserSettings(userID string, page, rpp int64) ([]*transfer.UserSetting, error) {
	rows, err := r.db.Query("SELECT * FROM fetch_user_settings ($1, $2, $3);",
		userID, page, rpp)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			if isNonexistentUserError(pqerr) {
				return nil, failure.ErrNotFound
			}
			log.Println(failure.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer rows.Close()
	settings := []*transfer.UserSetting{}
	if err = sqlscan.ScanAll(&settings, rows); err != nil {
		log.Println(err)
		return nil, err
	}
	return settings, nil
}

func (r *ur) FetchOneUserSetting(userID, settingKey string) (*transfer.UserSetting, error) {
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
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return nil, failure.ErrNotFound
			case isNonexistentPredefinedUserSettingError(pqerr):
				return nil, failure.ErrSettingNotFound
			}
		}
		return nil, err
	}
	defer result.Close()
	setting := transfer.UserSetting{}
	if err = sqlscan.ScanOne(&setting, result); err != nil {
		if sqlscan.NotFound(err) {
			return nil, failure.ErrSettingNotFound
		}
		log.Println(err)
		return nil, err
	}
	return &setting, nil
}

func (r *ur) UpdateUserSetting(userID, settingKey string, value string) (bool, error) {
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
				return false, failure.ErrNotFound
			case isNonexistentPredefinedUserSettingError(pqerr):
				return false, failure.ErrSettingNotFound
			}
			log.Println(failure.PQErrorToString(pqerr))
		}
		return false, err
	}
	if wasUpdated {
		return true, nil
	}
	return false, nil
}

func (r *ur) FetchBlockedUsers(page, rpp int64) ([]*transfer.User, error) {
	query := `
	SELECT "user_id" AS "id",
	       "role_id" AS "role",
	       "first_name",
	       "middle_name",
	       "last_name",
	       "surname",
	       "picture_url",
	       "email",
	       "is_blocked",
	       "created_at",
	       "updated_at"
	  FROM fetch_blocked_users ($1, $2);`
	rows, err := r.db.Query(query, page, rpp)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(failure.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer rows.Close()
	users := []*transfer.User{}
	if err = sqlscan.ScanAll(&users, rows); err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}

func (r *ur) FetchUserByID(userID string) (*model.User, error) {
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
				 "is_blocked",
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
				return nil, failure.ErrNotFound
			}
			log.Println(failure.PQErrorToString(pqerr))
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

func (r *ur) FetchUserByEmail(email string) (*model.User, error) {
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
				 "is_blocked",
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
				return nil, failure.ErrNotFound
			}
			log.Println(failure.PQErrorToString(pqerr))
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

func (r *ur) FetchTransferUserByEmail(email string) (*transfer.User, error) {
	user, err := r.FetchUserByEmail(email)
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

func (r *ur) FetchTransferUserByID(userID string) (*transfer.User, error) {
	query := `
	SELECT "user_id" AS "id",
	       "role_id" AS "role",
	       "first_name",
	       "middle_name",
	       "last_name",
	       "surname",
	       "picture_url",
	       "email",
				 "is_blocked",
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
				return nil, failure.ErrNotFound
			}
			log.Println(failure.PQErrorToString(pqerr))
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

func (r *ur) HardlyDeleteUser(userID string) error {
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
				return failure.ErrNotFound
			}
			log.Println(failure.PQErrorToString(pqerr))
		}
		return err
	}
	return nil
}

func (r *ur) SoftlyDeleteUser(userID string) (string, error) {
	var (
		query = `
		DELETE FROM "user"
					WHERE "user_id" = $1
			RETURNING "user_id";`
		row           = r.db.QueryRow(query, userID)
		deletedUserID = ""
	)
	if err := row.Scan(&deletedUserID); err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
			return "", err
		case errors.As(err, &pqerr):
			log.Println(failure.PQErrorToString(pqerr))
			return "", err
		case errors.Is(err, sql.ErrNoRows):
			return "", failure.ErrNotFound
		}
	}
	return deletedUserID, nil
}
