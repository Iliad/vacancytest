CREATE SEQUENCE vacancy_id;
CREATE TABLE IF NOT EXISTS vacancies
(
    id INT DEFAULT NEXTVAL('vacancy_id') PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    salary INT NOT NULL,
    experience TEXT NOT NULL,
    city TEXT NOT NULL
);
CREATE SEQUENCE user_id;
CREATE TYPE account_role AS ENUM ('viewer', 'editor');
CREATE TABLE IF NOT EXISTS users
(
    id INT DEFAULT NEXTVAL('user_id') PRIMARY KEY NOT NULL,
    role account_role NOT NULL DEFAULT 'viewer',
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);