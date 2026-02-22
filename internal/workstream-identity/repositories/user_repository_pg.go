package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/danilobml/workstream/internal/workstream-identity/db"
)

type UserPgRepository struct {
	db db.DBInterface
}

func NewUserPgRepository(db db.DBInterface) *UserPgRepository {
	return &UserPgRepository{
		db: db,
	}
}

func (ur *UserPgRepository) List(ctx context.Context) ([]*models.User, error) {
	query := `SELECT id, email, is_active
				FROM users`
	rows, err := ur.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(&user.ID, &user.Email, &user.IsActive)
		if err != nil {
			return nil, err
		}

		userRoles, err := ur.getUserRoles(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		user.Roles = userRoles

		users = append(users, user)
	}

	return users, nil
}

func (ur *UserPgRepository) FindById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User

	err := ur.db.QueryRow(ctx, `
		SELECT id, email, hashed_password, is_active
			FROM users
			WHERE id = $1
		`, id).Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	roles, err := ur.getUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return &user, nil
}

func (ur *UserPgRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := ur.db.QueryRow(ctx, `
		SELECT id, email, hashed_password, is_active
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	roles, err := ur.getUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return &user, nil
}

func (ur *UserPgRepository) Create(ctx context.Context, user models.User) error {
	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, email, hashed_password, is_active)
		VALUES ($1, $2, $3, $4)
	`, user.ID, user.Email, user.HashedPassword, true)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errs.ErrAlreadyExists
		}
		return err
	}

	for _, role := range user.Roles {
		_, err = tx.Exec(ctx, `
			INSERT INTO user_roles (user_id, role)
			VALUES ($1, $2)
		`, user.ID, role.GetName())
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (ur *UserPgRepository) Update(ctx context.Context, user models.User) error {
	tx, err := ur.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE users 
		SET email = $1, hashed_password = $2, is_active = $3
		WHERE id = $4
	`, user.Email, user.HashedPassword, user.IsActive, user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errs.ErrAlreadyExists
		}
		return err
	}

	_, err = tx.Exec(ctx, `
		DELETE FROM user_roles
		WHERE user_id = $1
	`, user.ID)
	if err != nil {
		return err
	}

	for _, role := range user.Roles {
		_, err = tx.Exec(ctx, `
			INSERT INTO user_roles (user_id, role)
			VALUES ($1, $2)
		`, user.ID, role.GetName())
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (ur *UserPgRepository) Delete(ctx context.Context, id uuid.UUID) error {

	return nil
}

func (ur *UserPgRepository) getUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error) {
	rows, err := ur.db.Query(ctx, `
		SELECT role
		FROM user_roles
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]models.Role, 0)

	for rows.Next() {
		var roleStr string
		if err := rows.Scan(&roleStr); err != nil {
			return nil, err
		}

		role, err := models.ParseRole(roleStr)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}
