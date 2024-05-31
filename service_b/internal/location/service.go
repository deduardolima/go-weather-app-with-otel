package location

import "context"

type LocationService interface {
	GetLocation(ctx context.Context, cep string) (string, error)
}

type locationService struct {
	client *LocationClient
}

func NewLocationService() LocationService {
	return &locationService{
		client: NewLocationClient(),
	}
}

func (ls *locationService) GetLocation(ctx context.Context, cep string) (string, error) {
	return ls.client.GetLocation(ctx, cep)
}
