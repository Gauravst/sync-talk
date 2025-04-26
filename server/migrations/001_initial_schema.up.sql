--- 001_initial_schema.up.sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  role TEXT DEFAULT 'USER',
  profilePic TEXT,
  createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chatRoom (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  description TEXT NOT NULL,
  private BOOLEAN DEFAULT TRUE,
  code TEXT,
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
  userId INTEGER NOT NULL,
  roomName TEXT NOT NULL,
  UNIQUE (userId, roomName),
  FOREIGN KEY (userId) REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (roomName) REFERENCES chatRoom (name) ON DELETE CASCADE
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
  chatRoom (name, description, userId, private)
VALUES
  (
    'general',
    'General discussion for all topics',
    1,
    false
  ),
  (
    'tech-talk',
    'Discussions about technology and programming',
    1,
    false
  ),
  (
    'golang',
    'Everything related to Go programming language',
    1,
    false
  ),
  (
    'react',
    'React.js discussions and help',
    1,
    false
  ),
  (
    'websockets',
    'WebSockets implementation and best practices',
    1,
    false
  ),
  (
    'gaming',
    'Gaming discussions and community',
    1,
    false
  ),
  (
    'music',
    'Music recommendations and discussions',
    1,
    false
  ),
  (
    'movies',
    'Movie discussions and recommendations',
    1,
    false
  ),
  (
    'books',
    'Book club and literature discussions',
    1,
    false
  ),
  ('design', 'UI/UX design discussions', 1, false),
  (
    'crypto',
    'Cryptocurrency and blockchain discussions',
    1,
    false
  ),
  (
    'fitness',
    'Fitness tips and motivation',
    1,
    false
  );
