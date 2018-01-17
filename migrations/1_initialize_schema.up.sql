CREATE TABLE records (
  id        SERIAL PRIMARY KEY,
  date      TIMESTAMP DEFAULT NOW(),
  comment   TEXT,
  sender    VARCHAR(100),
  document  VARCHAR(200),
  processed BOOLEAN   DEFAULT FALSE,
  escalated BOOLEAN   DEFAULT FALSE
);

CREATE TABLE pages (
  id        SERIAL PRIMARY KEY,
  record_id INT REFERENCES records ON DELETE CASCADE,
  content   TEXT
);
