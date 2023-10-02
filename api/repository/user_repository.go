package repository

import (
	"database/sql"
	"errors"
	"log"
	"math"
	"noda/api/data/model"
	"noda/api/data/transfer"
	"noda/api/data/types"
	"noda/failure"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Insert(next *transfer.UserCreation) (*transfer.User, error) {
	if yes, err := r.ExistsUserWithEmail(next.Email); err != nil {
		return nil, err
	} else if yes {
		return nil, failure.ErrSameEmail
	}
	query := `
	INSERT INTO "user" ("first_name", "middle_name", "last_name", "surname", "email", "password")
	     VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING "user_id" AS "id",
		          "first_name",
		          "middle_name",
		          "last_name",
		          "surname",
		          "picture_url",
		          "email",
		          "created_at",
		          "updated_at";`
	row, err := r.db.Query(query,
		next.FirstName, next.MiddleName, next.LastName, next.Surname, next.Email, next.Password)
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
	defer row.Close()
	user := transfer.User{}
	if err = sqlscan.ScanOne(&user, row); err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(userID string, up *transfer.UserUpdate) (bool, error) {
	if actual, err := r.SelectShallowUserByID(userID); err != nil {
		return false, err
	} else if actual.FirstName == up.FirstName &&
		actual.MiddleName == up.MiddleName &&
		actual.LastName == up.LastName &&
		actual.Surname == up.Surname {
		return false, nil
	}
	query := `
	   UPDATE "user"
	      SET "first_name" = COALESCE(NULLIF(trim($2), ''), "first_name"),
	          "middle_name" = COALESCE(NULLIF(trim($3), ''), "middle_name"),
	          "last_name" = COALESCE(NULLIF(trim($4), ''), "last_name"),
	          "surname" = COALESCE(NULLIF(trim($5), ''), "surname"),
						"updated_at" = 'now()'
			WHERE "user_id" = $1;`
	result, err := r.db.Exec(query, &userID, &up.FirstName, &up.MiddleName, &up.LastName, &up.Surname)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(failure.PQErrorToString(pqerr))
		}
		return false, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return count >= 1, nil
}

func (r *UserRepository) PromoteToAdmin(userID string) (bool, error) {
	if actual, err := r.SelectRawUserByID(userID); err != nil {
		return false, err
	} else if actual.Role == types.RoleAdmin {
		return false, nil
	}
	query := `
	UPDATE "user"
	   SET "role_id" = 1,
		     "updated_at" = 'now()'
	 WHERE "user_id" = $1;`
	result, err := r.db.Exec(query, &userID)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(failure.PQErrorToString(pqerr))
		}
		return false, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return count >= 1, nil
}

func (r *UserRepository) DegradeToNormalUser(userID string) (bool, error) {
	if actual, err := r.SelectRawUserByID(userID); err != nil {
		return false, err
	} else if actual.Role == types.RoleUser {
		return false, nil
	}
	query := `
	UPDATE "user"
	   SET "role_id" = 2,
		     "updated_at" = 'now()'
	 WHERE "user_id" = $1;`
	result, err := r.db.Exec(query, &userID)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(failure.PQErrorToString(pqerr))
		}
		return false, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return count >= 1, nil
}

func (r *UserRepository) Block(userID string) (bool, error) {
	if actual, err := r.SelectRawUserByID(userID); err != nil {
		return false, err
	} else if actual.IsBlocked {
		return false, nil
	}
	query := `
	UPDATE "user"
	   SET "is_blocked" = TRUE
	 WHERE "user_id" = $1;`
	result, err := r.db.Exec(query, &userID)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(failure.PQErrorToString(pqerr))
		}
		return false, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return count >= 1, nil
}

func (r *UserRepository) Unblock(userID string) (bool, error) {
	if actual, err := r.SelectRawUserByID(userID); err != nil {
		return false, err
	} else if !actual.IsBlocked {
		return false, nil
	}
	query := `
	UPDATE "user"
	   SET "is_blocked" = FALSE
	 WHERE "user_id" = $1;`
	result, err := r.db.Exec(query, &userID)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(failure.PQErrorToString(pqerr))
		}
		return false, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return count >= 1, nil
}

func (r *UserRepository) ExistsUserWithEmail(email string) (bool, error) {
	// TODO: Create an index on email.
	query := `
	SELECT "user_id"
	  FROM "user"
	 WHERE lower("email") = lower($1);`
	result, err := r.db.Exec(query, &email)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(failure.PQErrorToString(pqerr))
		}
		return false, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return count >= 1, nil
}

func (r *UserRepository) AssertUserExists(userID string) error {
	query := `
	SELECT "user_id"
	  FROM "user"
	 WHERE "user_id" = $1;`
	result, err := r.db.Exec(query, &userID)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(failure.PQErrorToString(pqerr))
		}
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	} else if count >= 1 {
		return nil
	}
	return failure.ErrNotFound
}

