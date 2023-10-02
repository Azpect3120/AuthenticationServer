# Go Authentication Server

A simple Go authentication server built with the Gin framework, PostgreSQL database, and UUID-keys for user authentication.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Create an Application](#create-an-application)
  - [Create a User](#create-a-user)
  - [Authenticate a User](#authenticate-a-user)
  - [Update a User](#update-a-user)
  - [Delete a User](#delete-a-user)
- [Contributing](#contributing)
- [License](#license)

## Introduction

This project provides a basic template for building a Go authentication server. It includes routes for user login, registration, and profile retrieval, all secured with UUID keys/ids and secrets keys. You can use this template as a starting point to create your own authentication system for web applications. Or you can connect to a live server running this application. Instructions for using the live server can be found below.

## Features

- User authentication using UUID keys.
- User registration with secure password storage.
- Retrieval of users.
- Multi-application storage.
- Basic error handling and logging.

## Getting Started

Follow these steps to get the project up and running on your local machine.

### Prerequisites

- Go (1.16 or higher)
- PostgreSQL database
- Git (optional)

### Installation

1. Clone the repository (if you haven't already):

```bash
    git clone https://github.com/yourusername/go-authentication-server.git
``` 

2. Set up your PostgreSQL database and configure the connection string in the `cmd/authServer/main.go` file:
```go
  // Replace with your PostgreSQL connection string.
  connectionString := "postgresql://username:password@localhost/dbname"
```

3. Install dependencies:
```bash
  go mod tidy
```

4. Build and run the server:
```bash
  go run cmd/AuthenticationServer/main.go
```
Your authentication server should now be running on `http://localhost:8080`.

## Usage

### <a id="create-an-application"></a> Create an Application

All users are stored within an `application`, which means for new users you must first create an application. To do this you can send a post request to the `/createApplication` endpoint.

```json
  {
    "name": "your-app-name"
  }
```

Example response:

NOTE: The `ID` returned by this request is very important! It will be REQUIRED to request any future data from the server (to ensure safety and security). So do not lose it!

```json
  {
    "status": 201,
    "name": "your-app-name",
    "ID": "00000000-0000-0000-00000000"
  }
```

### <a id="create-a-user"></a> Create a User

Once you have created an `application`, you are now set to begin creating and storing users!
To create a user in your new application you can send a post request to the `/createUser` endpoint.

The `ApplicationID` is the `ID` you were returned when you created an application.
Passwords are stored securely in the database using various hashing methods, so there is no need to hash or encrypt the passwords on the front-end. Though if you are concerned, you may still consider using your own hashing or encryption on the front-end BEFORE sending data to the server.

```json
  {
    "ApplicationID": "application-id-here",
    "Username": "your-username-here",
    "Password": "your-password-here"
  }
```

Example response:

NOTE: The `user.ID` and the `user.applicationID` should be saved for future use when updating the users information.

```json
  {
    "status": 201,
    "user": {
      "ID": "00000000-0000-0000-00000000",
      "username": "your-username-here",
      "password": "your-hashed-password-here",
      "applicationID": "00000000-0000-0000-00000000"
    }
  }
```

### <a id="authenticate-a-user"></a> Authenticate a User

You have now created your first user in an application! The user will be stored securely in the database and can now be used to authenticate logins. To do this, you can send a post request to the `/verifyUser` endpoint. The password that you pass will be sent in plain text and will be compared to the hashed password stored in the database. It will not be stored in the database along the way to ensure security. If you are building your own front-end security, you will need to match the data sent to the server.

```json
  {
    "ApplicationID": "application-id-here",
    "Username": "your-username-here",
    "Password": "your-password-here"
  }
```

Example responses:

NOTE: The password that is returned will still be hashed, to ensure no data leakage will occur.

User was verified successfully

```json
  {
    "status": 200,
    "user": {
      "ID": "00000000-0000-0000-00000000",
      "username": "your-username-here",
      "password": "your-hashed-password-here",
      "applicationID": "00000000-0000-0000-00000000"
    }
  }
```

User was not verified

```json
  {
    "status": 400,
    "error": "User was not verified"
  }
```

### <a id="update-a-user"></a> Update a User

Once you have created users you can update their username and password using the respective post endpoints, `/setUsername` and `/setPassword`.

NOTE: No validation is done on the server side, so any password validation should be handled on the front-end

`/setUsername`
```json
  {
    "ApplicationID": "00000000-0000-0000-00000000",
    "ID": "00000000-0000-0000-00000000",
    "Username": "new-username-here"
  }
```

Example responses: 

Username was successfully updated
```json
  {

  }
```

Username was not successfully updated
```json
  {

  }
```

`/setPassword`
```json
  {
    "ApplicationID": "00000000-0000-0000-00000000",
    "ID": "00000000-0000-0000-00000000",
    "Password": "new-password-here"
  }
```

Example responses: 

Password was successfully updated
```json
  {

  }
```

Password was not successfully updated
```json
  {

  }
```

### <a id="delete-a-user"></a> Delete a User

Finally, you can delete a user from an application by sending a post request to the `/deleteUser` endpoint.

```json
  {
    "ApplicationID": "00000000-0000-0000-00000000",
    "ID": "00000000-0000-0000-00000000",
  }
```

Example responses: 

User was successfully deleted
```json
  {

  }
```

User was not successfully deleted
```json
  {

  }
```



## Contributing

Contributions are welcome! If you'd like to contribute to this project, please follow these steps:

1. Fork the project.
2. Create a new branch for your feature or bug fix.
3. Make your changes.
4. Test your changes thoroughly.
5. Create a pull request.

## License

The project is licensed under the **MIT License**
