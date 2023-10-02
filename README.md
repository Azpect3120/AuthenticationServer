# Go Authentication Server

A simple Go authentication server built with the Gin framework, PostgreSQL database, and JWT (JSON Web Tokens) for user authentication.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - WIPPPP
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

2. Set up your PostgreSQL database and configure the connection string in the `main.go` file:
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

All users are stored within an `application`, which means for new users you must first create an application. To do this you can send a post request to the `/createApplication` endpoint.

```json
  {
    "name": "your-app-name"
  }
```
Example response:

```json
  {
    "status": 201,
    "name": "your-app-name",
    "ID:" "00000000-0000-0000-00000000"
  }
```




To authenticate users, you can send a POST request to `/login` with the following JSON data:

## Contributing

Contributions are welcome! If you'd like to contribute to this project, please follow these steps:

1. Fork the project.
2. Create a new branch for your feature or bug fix.
3. Make your changes.
4. Test your changes thoroughly.
5. Create a pull request.

## License

The project is licensed under the **MIT License**
