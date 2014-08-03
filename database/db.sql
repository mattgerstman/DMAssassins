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

--
-- Name: dm_user_role; Type: TYPE; Schema: public; Owner: dmassassins
--

CREATE TYPE dm_user_role AS ENUM (
    'dm_admin',
    'dm_captain',
    'dm_user'
);


ALTER TYPE public.dm_user_role OWNER TO "dmassassins";

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: dm_games; Type: TABLE; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE TABLE dm_games (
    game_id uuid NOT NULL,
    game_name character varying(100),
    game_started boolean
);


ALTER TABLE public.dm_games OWNER TO "dmassassins";

--
-- Name: dm_user_game_mapping; Type: TABLE; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE TABLE dm_user_game_mapping (
    user_id uuid,
    game_id uuid
);


ALTER TABLE public.dm_user_game_mapping OWNER TO "dmassassins";

--
-- Name: dm_user_properties; Type: TABLE; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE TABLE dm_user_properties (
    user_id uuid,
    key character varying(100),
    value bytea
);


ALTER TABLE public.dm_user_properties OWNER TO "dmassassins";

--
-- Name: dm_user_targets; Type: TABLE; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE TABLE dm_user_targets (
    user_id uuid,
    target_id uuid
);


ALTER TABLE public.dm_user_targets OWNER TO "dmassassins";

--
-- Name: dm_users; Type: TABLE; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE TABLE dm_users (
    user_id uuid NOT NULL,
    username character varying(256),
    secret character varying(100),
    alive boolean DEFAULT true,
    email character varying(256) DEFAULT ''::character varying NOT NULL,
    facebook_id bigint,
    facebook_token character varying,
    user_role dm_user_role
);


ALTER TABLE public.dm_users OWNER TO "dmassassins";

--
-- Data for Name: dm_games; Type: TABLE DATA; Schema: public; Owner: dmassassins
--

COPY dm_games (game_id, game_name, game_started) FROM stdin;
\.


--
-- Data for Name: dm_user_game_mapping; Type: TABLE DATA; Schema: public; Owner: dmassassins
--

COPY dm_user_game_mapping (user_id, game_id) FROM stdin;
\.


--
-- Data for Name: dm_user_properties; Type: TABLE DATA; Schema: public; Owner: dmassassins
--

