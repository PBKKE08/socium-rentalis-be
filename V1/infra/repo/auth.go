package repo

import (
	"database/sql"
	"errors"
	"github.com/PBKKE08/FP-BE/V1/entity"
	"github.com/rs/zerolog/log"
)

var (
	ErrInsertingUser = errors.New("error insert user")
	ErrNoUser        = errors.New("unknown user")
	ErrUnknown       = errors.New("internal error")
)

type AuthRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Insert(name, email, telephone, password, gender string) error {
	q := `INSERT INTO users(name, email, telephone, password, gender) values ($1, $2, $3, $4, $5)`

	result, err := r.db.Exec(q, name, email, telephone, password, gender)
	if err != nil {
		log.Debug().Msg(err.Error())
		return ErrInsertingUser
	}

	affected, err := result.RowsAffected()
	if err != nil {
		log.Debug().Err(err).Stack()
		return ErrInsertingUser
	}

	if affected == 0 {
		log.Debug().Err(err).Stack()
		return ErrInsertingUser
	}

	return nil
}

func (r *AuthRepository) GetByEmail(email string) (entity.User, error) {
	q := `SELECT id, name, telephone, email, password, gender FROM users WHERE email = $1`
	var user entity.User

	err := r.db.QueryRow(q, email).Scan(&user.ID, &user.Name, &user.Telephone, &user.Email, &user.Password, &user.Gender)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, nil
		}

		log.Debug().Err(err).Stack()
		return user, ErrUnknown
	}

	return user, nil
}