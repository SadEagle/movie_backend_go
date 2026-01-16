DROP TABLE IF EXISTS user_data;
DROP TABLE IF EXISTS movie;
DROP TABLE IF EXISTS favorite_movie;
DROP TABLE IF EXISTS rated_movie;

DROP MATERIALIZED VIEW IF EXISTS movie_rating_view;
DROP INDEX IF EXISTS movie_rating_view_index;
