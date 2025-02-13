// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/store layer. But at some poin you will
// want auditing or something that isn't specific to the data/store layer.
package user

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wirkijowski/ultimate-go/business/data/store/user"
	"github.com/wirkijowski/ultimate-go/business/sys/auth"
	"go.uber.org/zap"
)

// Core manages the set of API's for user accesss.
type Core struct {
	log  *zap.SugaredLogger
	user user.Store
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, db *sqlx.DB) Core {
	return Core{
		log:  log,
		user: user.NewStore(log, db),
	}
}

func (c Core) Create(ctx context.Context, nu user.NewUser, now time.Time) (user.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	usr, err := c.user.Create(ctx, nu, now)
	if err != nil {
		return user.User{}, fmt.Errorf("create: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return usr, nil
}

// Update replaces a user document in the database.
func (c Core) Update(ctx context.Context, claims auth.Claims, userID string, uu user.UpdateUser, now time.Time) error {

	// PERFORM PRE BUSINESS OPERATIONS

	if err := c.user.Update(ctx, claims, userID, uu, now); err != nil {
		return fmt.Errorf("updated: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return nil
}

// Delete removes a user document fromthe database.
func (c Core) Delete(ctx context.Context, claims auth.Claims, userID string) error {

	// PERFORM PRE BUSINESS OPERATIONS

	if err := c.user.Delete(ctx, claims, userID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return nil
}

// Query retrieves a list of existing users from the database.
func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]user.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	users, err := c.user.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return users, nil

}

// QueryByID retrieves user by ID from the database.
func (c Core) QueryByID(ctx context.Context, claims auth.Claims, userID string) (user.User, error) {

	usr, err := c.user.QueryByID(ctx, claims, userID)
	if err != nil {
		return user.User{}, fmt.Errorf("query by ID: %w", err)
	}

	return usr, nil
}

// QueryByEmail retrieves user by email from the database.
func (c Core) QueryByEmail(ctx context.Context, claims auth.Claims, email string) (user.User, error) {

	usr, err := c.user.QueryByEmail(ctx, claims, email)
	if err != nil {
		return user.User{}, fmt.Errorf("query by email: %w", err)
	}

	return usr, nil
}

// Authenticate finds a user by their emial and verifies their password. On
// Success it returns a Claims User representing this user, THe claims can be
// used to generate a token for future authentication.
func (c Core) Authenticate(ctx context.Context, now time.Time, email, password string) (auth.Claims, error) {

	claims, err := c.user.Authenticate(ctx, now, email, password)
	if err != nil {
		return auth.Claims{}, fmt.Errorf("claims: %w", err)
	}

	return claims, nil
}
