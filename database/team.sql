--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: dm_teams; Type: TABLE; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE TABLE dm_teams (
    team_id uuid NOT NULL,
    game_id uuid NOT NULL,
    team_name character varying(100)
);


ALTER TABLE public.dm_teams OWNER TO dmassassins;

--
-- Data for Name: dm_teams; Type: TABLE DATA; Schema: public; Owner: dmassassins
--

COPY dm_teams (team_id, game_id, team_name) FROM stdin;
fe9c4819-3509-11e4-9a6b-685b35b45205	9202fcd2-ccbd-42d4-8c54-99968a38e5e6	Tech
02cb5f5d-350a-11e4-9a6b-685b35b45205	9202fcd2-ccbd-42d4-8c54-99968a38e5e6	Morale
a4e7ad2a-9f40-4f35-9357-b9f7c3671c6a	9202fcd2-ccbd-42d4-8c54-99968a38e5e6	Art/Layout
7e170b2c-c0b2-48a5-85c2-0991abf9b7ae	9202fcd2-ccbd-42d4-8c54-99968a38e5e6	Operations
\.


--
-- Name: dm_teams_pkey; Type: CONSTRAINT; Schema: public; Owner: dmassassins; Tablespace: 
--

ALTER TABLE ONLY dm_teams
    ADD CONSTRAINT dm_teams_pkey PRIMARY KEY (team_id);


--
-- Name: dm_teams_game_id_team_name_idx; Type: INDEX; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE UNIQUE INDEX dm_teams_game_id_team_name_idx ON dm_teams USING btree (game_id, team_name);


--
-- Name: dm_teams_game_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dmassassins
--

ALTER TABLE ONLY dm_teams
    ADD CONSTRAINT dm_teams_game_id_fkey FOREIGN KEY (game_id) REFERENCES dm_games(game_id);


--
-- PostgreSQL database dump complete
--

