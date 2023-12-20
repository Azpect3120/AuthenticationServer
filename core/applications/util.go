package applications

import (
	"fmt"
	"strings"
)

// Match an array of inputted data columns
// and updates inputted array to validated
// array of data columns which can be stored
// in the db. A message will be returned
// which tells the called which columns were
// invalid. Valid column names are below,
// case insensitive.
// Username, First, Last, Full, Email, Password, Data
// Does not allow for duplicate columns.
func MatchColumns (c *[]string) (msg string) {
    var valid []string = make([]string, 0, len(*c))
    for _, col := range *c {
        switch strings.TrimSpace(strings.ToLower(col)) {
            case "username", "first", "last", "full", "email", "password", "data":
                valid = append(valid, strings.TrimSpace(strings.ToLower(col)))
            case "first name": 
                valid = append(valid, "first")
            case "last name":
                valid = append(valid, "last")
            case "full name":
                valid = append(valid, "full")
            default:
                msg += fmt.Sprintf("'%s' is invalid. ", col)
        }
    }
    *c = valid
    msg += removeDuplicates(c)
    msg = strings.TrimSpace(msg)
    return
}

// Remove duplicate inputs from slice.
// Used to remove duplicate columns when
// they're provided in the 'MatchColumns'
// function.
func removeDuplicates (s *[]string) (msg string) {
    var contents map[string]struct{} = make(map[string]struct{})
    var valid []string = make([]string, 0, len(*s))
    for _, col := range *s {
        fmt.Println(col)
        if _, exists := contents[col]; !exists {
            contents[col] = struct{}{}
            valid = append(valid, col)
        } else {
            msg += fmt.Sprintf("Duplicate '%s' was removed. ", col)
        }
    }
    *s = valid
    return
}
