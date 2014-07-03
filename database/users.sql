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
-- Name: dm_users; Type: TABLE; Schema: public; Owner: Matthew; Tablespace: 
--

CREATE TABLE dm_users (
    user_id uuid NOT NULL,
    username varchar(256),
    password character varying(100),
    secret varchar(100)
);


ALTER TABLE public.dm_users OWNER TO "Matthew";

--
-- Data for Name: dm_users; Type: TABLE DATA; Schema: public; Owner: Matthew
--

COPY dm_users (user_id, name, password, pin) FROM stdin;
070a6749-95a5-42b3-982f-5cbcae602ec7	Matt	\N	42
077af5eb-cec8-4407-9aaf-beb657d5e8b2	Taylor	\N	42
60c7f77d-d95b-436c-b695-81f6a1ac4fbb	Jorge	\N	42
cb355cb5-a01a-4b12-870a-925ac9ba18dd	Brandon	\N	42
b3d7dd5e-4f0c-4e5e-9f6f-c17102297028	Zac	\N	42
2b4ce513-0182-41c8-8a06-1dc9b92baf82	Jared	\N	42
7a4fb8d2-d15c-4fe0-a233-ae9b8d0868d8	Sydney	\N	42
25031d69-bccf-4294-9486-600a1890e79e	Steve	\N	42
dc20d2c8-3844-4269-9d67-028b66ab412a	Aly	\N	42
628f13a4-aac9-43af-af99-75585ea72fec	Eden	\N	42
73cb6afb-7cb3-4407-a9fa-f3118e93121e	Jill	\N	42
aedafc03-3741-40ae-b14b-8ef32ce5d3f4	Angela	\N	42
610a0174-4fa0-4a59-b10c-8212c349cfd0	Torin	\N	\N
\.


--
-- Name: dm_users_pkey; Type: CONSTRAINT; Schema: public; Owner: Matthew; Tablespace: 
--

ALTER TABLE ONLY dm_users
    ADD CONSTRAINT dm_users_pkey PRIMARY KEY (user_id);


--
-- PostgreSQL database dump complete
--

