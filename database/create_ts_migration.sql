--
-- Migration To Add Create Timestamps to Every Row of Every Current Table 
--

--
-- Add create_ts to table dm_game_properties
--
alter table dm_game_properties add column create_ts timestamp default current_timestamp;

--
-- Add create_ts to table dm_games
--
alter table dm_games add column create_ts timestamp default current_timestamp;

--
-- Add create_ts to table dm_teams
--
alter table dm_teams add column create_ts timestamp default current_timestamp;

--
-- Add create_ts to table dm_user_game_mapping
--
alter table dm_user_game_mapping add column create_ts timestamp default current_timestamp;

--
-- Add create_ts to table dm_user_properties
--
alter table dm_user_properties add column create_ts timestamp default current_timestamp;

--
-- Add create_ts to table dm_user_targets
--
alter table dm_user_targets add column create_ts timestamp default current_timestamp;

--
-- Add create_ts to table dm_users
--
alter table dm_users add column create_ts timestamp default current_timestamp;