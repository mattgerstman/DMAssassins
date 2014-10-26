---
--- Fixes game_id foreign keys not cascading deletes
---

---
--- Drop current constaint on dm_teams and create a new one that cascades deletes
---

alter table dm_teams drop constraint dm_teams_game_id_fkey;
alter table dm_teams add constraint dm_teams_game_id_fkey FOREIGN KEY (game_id) references dm_games(game_id) on delete cascade;

---
--- Drop current constaint on dm_kill_timers and create a new one that cascades deletes
---

alter table dm_kill_timers drop constraint dm_kill_timers_game_id_fkey;
alter table dm_kill_timers add constraint dm_kill_timers_game_id_fkey FOREIGN KEY (game_id) references dm_games(game_id) on delete cascade;

---
--- Drop current constaint on dm_user_targets and create a new one that cascades deletes
---

alter table dm_user_targets drop constraint dm_user_targets_game_id_fkey;
alter table dm_user_targets add constraint dm_user_targets_game_id_fkey FOREIGN KEY (game_id) references dm_games(game_id) on delete cascade;