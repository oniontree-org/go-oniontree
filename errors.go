package oniontree

import "fmt"

type ErrNotOnionTree struct {
	dir string
}

func (e *ErrNotOnionTree) Error() string {
	return fmt.Sprintf("oniontree: directory `%s` is not an OnionTree repository", e.dir)
}

type ErrIdExists struct {
	id string
}

func (e *ErrIdExists) Error() string {
	return fmt.Sprintf("oniontree: service with ID `%s` already exists", e.id)
}

type ErrIdNotExists struct {
	id string
}

func (e *ErrIdNotExists) Error() string {
	return fmt.Sprintf("oniontree: service with ID `%s` does not exist", e.id)
}

type ErrTagNotExists struct {
	name string
}

func (e *ErrTagNotExists) Error() string {
	return fmt.Sprintf("oniontree: tag with name `%s` does not exist", e.name)
}

type ErrInvalidID struct {
	id      string
	pattern string
}

func (e *ErrInvalidID) Error() string {
	return fmt.Sprintf("oniontree: service ID `%s` does not match the pattern \"%s\"", e.id, e.pattern)
}
