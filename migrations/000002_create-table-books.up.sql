CREATE TABLE IF NOT EXISTS books(
    filename text NOT NULL,
    link text NOT NULL,
    repository_id uuid references repository(id),
    id uuid PRIMARY KEY,
    link_expiry timestamptz DEFAULT NOW() + INTERVAL '1 weeks'
);

