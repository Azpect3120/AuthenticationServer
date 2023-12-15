package main

import (
	"fmt"
	"strings"

	"github.com/Azpect3120/AuthenticationServer/core/model"
	s "github.com/Azpect3120/AuthenticationServer/core/server"
	"github.com/gin-gonic/gin"
)

func main() {
	server := s.NewServer(3000, "")

    
    // `POST` v2/applications -> Create an application
    s.AddRoute(server, "post", "/v2/applications", func(ctx *gin.Context) {
        var req model.CreateApplicationRequest
        if err := ctx.ShouldBindJSON(&req); err != nil {
            ctx.JSON(404, gin.H{ "error": err.Error() })
        }

        message := MatchColumns(&req.Columns)


        ctx.JSON(200, gin.H{ "status": 200, "message": message, "data": req })
    })
    

    s.Listen(server)
}


// Match an array of inputted data columns 
// and updates inputted array to validated 
// array of data columns which can be stored 
// in the db. A message will be returned 
// which tells the called which columns were 
// invalid. Valid column names are below,
// case insensitive.
// Username, First, Last, Full, Email, Password, Data
func MatchColumns (c *[]string) (msg string) {
    var valid []string = make([]string, 0, len(*c))
    for _, col := range *c {
        switch strings.ToLower(col) {
            case "username", "first", "last", "full", "email", "password", "data":
                valid = append(valid, strings.ToLower(col))
            default:
                msg += fmt.Sprintf("'%s' is invalid. ", col)
        }
    }
    *c = valid
    return msg
}
