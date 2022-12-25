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
    email     VARCHAR(250) UNIQUE             NOT NULL CHECK ( email != '' ),
    password  VARCHAR(250)                    NOT NULL CHECK ( OCTET_LENGTH(password) != 0 )
);

COMMENT ON TABLE gophkeeper.credentials IS 'credentials for user authentication';
COMMENT ON COLUMN gophkeeper.credentials.uuid IS 'unique credentials ID';
COMMENT ON COLUMN gophkeeper.credentials.user_uuid IS 'user unique ID';
COMMENT ON COLUMN gophkeeper.credentials.email IS 'user login as login';
COMMENT ON COLUMN gophkeeper.credentials.password IS 'user password';
CREATE UNIQUE INDEX credentials_email_uindex
    ON gophkeeper.credentials (email);

CREATE TYPE gophkeeper.data_types AS ENUM ('ACCOUNT', 'TEXT', 'BINARY','CREDIT_CARD');
COMMENT ON TYPE gophkeeper.data_types IS 'allowed types of data';

CREATE TABLE gophkeeper.sync
(
    uuid       uuid        DEFAULT uuid_generate_v4() NOT NULL
        CONSTRAINT sync_pk
            PRIMARY KEY,
    user_uuid  uuid                                   NOT NULL
        CONSTRAINT sync_users_null_fk
            REFERENCES gophkeeper.users (uuid),
    machine_id VARCHAR(250) UNIQUE                    NOT NULL CHECK ( machine_id != '' ),
    data_uuid  uuid                                   NOT NULL,
    data_type  data_types                             NOT NULL,
    name       uuid                                   NOT NULL,
    sync_uuid  uuid                                   NOT NULL,
    created_at timestamptz DEFAULT NOW()              NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz
);

COMMENT ON TABLE gophkeeper.sync IS 'credentials for user authentication';
COMMENT ON COLUMN gophkeeper.sync.uuid IS 'unique sync ID';
COMMENT ON COLUMN gophkeeper.sync.user_uuid IS 'user unique ID';
COMMENT ON COLUMN gophkeeper.sync.machine_id IS 'client machine ID';
COMMENT ON COLUMN gophkeeper.sync.data_uuid IS 'data unique ID';
COMMENT ON COLUMN gophkeeper.sync.data_type IS 'type of data';
COMMENT ON COLUMN gophkeeper.sync.name IS 'name of data record, must be unique across user and data type';
COMMENT ON COLUMN gophkeeper.sync.sync_uuid IS 'sync ID marker to check collisions';
CREATE UNIQUE INDEX sync_data_uuid_uindex
    ON gophkeeper.sync (data_uuid);
CREATE UNIQUE INDEX sync_user_datatype_name_uindex
    ON gophkeeper.sync (user_uuid, data_type, name) WHERE deleted_at IS NOT NULL;
CREATE INDEX sync_machine_id_index
    ON gophkeeper.sync (machine_id);