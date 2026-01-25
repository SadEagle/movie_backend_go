CREATE TABLE IF NOT EXISTS user_data(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR NOT NULL,
  login VARCHAR NOT NULL UNIQUE,
  password VARCHAR NOT NULL,
  is_admin BOOL NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS movie(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS favorite_movie(
  user_id UUID REFERENCES user_data ON DELETE CASCADE,
  movie_id UUID REFERENCES movie ON DELETE CASCADE,
  PRIMARY KEY( user_id, movie_id)
);

CREATE TABLE IF NOT EXISTS rated_movie(
  user_id UUID REFERENCES user_data ON DELETE CASCADE,
  movie_id UUID REFERENCES movie ON DELETE CASCADE,
  rating SMALLINT NOT NULL CHECK(rating BETWEEN 1 AND 10),
  PRIMARY KEY( user_id, movie_id)
);

CREATE TABLE IF NOT EXISTS movie_comment(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES user_data ON DELETE CASCADE,
  movie_id UUID REFERENCES movie ON DELETE CASCADE,
  text VARCHAR NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE MATERIALIZED VIEW IF NOT EXISTS movie_rating_view AS
SELECT movie_id, COUNT(*) AS amount_rates, AVG(rating) AS rating
FROM rated_movie
GROUP BY movie_id;

CREATE UNIQUE INDEX movie_rating_view_index ON movie_rating_view(movie_id);
