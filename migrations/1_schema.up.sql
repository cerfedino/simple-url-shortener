CREATE TABLE long_urls (
  id SERIAL PRIMARY KEY,
  long_url VARCHAR NOT NULL
);

CREATE TABLE shortened_urls (
  short_url VARCHAR UNIQUE,
  long_url_id INT,
  private BOOLEAN NOT NULL DEFAULT FALSE,
  PRIMARY KEY (short_url, long_url_id),
  FOREIGN KEY (long_url_id) REFERENCES long_urls(id) ON DELETE CASCADE
);

CREATE TABLE log (
  id BIGSERIAL PRIMARY KEY,
  timestamp TIMESTAMP NOT NULL,
  ip VARCHAR(15) NOT NULL,
  success BOOLEAN NOT NULL,
  shortURL VARCHAR NOT NULL,
  redirectURL VARCHAR
);
