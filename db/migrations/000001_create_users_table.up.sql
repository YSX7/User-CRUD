CREATE TABLE IF NOT EXISTS users(
                                    id serial PRIMARY KEY,
                                    login VARCHAR (50) UNIQUE NOT NULL,
                                    password VARCHAR (100) NOT NULL
);