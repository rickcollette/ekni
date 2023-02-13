CREATE TABLE clients (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(255) NOT NULL,
  ip VARCHAR(255) NOT NULL,
  key VARCHAR(255) NOT NULL
);



CREATE TABLE WebUser (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    mfa BOOLEAN NOT NULL,
    active BOOLEAN NOT NULL,
    admin BOOLEAN NOT NULL
);


CREATE TABLE web_user_clients (
  web_user_id INTEGER NOT NULL,
  client_id INTEGER NOT NULL,
  PRIMARY KEY (web_user_id, client_id),
  FOREIGN KEY (web_user_id) REFERENCES web_users (id),
  FOREIGN KEY (client_id) REFERENCES clients (id)
);