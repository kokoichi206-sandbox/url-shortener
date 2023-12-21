CREATE TABLE shorturl (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    short TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO shorturl (url, short) VALUES ('https://www.google.com', 'google');
