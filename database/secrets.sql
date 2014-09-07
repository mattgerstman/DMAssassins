 alter table dm_user_game_mapping add column secret varchar(100) not null default '';
alter table dm_users drop column secret;
