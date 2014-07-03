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
-- Name: dm_user_targets; Type: TABLE; Schema: public; Owner: Matthew; Tablespace: 
--

CREATE TABLE dm_user_targets (
    user_id uuid,
    target_id uuid
);


ALTER TABLE public.dm_user_targets OWNER TO "Matthew";

--
-- Data for Name: dm_user_targets; Type: TABLE DATA; Schema: public; Owner: Matthew
--

COPY dm_user_targets (user_id, target_id) FROM stdin;
070a6749-95a5-42b3-982f-5cbcae602ec7	610a0174-4fa0-4a59-b10c-8212c349cfd0
610a0174-4fa0-4a59-b10c-8212c349cfd0	dc20d2c8-3844-4269-9d67-028b66ab412a
dc20d2c8-3844-4269-9d67-028b66ab412a	628f13a4-aac9-43af-af99-75585ea72fec
628f13a4-aac9-43af-af99-75585ea72fec	cb355cb5-a01a-4b12-870a-925ac9ba18dd
cb355cb5-a01a-4b12-870a-925ac9ba18dd	73cb6afb-7cb3-4407-a9fa-f3118e93121e
73cb6afb-7cb3-4407-a9fa-f3118e93121e	b3d7dd5e-4f0c-4e5e-9f6f-c17102297028
b3d7dd5e-4f0c-4e5e-9f6f-c17102297028	077af5eb-cec8-4407-9aaf-beb657d5e8b2
077af5eb-cec8-4407-9aaf-beb657d5e8b2	60c7f77d-d95b-436c-b695-81f6a1ac4fbb
60c7f77d-d95b-436c-b695-81f6a1ac4fbb	7a4fb8d2-d15c-4fe0-a233-ae9b8d0868d8
7a4fb8d2-d15c-4fe0-a233-ae9b8d0868d8	aedafc03-3741-40ae-b14b-8ef32ce5d3f4
aedafc03-3741-40ae-b14b-8ef32ce5d3f4	25031d69-bccf-4294-9486-600a1890e79e
25031d69-bccf-4294-9486-600a1890e79e	2b4ce513-0182-41c8-8a06-1dc9b92baf82
2b4ce513-0182-41c8-8a06-1dc9b92baf82	070a6749-95a5-42b3-982f-5cbcae602ec7
\.


--
-- Name: unique_target; Type: INDEX; Schema: public; Owner: Matthew; Tablespace: 
--

CREATE UNIQUE INDEX unique_target ON dm_user_targets USING btree (target_id);


--
-- Name: unique_user; Type: INDEX; Schema: public; Owner: Matthew; Tablespace: 
--

CREATE UNIQUE INDEX unique_user ON dm_user_targets USING btree (user_id);


--
-- Name: foreign_target; Type: FK CONSTRAINT; Schema: public; Owner: Matthew
--

ALTER TABLE ONLY dm_user_targets
    ADD CONSTRAINT foreign_target FOREIGN KEY (target_id) REFERENCES dm_users(user_id);


--
-- Name: foreign_user; Type: FK CONSTRAINT; Schema: public; Owner: Matthew
--

ALTER TABLE ONLY dm_user_targets
    ADD CONSTRAINT foreign_user FOREIGN KEY (user_id) REFERENCES dm_users(user_id);


--
-- PostgreSQL database dump complete
--

