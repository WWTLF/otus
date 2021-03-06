package UserHandlers

import (
	Core "user_list/CORE"
	"user_list/models"
	"user_list/restapi/operations/user"

	"github.com/go-openapi/runtime/middleware"
)

func CreateUser(params user.CreateUserParams) middleware.Responder {
	var id int32
	query := Core.GetInstance().DB.QueryRow(
		"INSERT INTO userlist(username, firstname, lastname, email, phone) VALUES ($1, $2, $3, $4, $5 ) RETURNING id;",
		params.Body.Username,
		params.Body.FirstName,
		params.Body.LastName,
		params.Body.Email,
		params.Body.Phone)

	err := query.Scan(&id)
	if err != nil {
		Core.HandelError(err, false)
		var code int32 = 500
		message := err.Error()
		return user.NewCreateUserDefault(int(code)).WithPayload(&models.Error{Code: &code, Message: &message})
	}
	createdUser := models.User{
		Email:     params.Body.Email,
		FirstName: params.Body.FirstName,
		ID:        int64(id),
		LastName:  params.Body.LastName,
		Phone:     params.Body.Phone,
		Username:  params.Body.Username,
	}

	return user.NewCreateUserOK().WithPayload(&createdUser)
}

func DeleteUser(params user.DeleteUserParams) middleware.Responder {
	row := Core.GetInstance().DB.QueryRow("DELETE FROM userlist WHERE id = $1 RETURNING id", params.UserID)
	var id int32
	err := row.Scan(&id)

	if err != nil {
		Core.HandelError(err, false)
		var code int32 = 404
		message := "user not found"
		return user.NewFindUserByIDDefault(int(code)).WithPayload(&models.Error{Code: &code, Message: &message})
	}
	return user.NewDeleteUserDefault(user.DeleteUserNoContentCode)
}

func FindUserById(params user.FindUserByIDParams) middleware.Responder {

	row := Core.GetInstance().DB.QueryRow("SELECT id, username, firstname, lastname, email, phone FROM userlist WHERE id = $1;", params.UserID)
	responseUser := models.User{}
	err := row.Scan(
		&responseUser.ID,
		&responseUser.Username,
		&responseUser.FirstName,
		&responseUser.LastName,
		&responseUser.Email,
		&responseUser.Phone,
	)

	if err != nil {
		Core.HandelError(err, false)
		var code int32 = 404
		message := "user not found"
		return user.NewFindUserByIDDefault(404).WithPayload(&models.Error{Code: &code, Message: &message})
	}

	return user.NewCreateUserOK().WithPayload(&responseUser)
}

func UpdateUser(params user.UpdateUserParams) middleware.Responder {
	row := Core.GetInstance().DB.QueryRow("UPDATE userlist SET"+
		" username = $1,"+
		" firstname = $2,"+
		" lastname = $3,"+
		" email = $4,"+
		" phone = $5"+
		" WHERE id = $6 RETURNING id;",
		params.Body.Username,
		params.Body.FirstName,
		params.Body.LastName,
		params.Body.Email,
		params.Body.Phone,
		params.UserID)
	var id int32
	err := row.Scan(&id)

	if err != nil {
		Core.HandelError(err, false)
		var code int32 = 404
		message := err.Error()
		return user.NewUpdateUserDefault(int(code)).WithPayload(&models.Error{Code: &code, Message: &message})
	}

	responseUser := models.User{
		Email:     params.Body.Email,
		FirstName: params.Body.FirstName,
		ID:        params.UserID,
		LastName:  params.Body.LastName,
		Phone:     params.Body.Phone,
		Username:  params.Body.Username,
	}
	return user.NewUpdateUserOK().WithPayload(&responseUser)
}
