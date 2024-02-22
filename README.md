<div align="center">
  <img alt="Locked Server Image" height="300" src="./docs/server.png" style="margin: 40px;" />

  [![GoLang](https://img.shields.io/badge/Go_1.21.5-white.svg?style=for-the-badge&logo=Go)](https://go.dev/)
  [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-lightblue.svg?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
  [![Gin Web Framework](https://img.shields.io/badge/gin-yellow.svg?style=for-the-badge&logo=gin)](https://gin-gonic.com/)

# Authentication Server Version 2.0.0
<p align="left">
  A simple user authentication server built with the Gin framework, PostgreSQL database, 
  and UUID-keys for user authentication. This server app is designed to authenticate 
  users using simple HTTP requests. It provides a secure and efficient way to verify 
  user identities and grant access to protected resources.
</p>
</div>

## 🗝 Table of Contents
-   [Features](#🗝-features)
-   [Getting Started](#🗝-getting-started)
    -   [Prerequisites](#prerequisites)
    -   [Installation](#installation)
    -   [Download Binary](#download-binary)
    -   [Database](#database)
-   [Usage](#🗝-usage)
    -  [Overview](#overview)
    -  [Applications](#applications)
    -  [Users](#users)
-   [Contributing](#🗝-contributing)
-   [License](#🗝-license)

## 🗝 Features
-   **User Authentication**: Authenticate users using a unique UUID-key.
-   **Application Management**: Create, read, update, and delete applications.
-   **User Management**: Create, read, update, and delete user accounts.
-   **Password Security**: Secure user passwords using a simple yet effective hashing algorithm.
-   **Logging**: Log user and application activity to the stdout.

## 🗝 Getting Started
Follow these steps to get the project up and running on your local machine.

### Prerequisites
-   **Go** (+1.16)
-   **PostgreSQL** database
-   **Git** (optional)

### Installation
1.  Clone the repository to your local machine using the following command:
```bash
git clone https://github.com/Azpect3120/AuthenticationServer.git && cd AuthenticationServer
```

2. Install the project dependencies using the following command:
```bash
go mod tidy
```

3. Setup your PostgreSQL database and configure the database connection in the `.env` file:
```.env
# This url can found in the dashboard of most PSQL hosts or can be constructed using the required pieces
# REQUIRED
DB_URL=postgresql://username:password@localhost:5432/Database

# The port the server will listen on. Default is 3000
# OPTIONAL
AUTH_SERVER_PORT=3000
```

4. Build and run the server:
```bash
go build -o ./bin/server ./cmd/main.go && ./bin/server
# or 
go run ./cmd/main.go
```

### Download Binary
If you do not have Go installed on your machine, you can download the binary from the releases page.
Select the appropriate binary for your operating system and architecture, then run the binary in your terminal.

### Database
Once the server is up and running you will need to connect to a PostgreSQL database.
If you would like the code to work out of the box, you may copy the database schema provided below.

```sql 
-- UUID extension for use in creating and storing
-- UUID value types
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table for storing applications. The applications
-- do not have any duplicate restraints beyond their
-- id (uuid).
CREATE TABLE IF NOT EXISTS applications (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT,
    columns TEXT[],
    createdat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    lastupdatedat TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table for storing users. The users do not have any
-- duplicate restraints beyond their id (uuid).
CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    applicationid UUID REFERENCES applications(id),
    username TEXT,
    firstname TEXT,
    lastname TEXT,
    fullname TEXT,
    email TEXT,
    password TEXT,
    data TEXT, -- Stringified JSON
    createdat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    lastupdatedat TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🗝 Usage

### Overview

### Applications

### Users


## 🗝 Contributing
This project is open source, therefore contributions are encouraged! If you'd like to contribute to this project, please follow these steps:

1. Fork the project.
1. Create a new branch for your feature or bug fix.
1. Make your changes.
1. Test your changes thoroughly.
1. Create a pull request.

## 🗝 License
This project is licensed under azpect3120 the **MIT License**

View [LICENSE](https://github.com/azpect3120/AuthenticationServer/blob/v2.0.0/LICENSE)
