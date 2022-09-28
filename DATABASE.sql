CREATE SEQUENCE companies_id_seq;
CREATE TABLE companies (
    id INTEGER PRIMARY KEY DEFAULT nextval('companies_id_seq'),
    name TEXT NOT NULL UNIQUE,
    code TEXT NOT NULL,
    country TEXT NOT NULL,
    website TEXT,
    phone TEXT
);
