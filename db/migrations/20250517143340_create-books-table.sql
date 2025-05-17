-- migrate:up
CREATE TABLE IF NOT EXISTS books(
    filename text NOT NULL,
    link text NOT NULL,
    repository_id uuid REFERENCES repository(id),
    id uuid PRIMARY KEY,
    link_expiry timestamptz DEFAULT NOW() + interval '1 weeks'
);

-- migrate:down
DROP TABLE books;

