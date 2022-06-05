package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

type users struct {
	db *sql.DB
}

//NewUserRepository create a user repository
func NewUserRepository(db *sql.DB) *users {
	return &users{db}
}

//Create insert a user in database
func (repository users) Create(user models.User) (uint64, error) {
	statement, err := repository.db.Prepare(
		"insert into users (name, nick, email, password) values(?, ?, ?, ?)",
	)

	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(user.Name, user.Nick, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastInsertId), nil
}

//Find return all user filtered by name or nick
func (repository users) Find(nameOrNick string) ([]models.User, error) {
	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick) //%nameOrNick%

	lines, err := repository.db.Query(
		"SELECT id, name, nick, email, createdAt FROM users WHERE name LIKE ? or nick LIKE ?",
		nameOrNick, nameOrNick,
	)

	if err != nil {
		return nil, err
	}

	defer lines.Close()

	var users []models.User

	for lines.Next() {
		var user models.User

		if err = lines.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

//FindByID return a user from database
func (repository users) FindByID(ID uint64) (models.User, error) {
	line, err := repository.db.Query(
		"SELECT id, name, nick, email, createAt FROM users WHERE id = ?",
		ID,
	)

	if err != nil {
		return models.User{}, err
	}
	defer line.Close()

	var user models.User

	if line.Next() {
		if err = line.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

//Update edit user in database
func (repository users) Update(ID uint64, user models.User) error {
	statement, err := repository.db.Prepare(
		"UPDATE users SET name = ?, nick = ?, email = ? WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(user.Name, user.Nick, user.Email, ID); err != nil {
		return err
	}

	return nil
}

//Delete remove user from database
func (repository users) Delete(ID uint64) error {
	statement, err := repository.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(ID); err != nil {
		return err
	}

	return nil
}

//FindByEmail find user by email and return user id and user password hash
func (repository users) FindByEmail(email string) (models.User, error) {
	line, err := repository.db.Query("SELECT id, password FROM users WHERE email = ?", email)
	if err != nil {
		return models.User{}, err
	}
	defer line.Close()

	var user models.User

	if line.Next() {
		if err = line.Scan(&user.ID, &user.Password); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

// Follower add follower user id
func (repository users) Follower(userId, followerId uint64) error {
	statement, err := repository.db.Prepare(
		"INSERT INTO followers (user_id, follower_id) values (?, ?)",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(userId, followerId); err != nil {
		return err
	}

	return nil
}

func (repository users) Unfollow(userId, followerId uint64) error {
	statement, err := repository.db.Prepare(
		"DELETE FROM followers WHERE user_id = ? AND follower_id = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(userId, followerId); err != nil {
		return err
	}

	return nil
}

// FindFollowersByUserId find all follow from user
func (repository users) FindFollowersByUserId(userId uint64) ([]models.User, error) {
	lines, err := repository.db.Query(`
		SELECT u.id, u.name, u.nick, u.email, u.createdAt
		FROM users u INNER JOIN followers f ON u.id = f.follower_id WHERE f.user_id = ?
	`, userId,
	)
	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var followers []models.User
	for lines.Next() {
		var follower models.User

		if err = lines.Scan(
			&follower.ID,
			&follower.Name,
			&follower.Nick,
			&follower.Email,
			&follower.CreatedAt,
		); err != nil {
			return nil, err
		}

		followers = append(followers, follower)
	}

	return followers, nil
}

// FindFollowingByUserId find all users that user is following
func (repository users) FindFollowingByUserId(userId uint64) ([]models.User, error) {
	lines, err := repository.db.Query(`
		SELECT u.id, u.name, u.nick, u.email, u.createdAt
		FROM users u INNER JOIN followers f ON u.id = f.user_id WHERE f.follower_id = ?
	`, userId,
	)
	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var users []models.User

	for lines.Next() {
		var user models.User

		if err = lines.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	return users, nil
}

// FindPasswordById find user password by user id
func (repository users) FindPasswordById(userId uint64) (string, error) {
	line, err := repository.db.Query("SELECT password FROM users where id = ?", userId)
	if err != nil {
		return "", err
	}
	defer line.Close()

	var user models.User

	if line.Next() {
		if err = line.Scan(&user.Password); err != nil {
			return "", err
		}
	}

	return user.Password, nil
}

//  UpdateUserPassword update user password by user id
func (repository users) UpdateUserPassword(userId uint64, password string) error {
	statement, err := repository.db.Prepare("UPDATE users SET password = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(password, userId); err != nil {
		return err
	}

	return nil
}
