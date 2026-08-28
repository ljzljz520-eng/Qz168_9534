package tea17

import (
	"tea17/internal/api"
	"tea17/internal/catalog"
	"tea17/internal/notify"
	"tea17/internal/service"
	"tea17/internal/store"
)

type Application struct {
	Store   *store.Store
	Service *service.Service
	API     *api.Server
	Sender  *notify.MemorySender
}

func Open(path string) (*Application, error) {
	db, err := store.Open(path)
	if err != nil {
		return nil, err
	}
	sender := notify.NewMemorySender()
	svc := service.New(db, catalog.New(), sender)
	return &Application{Store: db, Service: svc, API: api.New(svc), Sender: sender}, nil
}
func (a *Application) Close() error { return a.Store.Close() }
