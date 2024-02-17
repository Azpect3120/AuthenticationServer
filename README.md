# Media Storage Server

## Routes

### Users
`GET` v2/applications/:id/users/:id -> Return a specific user object

`GET` v2/applications/:id/users -> Returns an applications users list

`POST` v2/applications/:appid/users -> Create a user

`PATCH` v2/applications/:appid/users/:userid -> Update PIECES of a user object (does not require all fields)

`PUT` v2/applications/:appid/users/:userid -> Update an ENTIRE user object (requires all fields)

`DELETE` v2/applications/:appid/users/:userid -> Deletes a user

### Applications
`GET` v2/applications -> Returns a list of active applications

`GET` v2/applications/:id -> Return application object

`POST` v2/applications -> Create an application

`PATCH` v2/applications/:id -> Update PIECE of an application object (does not require all fields)

`PUT` v2/applications/:id -> Updates an ENTIRE application object (requires all fields)

`DELETE` v2/applications/:id -> Deletes an application

## Models

### User

    - ID
    - App ID
    - Username
    - First Name
    - Last Name
    - Full Name
    - Email
    - Password
    - Data (stringified)
    - Created At
    - Last Updated At

    Notes:
        - Use type sql.NullString instead of string for values that may be null in the tables. Like the user data that may not be stored

### Applications

    - ID
    - Name
    - Created At
    - Last Updated At
    
    NOTES:
        - Pass an array of columns into the `POST` request
        - Conert []string to pq.StringArray using `pq.StringArray()` method


### Testing requests

```sh
    curl -X POST -H "Content-Type: application/json" -d '{"name": "Hello World", "columns": ["data", "email", "HELLO WORLD"]}' localhost:3000/v2/applications | json
    curl -X POST -H "Content-Type: application/json" -d '{"name": "Hello World", "columns": ["data", "email", "password", "username", "username", "first", "first name", "last name", "hello world", "hi mom :)"]}' localhost:3000/v2/applications | json
```
