CREATE TABLE IF NOT EXISTS repository(
    name text,
    id uuid PRIMARY KEY,
    created timestamptz DEFAULT NOW()
);

