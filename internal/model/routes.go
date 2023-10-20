package model

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Load routes in to the server
func (server *Server) LoadRoutes (database Database) {
	server.router.GET("/users", func(c *gin.Context) { getUser(c, database) })
	server.router.POST("/users/create", func (c *gin.Context) { createUser(c, database) })
	server.router.POST("/users/delete", func(c *gin.Context) { deleteUser(c, database) })
	server.router.POST("/users/verify", func(c *gin.Context) { verifyUser(c, database) })
	server.router.POST("/users/username", func(c *gin.Context) { setUsername(c, database) }) 
	server.router.POST("/users/password", func(c *gin.Context) { setPassword(c, database) }) 

	server.router.GET("/applications/users", func(c *gin.Context) { getUsers(c, database) })
	server.router.POST("/applications/create", func (c *gin.Context) { createApplication(c, database) })

	server.router.GET("/testing/email", func(c *gin.Context) { emailTesting(c, database) })
}

func emailTesting (ctx *gin.Context, database Database) {
	res, err := SendEmail("hhargreaves2006@gmail.com", "This is a test email", "Hi mom!")

	if err != nil {
		ctx.JSON(err.Status, gin.H{ "status": err.Status, "error": err.Message })
		return
	}

	ctx.JSON(http.StatusOK, gin.H{ "status": 200, "message": res })
}

// Creates a new application in the database
func createApplication (ctx *gin.Context, database Database) {
	var appReq CreateApplicationRequest

	if err := ctx.ShouldBindJSON(&appReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	if err := Validate(appReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	ch := make (chan *AppResult)

	go database.CreateApplication(ch, appReq.Name)

	result := <- ch
	 
	if result.Error !=  nil {
		ctx.JSON(result.Error.Status, gin.H{ "status": result.Error.Status, "error": result.Error.Message })
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{ "status": 201, "application": result.Application })

	emailContent := `
	Dear ` + appReq.Email + `,

	We're excited to inform you that an application has been created for you!

	Application Details:
	- Application Name: ` + result.Application.Name + `
	- Application ID: ` + result.Application.ID.String() + `
	- Created Date: ` + time.Now().String() + `

	If you did not initiate this application, please contact our support team immediately.

	Thank you for using my server for your authentication needs.

	Best regards,
	Azpect3120
	https://github.com/Azpect3120/AuthenticationServer
	`

	SendEmail(appReq.Email, "Application Creation Notification", emailContent)
}

// Create a new user in an application in the database
// User the ApplicationID as a key
func createUser (ctx *gin.Context, database Database) {
	var userReq CreateUserRequest

	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	if err := Validate(userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	strCh := make(chan *StringResult)

	go HashString(strCh, userReq.Password)

	strResult := <- strCh

	if strResult.Error != nil {
		ctx.JSON(strResult.Error.Status, gin.H{ "status": strResult.Error.Status, "error": strResult.Error.Message })
		return
	}

	ch := make(chan *UserResult)

	go database.CreateUser(ch, userReq.ApplicationID, userReq.Username, strResult.String)

	result := <- ch

	if result.Error != nil {
		ctx.JSON(result.Error.Status, gin.H{ "status": result.Error.Status, "error": result.Error.Message })
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{ "status": 201, "user": result.User })
}

// Verify a user in the database
// Require ApplicationID as a key
/*
	body: {
		ApplicationID: string(uuid),
		Username: string,
		Password: string
	}

	return: {
		ApplicationID: string(uuid) || nil
		ID: string(uuid) || nil
	}
*/
func verifyUser (ctx *gin.Context, database Database) {
	var verifyReq VerifyUserRequest

	if err := ctx.ShouldBindJSON(&verifyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	if err := Validate(verifyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	ch := make(chan *UserResult)
	
	go database.VerifyUser(ch, verifyReq.ApplicationID, verifyReq.Username, verifyReq.Password)

	result := <- ch

	if result.Error != nil {
		ctx.JSON(result.Error.Status, gin.H{ "status": result.Error.Status, "error": result.Error.Message })
	} else {
		ctx.JSON(http.StatusOK, gin.H{ "status": 200, "user": result.User })
	}
}


// Get a user in the database
// Requires the UserID and the ApplicationID
// @param: app-id
// @param: user-id
func getUser (ctx *gin.Context, database Database) {
	var applicationID string = ctx.DefaultQuery("app-id", "")
	var userID string = ctx.DefaultQuery("user-id", "")
	
	ch := make(chan *UserResult)

	go database.GetUser(ch, applicationID, userID)

	result := <- ch

	if result.Error != nil {
		ctx.JSON(result.Error.Status, gin.H{ "status": result.Error.Status, "error": result.Error.Message })
	} else {
		ctx.JSON(http.StatusOK, gin.H{ "status": 200, "user": result.User })
	}
}

// Get all users stored in an application
// Requires the ApplicationID
// @param: app-id
func getUsers (ctx *gin.Context, database Database) {
	applicationId := ctx.DefaultQuery("app-id", "")


	ch := make(chan *UsersResult)

	go database.GetUsers(ch, applicationId)

	result := <- ch

	if result.Error != nil {
		ctx.JSON(result.Error.Status, gin.H{ "status": result.Error.Status, "error": result.Error.Message })
	} else {
		ctx.JSON(http.StatusOK, gin.H{ "status": 200, "users": result.Users })
	}
}

// Update a users username
// Requires the ApplicationID and the UserID
func setUsername (ctx *gin.Context, database Database) {
	var setRequest SetUsernameRequest	

	if err := ctx.ShouldBindJSON(&setRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	if err := Validate(setRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	ch := make(chan *UserResult)

	go database.SetUsername(ch, setRequest.ApplicationID, setRequest.ID, setRequest.Username)
	
	result := <- ch

	if result.Error != nil {
		ctx.JSON(result.Error.Status, gin.H{ "status": result.Error.Status, "error": result.Error.Message })	
	} else {
		ctx.JSON(http.StatusCreated, gin.H{ "status": 201, "user": result.User  })
	}
}

// Update a users password
// Requires the ApplicationID and the UserID
func setPassword (ctx *gin.Context, database Database) {
	var setRequest SetPasswordRequest

	if err := ctx.ShouldBindJSON(&setRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	ch := make(chan *UserResult)

	go database.SetPassword(ch, setRequest.ApplicationID, setRequest.ID, setRequest.Password)

	result := <- ch

	if result.Error != nil {
		ctx.JSON(result.Error.Status, gin.H{ "status": result.Error.Status, "error": result.Error.Message })	
	} else {
		ctx.JSON(http.StatusCreated, gin.H{ "status": 201, "user": result.User  })
	}
}

// Deletes a user from the database
// Requires the ApplicationID and the UserID
func deleteUser (ctx *gin.Context, database Database) {
	var deleteRequest DeleteUserReqest

	if err := ctx.ShouldBindJSON(&deleteRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	if err := Validate(deleteRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	ch := make(chan *ErrorResult)

	go database.DeleteUser(ch, deleteRequest.ApplicationID, deleteRequest.ID)

	result := <- ch

	if result.Error != nil {
		ctx.JSON(result.Error.Status, gin.H{ "status": result.Error.Status, "error": result.Error.Message })	
	} else {
		ctx.JSON(http.StatusOK, gin.H{ "status": 200, "message": "User was deleted" })
	}

}
