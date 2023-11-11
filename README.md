# User Authentication Server

A simple user authentication server built with the Gin framework, PostgreSQL database, and UUID-keys for user authentication.

## Table of Contents

-   [Introduction](#introduction)
-   [Features](#features)
-   [Getting Started](#getting-started)
    -   [Prerequisites](#prerequisites)
    -   [Installation](#installation)
    -   [Database](#database)
-   [Usage](#usage)
    -   [Create an Application](#create-an-application)
    -   [Create a User](#create-a-user)
    -   [Authenticate a User](#authenticate-a-user)
    -   [Update a User](#update-a-user)
    -   [Delete a User](#delete-a-user)
    -   [Get a User](#get-a-user)
-   [Contributing](#contributing)
-   [License](#license)

## Introduction

This project provides a basic template for building a Go authentication server. It includes routes for user login, registration, and profile retrieval, all secured with UUID keys/ids and secrets keys. You can use this template as a starting point to create your own authentication system for web applications.

## Features

-   User authentication using UUID keys.
-   User registration with secure password storage.
-   Retrieval of users.
-   Multi-application storage.
-   Basic error handling and logging.

## Getting Started

Follow these steps to get the project up and running on your local machine.

### Prerequisites

-   Go (1.16 or higher)
-   PostgreSQL database
-   Git (optional)

### Installation

1. Clone the repository (if you haven't already):

```bash
    git clone https://github.com/Azpect3120/AuthenticationServer.git
```

2. Set up your PostgreSQL database and configure the connection details in the `.env` file:

```.env
  # This url can found in the dashboard of most PSQL hosts or can be constructed using the required pieces
  db_url=your-connection-url-here

  # For more information visit this link: https://www.hostinger.com/tutorials/how-to-use-free-google-smtp-server
  smtp_email=your-email-here
  smtp_password=your-email-app-password-here
```

3. Install dependencies:

```bash
  go mod tidy
```

4. Build and run the server:

```bash
  go build -o ./bin/server cmd/authServer/main.go
  ./bin/server
```

Your authentication server should now be running on `http://localhost:8080`.

### Database

Once the server is up and running you will need to connect to a PostgreSQL database. If you would like the code to work out of the box, you may copy the database schema provided below.

```sql
  CREATE TABLE IF NOT EXISTS applications (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      name TEXT
  );

  CREATE TABLE IF NOT EXISTS users (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      applicationID UUID,
      username TEXT,
      password TEXT,
      data TEXT,
      FOREIGN KEY (applicationID) REFERENCES applications(ID)
  );

  CREATE TABLE IF NOT EXISTS logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    elapsedTime TEXT,
    method VARCHAR(16),
    endpoint TEXT,
    query TEXT,
    reqBody TEXT,
    code INT,
    resBody TEXT
  );
```

## Usage

### <a id="create-an-application"></a> Create an Application

All users are stored within an `application`, which means for new users you must first create an application. To do this you can send a post request to the `/applications/create` endpoint.

The data property can be used to store raw JSON data. It should be stringified before being sent to the server otherwise the server will not know how to handle the request. If you do not wish to hold any data along with the users you can send an empty string;

```json
{
    "name": "your-app-name",
    "email": "your-email-here"
}
```

Example response:

NOTE: The `ID` returned by this request is very important! It will be REQUIRED to request any future data from the server (to ensure safety and security). So do not lose it!

```json
{
    "status": 201,
    "application": {
        "name": "your-app-name",
        "ID": "00000000-0000-0000-00000000"
    }
}
```

### <a id="create-a-user"></a> Create a User

Once you have created an `application`, you are now set to begin creating and storing users!
To create a user in your new application you can send a post request to the `/users/create` endpoint.

The `ApplicationID` is the `ID` you were returned when you created an application.
Passwords are stored securely in the database using various hashing methods, so there is no need to hash or encrypt the passwords on the front-end. Though if you are concerned, you may still consider using your own hashing or encryption on the front-end BEFORE sending data to the server.

NOTE: `data` field must not be a blank string! You may provide a blank object `{}` but you cannot provide a blank string `""`

```json
{
    "applicationID": "application-id-here",
    "username": "your-username-here",
    "password": "your-password-here",
    "data": "{'data': 'your-json-data-here'}"
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
        "applicationID": "00000000-0000-0000-00000000",
        "data": "{'data': 'your-json-data-here'}"
    }
}
```

### <a id="authenticate-a-user"></a> Authenticate a User

You have now created your first user in an application! The user will be stored securely in the database and can now be used to authenticate logins. To do this, you can send a post request to the `/users/verify` endpoint. The password that you pass will be sent in plain text and will be compared to the hashed password stored in the database. It will not be stored in the database along the way to ensure security. If you are building your own front-end security, you will need to match the data sent to the server.

```json
{
    "applicationID": "application-id-here",
    "username": "your-username-here",
    "password": "your-password-here"
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
        "applicationID": "00000000-0000-0000-00000000",
        "data": "{'data': 'your-json-data-here'}"
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

Once you have created users you can update their username, password, and data using the respective post endpoints, `/users/username`, `/users/password`, and `/users/data`.

NOTE: No validation is done on the server side, so any password validation should be handled on the front-end

`/users/username`

```json
{
    "applicationID": "00000000-0000-0000-00000000",
    "ID": "00000000-0000-0000-00000000",
    "username": "new-username-here"
}
```

Example responses:

Username was successfully updated

```json
{
    "status": 201,
    "user": {
        "ID": "00000000-0000-0000-00000000",
        "username": "your-new-username-here",
        "password": "your-hashed-password-here",
        "applicationID": "00000000-0000-0000-00000000",
        "data": "{'data': 'your-json-data-here'}"
    }
}
```

Username was not successfully updated

```json
{
    "status": 400,
    "error": "The users username could not be change."
}
```

`/users/password`

```json
{
    "applicationID": "00000000-0000-0000-00000000",
    "ID": "00000000-0000-0000-00000000",
    "password": "new-password-here"
}
```

Example responses:

Password was successfully updated

```json
{
    "status": 201,
    "user": {
        "ID": "00000000-0000-0000-00000000",
        "username": "your-username-here",
        "password": "your-new-hashed-password-here",
        "applicationID": "00000000-0000-0000-00000000",
        "data": "{'data': 'your-json-data-here'}"
    }
}
```

Password was not successfully updated

```json
{
    "status": 400,
    "error": "The users password could not be changed."
}
```

`/users/data`

```json
{
    "applicationID": "00000000-0000-0000-00000000",
    "ID": "00000000-0000-0000-00000000",
    "data": "{'data': 'your-new-json-data-here'}"
}
```

Example responses:

Data was successfully updated

```json
{
    "status": 201,
    "user": {
        "ID": "00000000-0000-0000-00000000",
        "username": "your-new-username-here",
        "password": "your-hashed-password-here",
        "applicationID": "00000000-0000-0000-00000000",
        "data": "{'data': 'your-json-data-here'}"
    }
}
```

Data was not successfully updated

```json
{
    "status": 400,
    "error": "The users data could not be change."
}
```

### <a id="delete-a-user"></a> Delete a User

Finally, you can delete a user from an application by sending a post request to the `/users/delete` endpoint.

```json
{
    "applicationID": "00000000-0000-0000-00000000",
    "ID": "00000000-0000-0000-00000000"
}
```

Example responses:

User was successfully deleted

```json
{
    "status": 200,
    "message": "The user was deleted"
}
```

User was not successfully deleted

```json
{
    "status": 404,
    "message": "The user was not found"
}
```

### <a id="get-a-user"></a> Get a User

Now that you know how to create, update, delete, and verify users, you should know how to get the users. There are two ways, both which require the `ApplicationID`. Sending a get request to the `/users` endpoint will return a single user and sending a get request to the `/applications/users` endpoint will return a list of all the users in application.

NOTE: It is best to not allow users to send requests to the `/application/users` endpoint to prevent data leaks.

`/users`

Params:

-   `app-id` : The ID of the application the user was added into
-   `user-id` : The ID of the user you wish to delete

```bash
  /getUser?app-id=0000000-0000-0000-00000000&user-id=0000000-0000-0000-00000000
```

Example responses:

User was found

```json
{
    "status": 200,
    "user": {
        "ID": "00000000-0000-0000-00000000",
        "username": "a-username",
        "password": "a-hashed-password",
        "applicationID": "00000000-0000-0000-00000000",
        "data": "{'data': 'your-json-data-here'}"
    }
}
```

No user was found

```json
{
    "status": 404,
    "error": "User was not found."
}
```

`/application/users`

Params:

-   `app-id` : The ID of the application you want to view the users of

```bash
  /getUsers?app-id=0000000-0000-0000-00000000
```

Example responses:

Provided `applicationID` exists

```json
{
    "status": 200,
    "users": [
        {
            "ID": "00000000-0000-0000-00000000",
            "username": "a-username-1",
            "password": "a-hashed-password-1",
            "applicationID": "00000000-0000-0000-00000000",
            "data": "{'data':  'your-json-data-here'}"
        },
        {
            "ID": "00000000-0000-0000-00000000",
            "username": "a-username-2",
            "password": "a-hashed-password-2",
            "applicationID": "00000000-0000-0000-00000000",
            "data": "{'data':  'your-json-data-here'}"
        }
    ]
}
```

Provided `applicationID` does not exist

```json
{
    "status": 404,
    "error": "Application with the provided ID does not exist."
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

# The project is licensed under username the **MIT License**
