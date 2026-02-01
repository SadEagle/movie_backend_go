# About
Current project is try to implement simple golang backend of online movie theater with `chi`, `sqlc`, `migrate`, `swag`.

Currently it contains 5 tables: movie, user, favorite, rating, comment. We may add user, movie and favorite_movie/rated_movie relations.

# How to test
1. Run `docker compose up`
2. Run `docker compose exec backend ./scripts/migrate.bash ` for `dev-db` database migration
3. Go to the `http://localhost:8080/swagger/index.html` and play with swagger
