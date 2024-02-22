package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/Azpect3120/AuthenticationServer/core/applications"
	"github.com/Azpect3120/AuthenticationServer/core/database"
	"github.com/Azpect3120/AuthenticationServer/core/model"
	s "github.com/Azpect3120/AuthenticationServer/core/server"
	"github.com/Azpect3120/AuthenticationServer/core/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err.Error())
	}

	port := os.Getenv("AUTH_SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	server := s.NewServer(port)
	db := database.NewDatabase(os.Getenv("DB_URL"))

	// `GET` v2/applications/:id -> Gets an application
	s.AddRoute(server, "get", "/v2/applications/:id", func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
			return
		}
		app, code, err := applications.Retrieve(db, id)
		if err != nil {
			ctx.JSON(code, gin.H{"status": code, "error": err.Error()})
			return
		}
		ctx.JSON(code, gin.H{"status": code, "application": app})
	})

	// `POST` v2/applications -> Create an application
	s.AddRoute(server, "post", "/v2/applications", func(ctx *gin.Context) {
		var req model.CreateApplicationRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(404, gin.H{"status": 404, "error": err.Error()})
			return
		}

		message := applications.MatchColumns(&req.Columns)
		app := applications.New(req.Name, req.Columns)

		code, err := applications.Insert(db, app)
		if err != nil {
			ctx.JSON(code, gin.H{"status": code, "error": err.Error()})
			return
		}

		ctx.JSON(code, gin.H{"status": code, "message": message, "application": app})
	})

	// `PATCH` v2/applications/:id -> Updates an application.
	// Only provided fields will be updated.
	s.AddRoute(server, "patch", "/v2/applications/:id", func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
			return
		}

		var req model.ModifyApplicationRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(404, gin.H{"status": 404, "error": err.Error()})
			return
		}

		app, message, code, err := applications.Update(db, id, req.Name, req.Columns)
		if err != nil {
			ctx.JSON(code, gin.H{"status": code, "error": err.Error()})
			return
		}

		ctx.JSON(code, gin.H{"status": code, "message": message, "application": app})
	})

	// `PUT` v2/applications/:id -> Updates/Overwrites an application. Requires all fields.
	s.AddRoute(server, "put", "/v2/applications/:id", func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
			return
		}

		var req model.ModifyApplicationRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(404, gin.H{"status": 404, "error": err.Error()})
			return
		}

		app, message, code, err := applications.Overwrite(db, id, req.Name, req.Columns)
		if err != nil {
			ctx.JSON(code, gin.H{"status": code, "error": err.Error()})
			return
		}

		ctx.JSON(code, gin.H{"status": code, "message": message, "application": app})

	})

	// `DELETE` v2/applications/id -> Delete an application
	s.AddRoute(server, "delete", "/v2/applications/:id", func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
			return
		}

		code, err := applications.Delete(db, id)
		if err != nil {
			ctx.JSON(code, gin.H{"status": code, "error": err.Error()})
			return
		}

		ctx.JSON(204, gin.H{"status": 204, "message": "Application was deleted."})
	})

	// `GET` v2/applications -> Get all applications
	s.AddRoute(server, "get", "/v2/applications", func(ctx *gin.Context) {
		apps, code, err := applications.RetrieveAll(db)
		if err != nil {
			ctx.JSON(code, gin.H{"status": code, "error": err.Error()})
			return
		}
		ctx.JSON(code, gin.H{"status": code, "applications": apps, "count": len(apps)})
	})

	// `GET` v2/applications/:id/users/:id -> Get a user from an application
	s.AddRoute(server, "get", "/v2/applications/:id/users/:uid", func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
			return
		}

		uid, err := uuid.Parse(ctx.Param("uid"))
		if err != nil {
			ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
			return
		}

		data, code, err := users.Retrieve(db, id, uid)
		if err != nil {
			ctx.JSON(code, gin.H{"status": code, "error": err.Error()})
			return
		}

		ctx.JSON(code, gin.H{"status": code, "user": data})
	})

	// `GET` v2/applications/:id/users -> Get all users for an application
	s.AddRoute(server, "get", "/v2/applications/:id/users", func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
			return
		}

		data, code, err := users.RetrieveAll(db, id)
		if err != nil {
			ctx.JSON(code, gin.H{"status": code, "error": err.Error()})
			return
		}

		ctx.JSON(code, gin.H{"status": code, "users": data, "count": len(data)})

	})

	// `POST` v2/applications/:id/users -> Create a user for an applications
	s.AddRoute(server, "post", "/v2/applications/:id/users", func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
			return
		}

		var req model.UserData
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(404, gin.H{"status": 404, "error": err.Error()})
			return
		}

		user := users.New(id, &req)
		message, code, err := users.Validate(db, id, user)
		if err != nil {
			ctx.JSON(code, gin.H{"status": code, "message": message, "error": err.Error()})
			return
		}

		code, err = users.Insert(db, user)
		if err != nil {
			ctx.JSON(code, gin.H{"status": code, "error": err.Error()})
			return
		}

		appColumns, err := users.GetApplicationColumns(db, user.ApplicationID)
		if err != nil {
			ctx.JSON(500, gin.H{"status": 500, "error": err.Error()})
			return
		}

		providedColumns := make(map[string]string)

		for _, col := range appColumns {
			val := reflect.ValueOf(*user).FieldByName(users.COLUMNS[col])
			var fieldValue string
			switch val.Interface().(type) {
			case uuid.UUID:
				fieldValue = val.Interface().(uuid.UUID).String()
			case time.Time:
				fieldValue = val.Interface().(time.Time).String()
			default:
				fieldValue = val.String()
			}
			providedColumns[col] = fieldValue
		}

		ctx.JSON(code, gin.H{"status": code, "message": "User was created.", "user": providedColumns})
	})

	// `PATCH` v2/applications/:id/users/:uid -> Updates a user in an application. Does not require all fields
	s.AddRoute(server, "patch", "/v2/applications/:id/users/:uid", func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
			return
		}

		uid, err := uuid.Parse(ctx.Param("uid"))
		if err != nil {
			ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
			return
		}

		var req model.UserData
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(404, gin.H{"status": 404, "error": err.Error()})
			return
		}

		user, code, err := users.Update(db, id, uid, &req)
		if err != nil {
			ctx.JSON(code, gin.H{"status": code, "error": err.Error()})
			return
		}

		appColumns, err := users.GetApplicationColumns(db, user.ApplicationID)
		if err != nil {
			ctx.JSON(500, gin.H{"status": 500, "error": err.Error()})
			return
		}

		providedColumns := make(map[string]string)

		for _, col := range appColumns {
			val := reflect.ValueOf(*user).FieldByName(users.COLUMNS[col])
			var fieldValue string
			switch val.Interface().(type) {
			case uuid.UUID:
				fieldValue = val.Interface().(uuid.UUID).String()
			case time.Time:
				fieldValue = val.Interface().(time.Time).String()
			default:
				fieldValue = val.String()
			}
			providedColumns[col] = fieldValue
		}

		ctx.JSON(code, gin.H{"status": code, "message": "User credentials were successfully validated", "user": providedColumns})
	})

	// `DELETE` v2/applications/:id/users/:uid -> Delete a user from an application
	s.AddRoute(server, "delete", "/v2/applications/:id/users/:uid", func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
			return
		}

		uid, err := uuid.Parse(ctx.Param("uid"))
		if err != nil {
			ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
			return
		}

		code, err := users.Delete(db, id, uid)
		if err != nil {
			ctx.JSON(code, gin.H{"status": code, "error": err.Error()})
			return
		}

		ctx.JSON(204, gin.H{"status": 204, "message": "User was deleted."})
	})

	// `POST` v2/applications/:id/validation -> Validate a user's login credentials
	s.AddRoute(server, "post", "/v2/applications/:id/validate", func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, gin.H{"status": 400, "error": err.Error()})
			return
		}

		var req model.ValidateUserRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(404, gin.H{"status": 404, "error": err.Error()})
			return
		}

		user, message, code, err := users.ValidateLogin(db, id, req.Columns, &req.User)
		if err != nil {
			ctx.JSON(code, gin.H{"status": code, "message": message, "error": err.Error()})
			return
		}

		if message != "" {
			ctx.JSON(code, gin.H{"status": code, "message": fmt.Sprintf("User credentials were successfully validated. However, %s", message), "user": user})
			return
		}

		appColumns, err := users.GetApplicationColumns(db, user.ApplicationID)
		if err != nil {
			ctx.JSON(500, gin.H{"status": 500, "error": err.Error()})
			return
		}

		providedColumns := make(map[string]string)

		for _, col := range appColumns {
			val := reflect.ValueOf(*user).FieldByName(users.COLUMNS[col])
			var fieldValue string
			switch val.Interface().(type) {
			case uuid.UUID:
				fieldValue = val.Interface().(uuid.UUID).String()
			case time.Time:
				fieldValue = val.Interface().(time.Time).String()
			default:
				fieldValue = val.String()
			}
			providedColumns[col] = fieldValue
		}

		ctx.JSON(code, gin.H{"status": code, "message": "User credentials were successfully validated", "user": providedColumns})
	})

	s.Listen(server)
}
