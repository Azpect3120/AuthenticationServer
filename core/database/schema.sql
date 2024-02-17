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
