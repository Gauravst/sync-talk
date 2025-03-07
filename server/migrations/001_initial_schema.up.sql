--- 001_initial_schema.up.sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  role TEXT DEFAULT 'USER',
  createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chatRoom (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  members INTEGER,
  private BOOLEAN DEFAULT TRUE,
  userId INTEGER NOT NULL,
  createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_user FOREIGN KEY (userId) REFERENCES users (id)
);

CREATE TABLE loginSession (
  id SERIAL PRIMARY KEY,
  token TEXT NOT NULL,
  userId INTEGER NOT NULL,
  createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE groupMembers (
  id SERIAL PRIMARY KEY,
  userId INT NOT NULL,
  roomName TEXT NOT NULL,
  UNIQUE (userId, roomName)
);

CREATE TABLE messages (
  id SERIAL PRIMARY KEY,
  roomName TEXT NOT NULL,
  userId INT NOT NULL,
  content TEXT NOT NULL,
  createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO
  users (username, password)
VALUES
  ('test_user', 'test_password');

INSERT INTO
  chatRoom (name, members, description, userId)
VALUES
  (
    'general',
    24,
    'General discussion for all topics',
    1
  ),
  (
    'tech-talk',
    18,
    'Discussions about technology and programming',
    1
  ),
  (
    'golang',
    12,
    'Everything related to Go programming language',
    1
  ),
  ('react', 15, 'React.js discussions and help', 1),
  (
    'websockets',
    8,
    'WebSockets implementation and best practices',
    1
  ),
  (
    'gaming',
    32,
    'Gaming discussions and community',
    1
  ),
  (
    'music',
    20,
    'Music recommendations and discussions',
    1
  ),
  (
    'movies',
    16,
    'Movie discussions and recommendations',
    1
  ),
  (
    'books',
    14,
    'Book club and literature discussions',
    1
  ),
  ('design', 22, 'UI/UX design discussions', 1),
  (
    'crypto',
    19,
    'Cryptocurrency and blockchain discussions',
    1
  ),
  ('fitness', 17, 'Fitness tips and motivation', 1);
