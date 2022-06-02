package album

import (
	"context"
	"time"

	"github.com/sangianpatrick/go-otl-demo/exception"
	"github.com/sangianpatrick/go-otl-demo/response"
)

type AlbumService interface {
	GetMany(ctx context.Context) (resp response.Response)
}

type albumServiceImpl struct {
	repository AlbumRepository
}

func (s *albumServiceImpl) GetMany(ctx context.Context) (resp response.Response) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	bunchOfAlbums, err := s.repository.FindMany(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			return response.ResponseError(response.StatusRequestTimeout, err, nil, nil, "")
		}

		if err == exception.ErrBadGateway {
			return response.ResponseError(response.StatusBadGateway, err, nil, nil, "")
		}

		return response.ResponseError(response.StatusInternalServerError, err, nil, nil, "")
	}

	return response.ResponseSuccess(response.StatusOK, bunchOfAlbums, nil, "bunch of albums")
}

func NewAlbumService(repository AlbumRepository) (service AlbumService) {
	return &albumServiceImpl{repository: repository}
}
