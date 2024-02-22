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

### Create an Application

### Update an Application (Part)

### Update an Application (Full)

### Delete an Application

## User Examples

### Get All Users

### Get a User

### Create a User

### Update a User

### Delete a User
