CREATE TABLE user_data(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR NOT NULL,
  login VARCHAR NOT NULL UNIQUE,
  encoded_password BYTEA NOT NULL,
  is_admin BOOL NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE movie(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE favorite(
  user_id UUID REFERENCES user_data ON DELETE CASCADE,
  movie_id UUID REFERENCES movie,
  PRIMARY KEY( user_id, movie_id)
);

CREATE TABLE rating(
  user_id UUID REFERENCES user_data ON DELETE CASCADE,
  movie_id UUID REFERENCES movie ON DELETE CASCADE,
  rating SMALLINT NOT NULL CHECK(rating BETWEEN 1 AND 10),
  PRIMARY KEY( user_id, movie_id)
);

CREATE TABLE comment(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES user_data,
  movie_id UUID REFERENCES movie ON DELETE CASCADE,
  text VARCHAR NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE MATERIALIZED VIEW total_rating_mview AS
SELECT movie_id, COUNT(*) AS amount_rates, AVG(rating) AS rating
FROM rating
GROUP BY movie_id;

CREATE UNIQUE INDEX total_rating_mview_index ON total_rating_mview(movie_id);
