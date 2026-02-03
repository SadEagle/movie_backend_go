#! /bin/sh
DB_HOST="dev-db"
DB_PASSWORD=$(cat "$DB_PASSWORD_FILE")
SSLMODE=${SSLMODE:-disable}
# Init admin user 
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" << EOF

insert into user_data(name, login, encoded_password, is_admin) 
select 'Sr. Admin', 'admin','\x251420b0f1fd1944c2d3c38536686039804d33d19a013872552ba771c950307d'::bytea, true 
where not exists(
  select null
  from user_data
  where login = 'admin'
);
EOF