COPY dm_user_properties (user_id, key, value) FROM stdin;
f7fe373c-aad2-4794-98da-701dc8ffbce9	Twitter	\\x687474703a2f2f747769747465722e636f6d2f6a696d6d79796973666c7979
f7fe373c-aad2-4794-98da-701dc8ffbce9	Facebook	\\x687474703a2f2f66616365626f6f6b2e636f6d2f6a696d6d79796973666c7979
f7fe373c-aad2-4794-98da-701dc8ffbce9	Name	\\x4a696d6d79
f7fe373c-aad2-4794-98da-701dc8ffbce9	photo_thumb	\\x68747470733a2f2f67726170682e66616365626f6f6b2e636f6d2f6a696d6d79796973666c79792f706963747572653f77696474683d333030266865696768743d333030
f7fe373c-aad2-4794-98da-701dc8ffbce9	photo	\\x68747470733a2f2f67726170682e66616365626f6f6b2e636f6d2f6a696d6d79796973666c79792f706963747572653f77696474683d31303030
9c83e902-4b0c-4ea1-8fe6-682c71be67c4	Facebook	\\x68747470733a2f2f66616365626f6f6b2e636f6d2f3130313532363232303230343831393133
9c83e902-4b0c-4ea1-8fe6-682c71be67c4	photo	\\x68747470733a2f2f67726170682e66616365626f6f6b2e636f6d2f31303135323632323032303438313931332f706963747572653f77696474683d31303030
9c83e902-4b0c-4ea1-8fe6-682c71be67c4	photo_thumb	\\x68747470733a2f2f67726170682e66616365626f6f6b2e636f6d2f31303135323632323032303438313931332f706963747572653f77696474683d333030266865696768743d333030
9c83e902-4b0c-4ea1-8fe6-682c71be67c4	first_name	\\x4d617474
9c83e902-4b0c-4ea1-8fe6-682c71be67c4	last_name	\\x47657273746d616e
d3a83c74-bc6f-430e-9afa-7b17a4b0ea5d	Facebook	\\x68747470733a2f2f66616365626f6f6b2e636f6d2f31353033363137383236353231383333
d3a83c74-bc6f-430e-9afa-7b17a4b0ea5d	photo	\\x68747470733a2f2f67726170682e66616365626f6f6b2e636f6d2f313530333631373832363532313833332f706963747572653f77696474683d31303030
d3a83c74-bc6f-430e-9afa-7b17a4b0ea5d	photo_thumb	\\x68747470733a2f2f67726170682e66616365626f6f6b2e636f6d2f313530333631373832363532313833332f706963747572653f77696474683d333030266865696768743d333030
d3a83c74-bc6f-430e-9afa-7b17a4b0ea5d	first_name	\\x4861727279
d3a83c74-bc6f-430e-9afa-7b17a4b0ea5d	last_name	\\x506f74746572
9e8a2fe5-424f-413a-b5a9-a36768b2ff96	Facebook	\\x68747470733a2f2f66616365626f6f6b2e636f6d2f31353137323631313538343839383232
9e8a2fe5-424f-413a-b5a9-a36768b2ff96	photo	\\x68747470733a2f2f67726170682e66616365626f6f6b2e636f6d2f313531373236313135383438393832322f706963747572653f77696474683d31303030
9e8a2fe5-424f-413a-b5a9-a36768b2ff96	photo_thumb	\\x68747470733a2f2f67726170682e66616365626f6f6b2e636f6d2f313531373236313135383438393832322f706963747572653f77696474683d333030266865696768743d333030
9e8a2fe5-424f-413a-b5a9-a36768b2ff96	first_name	\\x526f6e
9e8a2fe5-424f-413a-b5a9-a36768b2ff96	last_name	\\x576561736c6579
e9950266-0bf5-459e-9f43-0e23e7057e16	Facebook	\\x68747470733a2f2f66616365626f6f6b2e636f6d2f31343837323336303131353135363934
e9950266-0bf5-459e-9f43-0e23e7057e16	photo	\\x68747470733a2f2f67726170682e66616365626f6f6b2e636f6d2f313438373233363031313531353639342f706963747572653f77696474683d31303030
e9950266-0bf5-459e-9f43-0e23e7057e16	photo_thumb	\\x68747470733a2f2f67726170682e66616365626f6f6b2e636f6d2f313438373233363031313531353639342f706963747572653f77696474683d333030266865696768743d333030
e9950266-0bf5-459e-9f43-0e23e7057e16	first_name	\\x5279616e
e9950266-0bf5-459e-9f43-0e23e7057e16	last_name	\\x4c65776973
\.


--
-- Data for Name: dm_user_targets; Type: TABLE DATA; Schema: public; Owner: dmassassins
--

COPY dm_user_targets (user_id, target_id) FROM stdin;
e9950266-0bf5-459e-9f43-0e23e7057e16	d3a83c74-bc6f-430e-9afa-7b17a4b0ea5d
d3a83c74-bc6f-430e-9afa-7b17a4b0ea5d	f7fe373c-aad2-4794-98da-701dc8ffbce9
f7fe373c-aad2-4794-98da-701dc8ffbce9	9c83e902-4b0c-4ea1-8fe6-682c71be67c4
9c83e902-4b0c-4ea1-8fe6-682c71be67c4	9e8a2fe5-424f-413a-b5a9-a36768b2ff96
9e8a2fe5-424f-413a-b5a9-a36768b2ff96	e9950266-0bf5-459e-9f43-0e23e7057e16
\.


--
-- Data for Name: dm_users; Type: TABLE DATA; Schema: public; Owner: dmassassins
--

