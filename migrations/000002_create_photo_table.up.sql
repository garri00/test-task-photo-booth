CREATE TABLE IF NOT EXISTS service.photos
(
    id           UUID PRIMARY KEY DEFAULT service.gen_random_uuid(),
    data_origin TEXT,
    data_75     TEXT,
    data_50     TEXT,
    data_25     TEXT,
    is_deleted   BOOLEAN NOT NULL
);
ALTER TABLE service.photos
    OWNER TO "serviceadmin";

