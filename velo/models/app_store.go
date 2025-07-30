package models

type AppStore interface {
	GetItem(id string) (Item, error)
	GetItems() ([]Item, error)
	DeleteItem(id string) error
	UpdateItem(id string, update map[string]any) (Item, error)
}
