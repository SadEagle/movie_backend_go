insert into user_data(name, login, encoded_password, is_admin) 
select 'Sr. Admin', 'admin','\x251420b0f1fd1944c2d3c38536686039804d33d19a013872552ba771c950307d'::bytea, true 
where not exists(
  select null
  from user_data
  where login = 'admin'
);
