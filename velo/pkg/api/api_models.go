package api

import "errors"

// CreateItemRequest represents the data for creating a new item.
type CreateItemRequest struct {
	Name string `json:"name"`
}

// Validate ensures the request data is valid.
func (c *CreateItemRequest) Validate() error {
	if c.Name == "" {
		return errors.New("name is a required field")
	}
	return nil
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

func (c *CreateUserRequest) Validate() error {
	if c.Name == "" {
		return errors.New("name ise a required field")
	}
	return nil
}
