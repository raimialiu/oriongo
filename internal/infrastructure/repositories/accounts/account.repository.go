package accounts

import "oriongo/internal/infrastructure/repositories"

type (
	AccountRepository struct {
		_base repositories.BaseRepository
	}
)

func NewAccountRepository(base repositories.BaseRepository) *AccountRepository {
	return &AccountRepository{
		_base: base,
	}
}

func NewAccount() {

}
