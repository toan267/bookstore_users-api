package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"toan267/bookstore_users-api/domain/users"
	"toan267/bookstore_users-api/services"
	"toan267/bookstore_users-api/utils/errors"
)

func CreateUser(c *gin.Context) {
	var user users.User
	//fmt.Println(user)
	//bytes, err := ioutil.ReadAll(c.Request.Body)
	//if err != nil {
	//	//TODO: handle error
	//	return
	//}
	//if err := json.Unmarshal(bytes, &user); err != nil {
	//	fmt.Println(err.Error())
	//	//TODO: Handle json error
	//	return
	//}
	if err := c.ShouldBindJSON(&user); err != nil {
		//restErr := errors.RestErr{
		//	Message: "Invalid json body",
		//	Status: http.StatusBadRequest,
		//	Error: "Bad_request",
		//}
		restErr := errors.NewBadRequestError("Invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	//c.String(http.StatusOK, "Not implement yet")
	c.JSON(http.StatusCreated, result)
}

func SearchUser(c *gin.Context) {
	c.String(http.StatusOK, "Not implement yet")
}

func GetUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("user id should be a number")
		c.JSON(err.Status, err)
		return
	}
	user, getErr := services.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, user)
}