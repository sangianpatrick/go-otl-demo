package album

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sangianpatrick/go-otl-demo/exception"
	"github.com/sirupsen/logrus"
)

var (
	albumPath = "/albums"
)

type AlbumRepository interface {
	FindMany(ctx context.Context) (bunchOfAlbums []Album, err error)
}

type albumRepositoryImpl struct {
	logger *logrus.Logger
	host   string
	c      *http.Client
}

func (r *albumRepositoryImpl) FindMany(ctx context.Context) (bunchOfAlbums []Album, err error) {
	url := fmt.Sprintf("%s%s", r.host, albumPath)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	resp, err := r.c.Do(req)
	if err != nil {
		r.logger.WithContext(ctx).WithError(err).Error(err.Error())
		err = exception.ErrBadGateway
		return
	}

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode > 500 {
		err = exception.ErrBadGateway
		r.logger.WithContext(ctx).WithError(err).WithFields(logrus.Fields{
			"body": string(bodyBytes),
		}).Error(err.Error())
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = exception.ErrBadGateway
		r.logger.WithContext(ctx).WithError(err).WithFields(logrus.Fields{
			"body": string(bodyBytes),
		}).Error(err.Error())
		return
	}

	json.Unmarshal(bodyBytes, &bunchOfAlbums)

	return
}

func NewAlbumRepository(logger *logrus.Logger, c *http.Client, host string) AlbumRepository {
	return &albumRepositoryImpl{
		logger: logger,
		host:   host,
		c:      c,
	}
}
