package models

import (
	"lambda.sx/marcus/lambdago/sql"
	"log"
)

// Model is an interface that should be implemented by anything that should be saved to the database.
type Model interface {
	// TableName returns a string that will be used in the database as the name of the table.
	TableName() string
}

// save puts the current state of a model in the database
// TODO update existing instead of adding!
func Save(m Model) error {
	// Open the table for the model
	col, err := sql.Connection().Collection(m.TableName())
	if err != nil {
		log.Fatalf("sess.Collection(): %q\n", err)
	}
	// Add the model to the table
	col.Append(m)
	// If we failed, let the caller know
	return err
}
