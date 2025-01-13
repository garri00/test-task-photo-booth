BEGIN;

CREATE SCHEMA IF NOT EXISTS service;

ALTER SCHEMA service OWNER TO "serviceadmin";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION service.gen_random_uuid(
)
    RETURNS uuid
    LANGUAGE 'c'
    COST 1
    VOLATILE PARALLEL SAFE
AS '$libdir/pgcrypto', 'pg_random_uuid'
;

ALTER FUNCTION service.gen_random_uuid()
    OWNER TO "serviceadmin";

COMMIT;



