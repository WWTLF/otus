CREATE TABLE IF NOT EXISTS userlist (
   id serial PRIMARY KEY,
   username VARCHAR (256) UNIQUE NOT NULL,
   firstname VARCHAR (256) NOT NULL,
   lastname VARCHAR (256) NOT NULL,
   email VARCHAR (256) NOT NULL,
   phone VARCHAR (11) NOT NULL
);