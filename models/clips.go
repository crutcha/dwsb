package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"time"
)

type Clip struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name      string    `json:"name" db:"name"`
	Tag       string    `json:"tag" db:"tag"`
	File      string    `json:"file" db:"file"`
	Guild     string    `json:"guild" db:"guild"`
	User      User      `json:"user" has_one:"user" fk_id:"id"`
}

// String is not required by pop and may be deleted
func (c Clip) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Clips is not required by pop and may be deleted
type Clips []Clip

// String is not required by pop and may be deleted
func (c Clips) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Clip) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Name, Name: "Name"},
		&validators.StringIsPresent{Field: c.Tag, Name: "Tag"},
		&validators.StringIsPresent{Field: c.File, Name: "File"},
		&validators.StringIsPresent{Field: c.Guild, Name: "Guild"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Clip) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Clip) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
