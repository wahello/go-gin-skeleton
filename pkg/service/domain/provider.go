package domain

import "context"

type Provider struct {
	UUID      string `json:"uuid"`
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
}

// ProviderService provides methods for interacting with Provider service.
type ProviderService interface {
	CreateProvider(ctx context.Context, shortName, longName string) (*Provider, error)
	UpdateProvider(ctx context.Context, uuid, shortName, longName string) (*Provider, error)
	GetProviderByUUID(ctx context.Context, uuid string) (*Provider, error)
	GetProviders(ctx context.Context, limit int) ([]Provider, error)
	DeleteProviderByUUID(ctx context.Context, uuid string) error
}

// ProviderRepository provides methods for interacting with Provider repository.
type ProviderRepository interface {
	CreateProvider(ctx context.Context, p Provider) error
	UpdateProvider(ctx context.Context, p Provider) error
	DeleteProviderByUUID(ctx context.Context, uuid string) error
	GetProviderByUUID(ctx context.Context, uuid string) (*Provider, error)
	GetProviderByShortName(ctx context.Context, shortName string) (*Provider, error)
	GetProviders(ctx context.Context, limit int) ([]Provider, error)
}
