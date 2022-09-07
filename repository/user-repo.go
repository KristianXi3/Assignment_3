package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/KristianXi3/Assignment_3/entity"
	"github.com/KristianXi3/Assignment_3/helper"
	"github.com/KristianXi3/Assignment_3/model"
)

type UserRepoIface interface {
	GetUsers(ctx context.Context) ([]*model.User, error)
	GetUserById(ctx context.Context, id int) (*model.User, error)
	CreateUser(ctx context.Context, user entity.User) (string, error)
	UpdateUser(ctx context.Context, id int, user model.User) (string, error)
	DeleteUser(ctx context.Context, id int) (string, error)
	LoginUser(ctx context.Context, login model.Login) (*entity.User, string)
}

type UserRepo struct {
	sql *sql.DB
}

func NewUserRepo(context *sql.DB) UserRepoIface {
	return &UserRepo{sql: context}
}

func (repo *UserRepo) LoginUser(ctx context.Context, loginUser model.Login) (*entity.User, string) {
	var user entity.User

	data, err := repo.sql.QueryContext(ctx, "SELECT id, username, password, email, age FROM USERS WHERE email=@email",
		sql.Named("email", loginUser.Email))
	defer data.Close()
	if err != nil {
		log.Fatal(err)
		return nil, "User Not Found"
	}

	for data.Next() {
		err := data.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Age)
		if err != nil {
			log.Fatal(err)
			return nil, "Get User Failed"
		}
	}

	check := helper.CheckPasswordHash(loginUser.Password, user.Password)
	if !check {
		return nil, "Invalid Password"
	}
	return &user, ""
}

func (repo *UserRepo) GetUsers(ctx context.Context) ([]*model.User, error) {
	users := []*model.User{}

	data, err := repo.sql.QueryContext(ctx, "SELECT id, username, email, age FROM USERS")
	defer data.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for data.Next() {
		var user model.User
		err := data.Scan(&user.Id, &user.Username, &user.Email, &user.Age)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (repo *UserRepo) GetUserById(ctx context.Context, id int) (*model.User, error) {
	var user model.User

	data, err := repo.sql.QueryContext(ctx, "SELECT id, username, email, age FROM USERS WHERE Id = @Id",
		sql.Named("Id", id))
	defer data.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for data.Next() {
		err := data.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Age)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	return &user, nil
}

func (repo *UserRepo) CreateUser(ctx context.Context, user entity.User) (string, error) {
	var result string
	var err error

	user.Password, err = helper.GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatal("Error in password hashing")
		return "", err
	}

	_, err = repo.sql.ExecContext(ctx, "INSERT into USERS (username, email, password, age) values (@username, @email, @password, @age)",
		sql.Named("username", user.Username),
		sql.Named("email", user.Email),
		sql.Named("password", user.Password),
		sql.Named("age", user.Age))

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	result = "User created successfully"
	return result, nil
}

func (repo *UserRepo) UpdateUser(ctx context.Context, id int, user model.User) (string, error) {
	var result string

	_, err := repo.sql.ExecContext(ctx, "UPDATE USERS set username = @username, email = @email, age = @age where id = @id",
		sql.Named("id", id),
		sql.Named("username", user.Username),
		sql.Named("email", user.Email),
		sql.Named("age", user.Age))

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	result = "User updated successfully"
	return result, nil
}

func (repo *UserRepo) DeleteUser(ctx context.Context, id int) (string, error) {
	var result string

	_, err := repo.sql.ExecContext(ctx, "DELETE from USERS where id=@id",
		sql.Named("id", id))

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	result = fmt.Sprintf("User %d deleted", id)
	return result, nil
}
