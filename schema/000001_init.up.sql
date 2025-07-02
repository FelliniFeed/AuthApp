CREATE TABLE users (
    Id uuid default gen_random_uuid() PRIMARY KEY,
    UserName VARCHAR(255) UNIQUE NOT NULL,
    Password VARCHAR(255) NOT NULL
);

CREATE TABLE refresh_tokens (
    User_Id UUID NOT NULL,
    Refresh_Token VARCHAR(255) NOT NULL
);