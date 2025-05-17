-- migrate:up
CREATE TABLE IF NOT EXISTS repository(
    name text,
    id uuid PRIMARY KEY,
    created timestamptz DEFAULT NOW()
);

-- migrate:down
DROP TABLE repository;

