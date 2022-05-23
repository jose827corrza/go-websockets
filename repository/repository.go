package repository

import (
	"context"

	"github.com/jose827corrza/go-websockets/models"
)

type Repository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	InsertPost(ctx context.Context, post *models.Post) error
	GetPostById(ctx context.Context, id string) (*models.Post, error)
	UpdatePostById(ctx context.Context, post *models.Post) error
	DeletePostById(ctx context.Context, id string, userId string) error
	ListPosts(ctx context.Context, page uint64, limit uint64) ([]*models.Post, error)
	Close() error
}

/*
PATRON REPOSITORY
Se lleva muy bien con los principios SOLID, por que abstrae lo mas que se puede
el codigo, independientemente de que persistencia de datos se va a utilizar,
digamos aca solo se ven las interfaces del crud, mas no el como se ejecute para la DB.
*/
//+++Abstracciones
//---Concretas

var implementation Repository

func SetRepository(repository Repository) {
	implementation = repository
}

func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}
func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return implementation.GetUserById(ctx, id)
}
func Close() error {
	return implementation.Close()
}
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, email)
}

func InsertPost(ctx context.Context, post *models.Post) error {
	return implementation.InsertPost(ctx, post)
}
func GetPostById(ctx context.Context, id string) (*models.Post, error) {
	return implementation.GetPostById(ctx, id)
}
func UpdatePostById(ctx context.Context, post *models.Post) error {
	return implementation.UpdatePostById(ctx, post)
}
func DeletePostById(ctx context.Context, id string, userId string) error {
	return implementation.DeletePostById(ctx, id, userId)
}
func ListPosts(ctx context.Context, page uint64, limit uint64) ([]*models.Post, error) {
	return implementation.ListPosts(ctx, page, limit)
}
