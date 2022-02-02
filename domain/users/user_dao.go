package users

import (
	"fmt"
	"toan267/bookstore_users-api/datasources/mysql/users_db"
	"toan267/bookstore_users-api/logger"
	"toan267/bookstore_users-api/utils/errors"
	"toan267/bookstore_users-api/utils/mysql_utils"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created, password, status) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser = "SELECT id, first_name, last_name, email, date_created, status FROM users where id=?;"
	queryUpdateUser = "UPDATE users set first_name=?, last_name=?, email=? WHERE id=?"
	queryDeleteUser = "DELETE FROM users where id = ?"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when try to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		logger.Error("error when try to get user by id", err)
		return errors.NewInternalServerError("database error")
		//return mysql_utils.ParseError(err)
/*		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to get user %d: %s", user.Id, err.Error()))
*/	}
/*	result := usersDB[user.Id]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
	}
	user.Id = result.Id
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated
*/	return nil
}
func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when try to prepare save user statement", err)
		return errors.NewInternalServerError("database error")
		//return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	//user.DateCreated = date_utils.GetNowString()
	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Password, user.Status)
	if saveErr != nil {
		logger.Error("error when try to save user", saveErr)
		return errors.NewInternalServerError("database error")
		//return mysql_utils.ParseError(saveErr)
/*		sqlErr, ok := saveErr.(*mysql.MySQLError)
		if ok {
			//mysql errors
			switch sqlErr.Number {
			case 1062:
				return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", user.Email))
			}
			//if strings.Contains(err.Error(), error_unique) {
			//				return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", user.Email))
			//}
			return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", saveErr.Error()))
		}
		//not mysql error
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", saveErr.Error()))
*/	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when try to get last insert id after creating a new user", saveErr)
		return errors.NewInternalServerError("database error")
		//return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}
	user.Id = userId
	/*current:=usersDB[user.Id]
	if current != nil {
		if current.Email ==user.Email {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already registered", user.Email))
		}
		return errors.NewBadRequestError(fmt.Sprintf("user %d already exists", user.Id))
	}
	user.DateCreated = date_utils.GetNowString()
	usersDB[user.Id] = user*/
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when try to prepare update statement", err)
		return errors.NewInternalServerError("database error")
		//return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when try to update user", err)
		return errors.NewInternalServerError("database error")
		//return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when try to prepare delete statement", err)
		return errors.NewInternalServerError("database error")
		//return mysql_utils.ParseError(err)
	}
	defer stmt.Close()
	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when try to delete user", err)
		return errors.NewInternalServerError("database error")
		//return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	rows, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()
	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName,
			&user.Email, &user.DateCreated, &user.Status); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}
