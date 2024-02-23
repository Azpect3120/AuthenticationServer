# Authentication Server Example Documentation

## Table of Contents
-   [Application Examples](#application-examples)
    -   [Get All Applications](#get-all-applications)
    -   [Get an Application](#get-an-application)
    -   [Create an Application](#create-an-application)
    -   [Update an Application (Part)](#update-an-application-(part))
    -   [Update an Application (Full)](#update-an-application-(full))
    -   [Delete an Application](#delete-an-application)
-   [User Examples](#user-examples)
    -   [Get All Users](#get-all-users)
    -   [Get a User](#get-a-user)
    -   [Create a User](#create-a-user)
    -   [Update a User](#update-a-user)
    -   [Delete a User](#delete-a-user)

## Application Examples

### Get All Applications

Gets all applications in the database.

```http
GET /v2/applications
```

Example Response:
```json
{
  "applications": [
    {
      "id": "7da29198-4c18-48b4-893a-3ae4b2ddcbc0",
      "name": "Test Application",
      "columns": [
        "id",
        "applicationid",
        "username",
        "password",
        "email",
        "createdat",
        "lastupdatedat"
      ],
      "createdat": "2024-02-17T04:00:24.716294Z",
      "lastupdatedat": "2024-02-17T04:00:24.716294Z"
    }
  ],
  "count": 1,
  "status": 200
}
```

### Get an Application

Gets an application using its ID.

```http
GET /v2/applications/:id
```

Example Response:
```json
{
  "application": {
    "id": "7da29198-4c18-48b4-893a-3ae4b2ddcbc0",
    "name": "Test Application",
    "columns": [
      "id",
      "applicationid",
      "username",
      "password",
      "email",
      "createdat",
      "lastupdatedat"
    ],
    "createdat": "2024-02-17T04:00:24.716294Z",
    "lastupdatedat": "2024-02-17T04:00:24.716294Z"
  },
  "status": 200
}
```

### Create an Application

Creates a new application. List of valid columns can be found in the [README](https://github.com/Azpect3120/AuthenticationServer/blob/v2.0.0/README.md#overview).

```http
POST /v2/applications
```

Request Body:

-  `name` (string) : The name of the application. **REQUIRED**
-  `columns` ([]string) : The columns of the application. **REQUIRED**

```json
{
  "name": "Test Application",
  "columns": [
    "username",
    "password",
    "email"
  ]
}
```
Example Response:
```json
{
  "application": {
    "id": "7da29198-4c18-48b4-893a-3ae4b2ddcbc0",
    "name": "Test Application",
    "columns": [
      "id",
      "applicationid",
      "username",
      "password",
      "email",
      "createdat",
      "lastupdatedat"
    ],
    "createdat": "2024-02-17T04:00:24.716294Z",
    "lastupdatedat": "2024-02-17T04:00:24.716294Z"
  },
  "status": 201
}
```

### Update an Application (Part)
Updates part of an application. Can be used to ADD columns or update the name.
```http
PATCH /v2/applications/:id
```

Request Body:
- `name` (string) : The name of the application.
- `columns` ([]string) : The columns of the application.

```json
{
  "name": "Test Application (UPDATED)",
  "columns": [
      "data"
  ]
}
```

Example Response:
```json
{
  "application": {
    "id": "7da29198-4c18-48b4-893a-3ae4b2ddcbc0",
    "name": "Test Application (UPDATED)",
    "columns": [
      "id",
      "applicationid",
      "username",
      "password",
      "first",
      "last",
      "data",
      "createdat",
      "lastupdatedat"
    ],
    "createdat": "2024-02-17T04:00:24.716294Z",
    "lastupdatedat": "2024-02-23T19:21:22.832849Z"
  },
  "message": "",
  "status": 200
}
```

### Update an Application (Full)
Updates an entire application. Can be used to OVERWRITE columns or update the name.
Columns provided will be set as the new columns, any columns not provided will be removed.

```http
PUT /v2/applications/:id
```

Request Body:
- `name` (string) : The name of the application. **REQUIRED**
- `columns` ([]string) : The columns of the application. **REQUIRED(())

```json
{
  "name": "New Test Application",
  "columns": [
      "first",
      "last"
  ]
}
```

Example Response:
```json
{
  "application": {
    "id": "7da29198-4c18-48b4-893a-3ae4b2ddcbc0",
    "name": "Test Application (UPDATED)",
    "columns": [
      "id",
      "applicationid",
      "username",
      "password",
      "first",
      "last",
      "data",
      "createdat",
      "lastupdatedat"
    ],
    "createdat": "2024-02-17T04:00:24.716294Z",
    "lastupdatedat": "2024-02-23T19:21:22.832849Z"
  },
  "message": "",
  "status": 200
}
```

Example Response:
```json
{
  "application": {
    "id": "7da29198-4c18-48b4-893a-3ae4b2ddcbc0",
    "name": "New Test Application",
    "columns": [
      "id",
      "applicationid",
      "first",
      "last",
      "createdat",
      "lastupdatedat"
    ],
    "createdat": "2024-02-17T04:00:24.716294Z",
    "lastupdatedat": "2024-02-23T19:24:22.669494Z"
  },
  "message": "",
  "status": 200
}
```

### Delete an Application
Deletes an application using its ID.

```http
DELETE /v2/applications/:id
```

No response is returned if successful. (Status 204)

## User Examples

### Get All Users

### Get a User

### Create a User

### Update a User

### Delete a User
