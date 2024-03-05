-- Users Table
CREATE TABLE users (
                           id UNIQUEIDENTIFIER DEFAULT NEWID() PRIMARY KEY,
                           username VARCHAR(MAX) NOT NULL,
    email VARCHAR(MAX) NOT NULL,
    full_name VARCHAR(MAX) NOT NULL,
    hashed_password VARCHAR(MAX) NOT NULL,
    is_email_verified BIT NOT NULL DEFAULT 0,
    hashed_pin VARCHAR(MAX) DEFAULT NULL,
    created_at DATETIMEOFFSET DEFAULT SYSDATETIMEOFFSET(),
    updated_at DATETIMEOFFSET DEFAULT SYSDATETIMEOFFSET(),
    deleted_at DATETIMEOFFSET DEFAULT NULL
);

-- Accounts Table
CREATE TABLE accounts (
                              id BIGINT IDENTITY(1,1) PRIMARY KEY,
                              user_id UNIQUEIDENTIFIER FOREIGN KEY REFERENCES users(id),
                              balance BIGINT NOT NULL,
                              created_at DATETIMEOFFSET DEFAULT SYSDATETIMEOFFSET(),
                              updated_at DATETIMEOFFSET DEFAULT SYSDATETIMEOFFSET(),
                              deleted_at DATETIMEOFFSET DEFAULT NULL
);

-- Verify Emails Table
CREATE TABLE verify_emails (
                                  id BIGINT IDENTITY(1,1) PRIMARY KEY,
                                  user_id UNIQUEIDENTIFIER FOREIGN KEY REFERENCES users(id),
                                  email VARCHAR(MAX) NOT NULL,
    secret_code VARCHAR(MAX) NOT NULL,
    is_used BIT NOT NULL DEFAULT 0,
    expired_at DATETIMEOFFSET NOT NULL,
    created_at DATETIMEOFFSET DEFAULT SYSDATETIMEOFFSET(),
    updated_at DATETIMEOFFSET DEFAULT SYSDATETIMEOFFSET(),
    deleted_at DATETIMEOFFSET DEFAULT NULL
);

-- Session table
CREATE TABLE sessions (
                              id UNIQUEIDENTIFIER PRIMARY KEY,
                              user_id UNIQUEIDENTIFIER FOREIGN KEY REFERENCES users(id),
                              refresh_token NVARCHAR(MAX) NOT NULL,
                              user_agent VARCHAR(MAX) NOT NULL,
    client_ip VARCHAR(MAX) NOT NULL,
    is_blocked BIT NOT NULL DEFAULT 0,
    expired_at DATETIMEOFFSET NOT NULL,
    created_at DATETIMEOFFSET DEFAULT SYSDATETIMEOFFSET(),
    deleted_at DATETIMEOFFSET DEFAULT NULL
);

CREATE TABLE forgot_password (
                                    id BIGINT IDENTITY(1,1) PRIMARY KEY,
                                    user_id UNIQUEIDENTIFIER FOREIGN KEY REFERENCES users(id),
                                    reset_token NVARCHAR(MAX) NOT NULL,
                                    is_used BIT NOT NULL DEFAULT 0,
                                    request_timestamp DATETIMEOFFSET DEFAULT SYSDATETIMEOFFSET(),
                                    expiration_timestamp DATETIMEOFFSET,
                                    deleted_at DATETIMEOFFSET DEFAULT NULL
);

CREATE TABLE entries (
                             id BIGINT IDENTITY(1,1) PRIMARY KEY,
                             account_id BIGINT NOT NULL REFERENCES accounts(id),
                             date DATE NOT NULL DEFAULT GETDATE(),
                             amount BIGINT NOT NULL,
                             entry_type VARCHAR(MAX) NOT NULL,
    created_at DATETIMEOFFSET DEFAULT SYSDATETIMEOFFSET(),
    deleted_at DATETIMEOFFSET DEFAULT NULL
);

CREATE TABLE transfers (
                               id BIGINT IDENTITY(1,1) PRIMARY KEY,
                               from_account_id BIGINT NOT NULL REFERENCES accounts(id),
                               to_account_id BIGINT NOT NULL REFERENCES accounts(id),
                               amount BIGINT NOT NULL,
                               created_at DATETIMEOFFSET DEFAULT SYSDATETIMEOFFSET(),
                               deleted_at DATETIMEOFFSET DEFAULT NULL
);

CREATE TABLE dbo.expenses_plans (
                                   id BIGINT IDENTITY(1,1) PRIMARY KEY,
                                   user_id UNIQUEIDENTIFIER FOREIGN KEY REFERENCES dbo.users(id),
                                   title VARCHAR(MAX) NOT NULL,
    description VARCHAR(MAX) DEFAULT NULL,
    amount BIGINT NOT NULL,
    date DATE NOT NULL DEFAULT GETDATE(),
    created_at DATETIMEOFFSET DEFAULT SYSDATETIMEOFFSET(),
    updated_at DATETIMEOFFSET DEFAULT SYSDATETIMEOFFSET(),
    deleted_at DATETIMEOFFSET DEFAULT NULL
);

CREATE TABLE notifications (
                                   id BIGINT IDENTITY(1,1) PRIMARY KEY,
                                   user_id UNIQUEIDENTIFIER FOREIGN KEY REFERENCES users(id),
                                   title VARCHAR(MAX) NOT NULL,
    description VARCHAR(MAX) NOT NULL,
    category VARCHAR(MAX) NOT NULL,
    created_at DATETIMEOFFSET DEFAULT SYSDATETIMEOFFSET(),
    deleted_at DATETIMEOFFSET DEFAULT NULL
);
