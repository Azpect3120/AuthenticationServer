-- @block
CREATE TABLE IF NOT EXISTS Applications (
    ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    Name TEXT UNIQUE,
    Key UUID DEFAULT uuid_generate_v4()
);

CREATE TABLE IF NOT EXISTS Users (
    ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ApplicationID UUID,
    Username TEXT,
    Password TEXT,
    FOREIGN KEY (ApplicationID) REFERENCES Applications(ID)
);

-- @block
SELECT * FROM Applications;
SELECT * FROM Users;

-- @block
DROP TABLE Users;
DROP TABLE Applications;