COPY dm_users (user_id, username, secret, alive, email, facebook_id, facebook_token, user_role) FROM stdin;
f7fe373c-aad2-4794-98da-701dc8ffbce9	Jimmy	muggle	f		1584513463	CAAFsDY1v8skBANZBZCwoZC6FfFi5o7ckncvgsLZCa3F3ryLbXNtLjLRhPnwaSFvyVd2JcnbmNRYLXsS73RtuHfQT688FZAeaAj3qBZCg7HjvWgnn3vq1EIU1NIXl033CjU5iDp7ZAxO2fWqSl5hNqs5Yk66FuMlkWrS4JvrD06aWeOE4vjoQpCfdvrAxyKtm2j5i3XaUOElrFjQdCSit8nhMfZAUm8D4LcIZD	dm_user
d3a83c74-bc6f-430e-9afa-7b17a4b0ea5d	HarryPotter	muggle	f	harry_kcoxeji_potter@tfbnw.net	1503617826521833	\N	dm_user
9e8a2fe5-424f-413a-b5a9-a36768b2ff96	RonWeasley	muggle	f	ron_uphoxsv_weasley@tfbnw.net	1517261158489822	\N	dm_user
e9950266-0bf5-459e-9f43-0e23e7057e16	RyanLewis	muggle	t	ryan_wiuwnyz_lewis@tfbnw.net	1487236011515694	\N	dm_user
9c83e902-4b0c-4ea1-8fe6-682c71be67c4	MattGerstman	muggle	f	imatt711@me.com	10152622020481913	\N	dm_admin
\.


--
-- Name: dm_games_pkey; Type: CONSTRAINT; Schema: public; Owner: dmassassins; Tablespace: 
--

ALTER TABLE ONLY dm_games
    ADD CONSTRAINT dm_games_pkey PRIMARY KEY (game_id);


--
-- Name: dm_users_pkey; Type: CONSTRAINT; Schema: public; Owner: dmassassins; Tablespace: 
--

ALTER TABLE ONLY dm_users
    ADD CONSTRAINT dm_users_pkey PRIMARY KEY (user_id);


--
-- Name: dm_user_game_mapping_game_id_idx; Type: INDEX; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE INDEX dm_user_game_mapping_game_id_idx ON dm_user_game_mapping USING btree (game_id);


--
-- Name: dm_user_game_mapping_user_id_game_id_idx; Type: INDEX; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE UNIQUE INDEX dm_user_game_mapping_user_id_game_id_idx ON dm_user_game_mapping USING btree (user_id, game_id);


--
-- Name: dm_user_properties_user_id_key_idx; Type: INDEX; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE UNIQUE INDEX dm_user_properties_user_id_key_idx ON dm_user_properties USING btree (user_id, key);


--
-- Name: dm_users_email_idx; Type: INDEX; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE UNIQUE INDEX dm_users_email_idx ON dm_users USING btree (username);


--
-- Name: dm_users_facebook_id_idx; Type: INDEX; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE UNIQUE INDEX dm_users_facebook_id_idx ON dm_users USING btree (facebook_id);


--
-- Name: unique_target; Type: INDEX; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE UNIQUE INDEX unique_target ON dm_user_targets USING btree (target_id);


--
-- Name: unique_user; Type: INDEX; Schema: public; Owner: dmassassins; Tablespace: 
--

CREATE UNIQUE INDEX unique_user ON dm_user_targets USING btree (user_id);


--
-- Name: dm_user_game_mapping_game_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dmassassins
--

ALTER TABLE ONLY dm_user_game_mapping
    ADD CONSTRAINT dm_user_game_mapping_game_id_fkey FOREIGN KEY (game_id) REFERENCES dm_games(game_id);


--
-- Name: dm_user_game_mapping_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dmassassins
--

ALTER TABLE ONLY dm_user_game_mapping
    ADD CONSTRAINT dm_user_game_mapping_user_id_fkey FOREIGN KEY (user_id) REFERENCES dm_users(user_id);


--
-- Name: foreign_target; Type: FK CONSTRAINT; Schema: public; Owner: dmassassins
--

ALTER TABLE ONLY dm_user_targets
    ADD CONSTRAINT foreign_target FOREIGN KEY (target_id) REFERENCES dm_users(user_id);


--
-- Name: foreign_user; Type: FK CONSTRAINT; Schema: public; Owner: dmassassins
--

ALTER TABLE ONLY dm_user_targets
    ADD CONSTRAINT foreign_user FOREIGN KEY (user_id) REFERENCES dm_users(user_id);


--
-- Name: foreign_user; Type: FK CONSTRAINT; Schema: public; Owner: dmassassins
--

ALTER TABLE ONLY dm_user_properties
    ADD CONSTRAINT foreign_user FOREIGN KEY (user_id) REFERENCES dm_users(user_id);


--
-- Name: public; Type: ACL; Schema: -; Owner: dmassassins
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM "dmassassins";
GRANT ALL ON SCHEMA public TO "dmassassins";
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

