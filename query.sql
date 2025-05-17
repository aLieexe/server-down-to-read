drop table schema_migrations;
drop table books;

CREATE TABLE IF NOT EXISTS books(
    title text,
    link text,
    repository_id serial,
    id serial PRIMARY KEY,
    link_expiry timestamptz DEFAULT NOW() + INTERVAL '1 weeks'
);


SELECT * FROM repository WHERE id = 'b8a1f35e-3f47-4c26-87f3-9a4b99b0bb12';