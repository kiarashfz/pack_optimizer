CREATE TABLE IF NOT EXISTS packs (
    id SERIAL PRIMARY KEY,
    size INTEGER UNIQUE NOT NULL
);

-- Sample insert
INSERT INTO packs (size)
VALUES (250),
       (500),
       (1000),
       (2000),
       (5000);
