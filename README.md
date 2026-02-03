# About
Current project is try to implement simple golang backend with using `jwt`, `chi`, `sqlc`, `migrate`, `swagger` libraries and much more.

Idea was to make some simple golang backend. Main orientation was online cinema, where only admin may manage films, but user may rate them, add favorites and write comments.

It's learning project so... nothing serious.

# How to test
1. Run `docker compose up`
2. Run migration 
```
docker compose exec migration /migrate/migrate.sh
```
3. Init admin user (Login: admin, Password: passwd)
```
docker compose exec migration /migrate/init_admin.sh
```
4. Go to the `http://localhost:8080/swagger/index.html` and play with swagger
