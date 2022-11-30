CREATE SCHEMA IF NOT EXISTS gophkeeper;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE gophkeeper.users
(
    uuid       uuid        DEFAULT uuid_generate_v4() NOT NULL
        CONSTRAINT users_pk
            PRIMARY KEY,
    created_at timestamptz DEFAULT NOW()              NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE gophkeeper.users IS 'users of gophkeepers of the seven keys';
COMMENT ON COLUMN gophkeeper.users.uuid IS 'unique user ID';
COMMENT ON COLUMN gophkeeper.users.created_at IS 'timestamp when user were created';
COMMENT ON COLUMN gophkeeper.users.updated_at IS 'timestamp when user were updated last time';

CREATE TABLE gophkeeper.credentials
(
    uuid      uuid DEFAULT uuid_generate_v4() NOT NULL
        CONSTRAINT credentials_pk
            PRIMARY KEY,
    user_uuid uuid                            NOT NULL
        CONSTRAINT credentials_users_null_fk
            REFERENCES gophkeeper.users (uuid),
    email     VARCHAR(250) UNIQUE             NOT NULL CHECK ( email <> '' ),
    password  VARCHAR(250)                    NOT NULL CHECK ( OCTET_LENGTH(password) <> 0 )
);

COMMENT ON TABLE gophkeeper.credentials IS 'credentials for user authentication';
COMMENT ON COLUMN gophkeeper.credentials.uuid IS 'unique credentials ID';
COMMENT ON COLUMN gophkeeper.credentials.user_uuid IS 'user unique ID';
COMMENT ON COLUMN gophkeeper.credentials.email IS 'user login as login';
COMMENT ON COLUMN gophkeeper.credentials.password IS 'user password';
CREATE UNIQUE INDEX credentials_email_uindex
    ON gophkeeper.credentials (email);