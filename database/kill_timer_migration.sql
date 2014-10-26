--
-- Migration To Create the dm_kill_timers table
--

--
-- Create dm_kill_timers
--
CREATE TABLE dm_kill_timers (game_id uuid REFERENCES dm_games(game_id), create_ts timestamp default now(), execute_ts timestamp NOT NULL);

--
-- Create index for unique game_id
--
CREATE UNIQUE INDEX single_kill_timer ON dm_kill_timers(game_id);