func (r *UserRepository) SelectAll(limit, page int64) (*[]*transfer.User, error) {
	maxValidBeforeOverflow := (math.MaxInt64 / limit) - 1
	if page > maxValidBeforeOverflow {
		page = maxValidBeforeOverflow
	}
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
	  FROM "user"
	 WHERE "is_blocked" IS FALSE
ORDER BY "created_at" DESC
   LIMIT $1
	OFFSET ($1 * ($2::BIGINT - 1));`
	rows, err := r.db.Query(query, &limit, &page)
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
	return &users, nil
}

func (r *UserRepository) SelectAllSettings(limit, page int64, userID string) (*[]*transfer.UserSetting, error) {
	err := r.AssertUserExists(userID)
	if err != nil {
		return nil, err
	}
	maxValidBeforeOverflow := (math.MaxInt64 / limit) - 1
	if page > maxValidBeforeOverflow {
		page = maxValidBeforeOverflow
	}
	query := `
    SELECT "us"."key",
           "df"."description",
           "us"."value",
           "us"."created_at",
           "us"."updated_at"
      FROM "user_setting" "us"
INNER JOIN "predefined_user_setting" "df"
        ON "us"."key" = "df"."key"
     WHERE "us"."user_id" = $1
  ORDER BY "created_at" DESC
     LIMIT $2
	  OFFSET ($2 * ($3::BIGINT - 1));`
	rows, err := r.db.Query(query, &userID, &limit, &page)
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
	settings := []*transfer.UserSetting{}
	if err = sqlscan.ScanAll(&settings, rows); err != nil {
		log.Println(err)
		return nil, err
	}
	return &settings, nil
}

func (r *UserRepository) SelectOneSetting(userID, settingKey string) (*transfer.UserSetting, error) {
	err := r.AssertUserExists(userID)
	if err != nil {
		return nil, err
	}
	query := `
    SELECT "us"."key",
           "df"."description",
           "us"."value",
           "us"."created_at",
           "us"."updated_at"
      FROM "user_setting" "us"
INNER JOIN "predefined_user_setting" "df"
        ON "us"."key" = "df"."key"
     WHERE "us"."user_id" = $1 AND
		        "us"."key" = $2;`
	result, err := r.db.Query(query, &userID, &settingKey)
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

func (r *UserRepository) SelectAllBlocked(limit, page int64) (*[]*transfer.User, error) {
	maxValidBeforeOverflow := (math.MaxInt64 / limit) - 1
	if page > maxValidBeforeOverflow {
		page = maxValidBeforeOverflow
	}
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
	  FROM "user"
	 WHERE "is_blocked" IS TRUE
ORDER BY "created_at" DESC
   LIMIT $1
	OFFSET ($1 * ($2::BIGINT - 1));`
	rows, err := r.db.Query(query, &limit, &page)
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
	return &users, nil
}

func (r *UserRepository) SelectShallowUserByEmail(email string) (*transfer.User, error) {
	user, err := r.SelectRawUserByEmail(email)
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

func (r *UserRepository) SelectShallowUserByID(id string) (*transfer.User, error) {
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
	  FROM "user"
	 WHERE "user_id" = $1;`
	row, err := r.db.Query(query, &id)
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
	defer row.Close()
	user := transfer.User{}
	if err := sqlscan.ScanOne(&user, row); err != nil {
		switch {
		default:
			log.Println(err)
			return nil, err
		case sqlscan.NotFound(err):
			return nil, failure.ErrNotFound
		}
	}
	return &user, nil
}

func (r *UserRepository) SelectRawUserByID(userID string) (*model.User, error) {
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
	  FROM "user"
	 WHERE "user_id" = $1;`

	row, err := r.db.Query(query, &userID)
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
	defer row.Close()

	user := model.User{}
	if err := sqlscan.ScanOne(&user, row); err != nil {
		switch {
		default:
			log.Println(err)
			return nil, err
		case sqlscan.NotFound(err):
			return nil, failure.ErrNotFound
		}
	}
	return &user, nil
}

func (r *UserRepository) SelectRawUserByEmail(email string) (*model.User, error) {
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
	  FROM "user"
	 WHERE lower("email") = lower($1);`
	row, err := r.db.Query(query, &email)
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
	defer row.Close()

	user := model.User{}
	if err := sqlscan.ScanOne(&user, row); err != nil {
		switch {
		default:
			log.Println(err)
			return nil, err
		case sqlscan.NotFound(err):
			return nil, failure.ErrNotFound
		}
	}
	return &user, nil
}

func (r *UserRepository) HardDelete(userID string) error {
	if err := r.AssertUserExists(userID); err != nil {
		return err
	}
	query := `
	DELETE FROM "user"
	      WHERE "user_id" = $1;`
	result, err := r.db.Exec(query, &userID)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(failure.PQErrorToString(pqerr))
		}
		return err
	}
	if _, err := result.RowsAffected(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *UserRepository) SoftDelete(id string) (string, error) {
	var (
		query = `
		DELETE FROM "user"
					WHERE "user_id" = $1
			RETURNING "user_id";`
		row           = r.db.QueryRow(query, id)
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
