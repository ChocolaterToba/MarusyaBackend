--
-- PostgreSQL database dump
--

-- Dumped from database version 14.1 (Ubuntu 14.1-201-yandex.52142.672784f35a)
-- Dumped by pg_dump version 14.3 (Ubuntu 14.3-1.pgdg20.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: account; Type: TABLE; Schema: public; Owner: mikhail
--

CREATE TABLE public.account (
    user_id integer NOT NULL,
    application_id text NOT NULL,
    session_id text DEFAULT ''::text NOT NULL,
    username text,
    current_question_id integer DEFAULT 0 NOT NULL,
    past_questions numeric[] DEFAULT '{}'::numeric[] NOT NULL,
    past_answers jsonb DEFAULT '[]'::jsonb NOT NULL
);


ALTER TABLE public.account OWNER TO mikhail;

--
-- Name: question; Type: TABLE; Schema: public; Owner: mikhail
--

CREATE TABLE public.question (
    question_id integer NOT NULL,
    test_id integer DEFAULT 0 NOT NULL,
    text text DEFAULT ''::text NOT NULL,
    next_question_ids jsonb DEFAULT '{}'::jsonb NOT NULL,
    question_in_test_id numeric DEFAULT 0 NOT NULL
);


ALTER TABLE public.question OWNER TO mikhail;

--
-- Name: question_question_id_seq; Type: SEQUENCE; Schema: public; Owner: mikhail
--

CREATE SEQUENCE public.question_question_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.question_question_id_seq OWNER TO mikhail;

--
-- Name: question_question_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mikhail
--

ALTER SEQUENCE public.question_question_id_seq OWNED BY public.question.question_id;


--
-- Name: quiz; Type: TABLE; Schema: public; Owner: mikhail
--

CREATE TABLE public.quiz (
    title text DEFAULT ''::text NOT NULL,
    backtracking_enabled boolean DEFAULT true NOT NULL,
    calculate_correctness boolean DEFAULT false NOT NULL,
    id integer NOT NULL
);


ALTER TABLE public.quiz OWNER TO mikhail;

--
-- Name: quiz_id_seq; Type: SEQUENCE; Schema: public; Owner: mikhail
--

CREATE SEQUENCE public.quiz_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.quiz_id_seq OWNER TO mikhail;

--
-- Name: quiz_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mikhail
--

ALTER SEQUENCE public.quiz_id_seq OWNED BY public.quiz.id;


--
-- Name: user_user_id_seq; Type: SEQUENCE; Schema: public; Owner: mikhail
--

CREATE SEQUENCE public.user_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_user_id_seq OWNER TO mikhail;

--
-- Name: user_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mikhail
--

ALTER SEQUENCE public.user_user_id_seq OWNED BY public.account.user_id;


--
-- Name: account user_id; Type: DEFAULT; Schema: public; Owner: mikhail
--

ALTER TABLE ONLY public.account ALTER COLUMN user_id SET DEFAULT nextval('public.user_user_id_seq'::regclass);


--
-- Name: question question_id; Type: DEFAULT; Schema: public; Owner: mikhail
--

ALTER TABLE ONLY public.question ALTER COLUMN question_id SET DEFAULT nextval('public.question_question_id_seq'::regclass);


--
-- Name: quiz id; Type: DEFAULT; Schema: public; Owner: mikhail
--

ALTER TABLE ONLY public.quiz ALTER COLUMN id SET DEFAULT nextval('public.quiz_id_seq'::regclass);


--
-- Name: account account_pk; Type: CONSTRAINT; Schema: public; Owner: mikhail
--

ALTER TABLE ONLY public.account
    ADD CONSTRAINT account_pk PRIMARY KEY (user_id);


--
-- Name: question question_pk; Type: CONSTRAINT; Schema: public; Owner: mikhail
--

ALTER TABLE ONLY public.question
    ADD CONSTRAINT question_pk PRIMARY KEY (question_id);


--
-- Name: question question_un; Type: CONSTRAINT; Schema: public; Owner: mikhail
--

ALTER TABLE ONLY public.question
    ADD CONSTRAINT question_un UNIQUE (test_id, question_in_test_id);


--
-- Name: quiz quiz_pk; Type: CONSTRAINT; Schema: public; Owner: mikhail
--

ALTER TABLE ONLY public.quiz
    ADD CONSTRAINT quiz_pk PRIMARY KEY (id);


--
-- Name: user_application_id_idx; Type: INDEX; Schema: public; Owner: mikhail
--

CREATE UNIQUE INDEX user_application_id_idx ON public.account USING btree (application_id);


--
-- Name: user_session_id_idx; Type: INDEX; Schema: public; Owner: mikhail
--

CREATE INDEX user_session_id_idx ON public.account USING btree (session_id) WHERE (session_id <> ''::text);


--
-- Name: account account_fk; Type: FK CONSTRAINT; Schema: public; Owner: mikhail
--

ALTER TABLE ONLY public.account
    ADD CONSTRAINT account_fk FOREIGN KEY (current_question_id) REFERENCES public.question(question_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: question question_fk; Type: FK CONSTRAINT; Schema: public; Owner: mikhail
--

ALTER TABLE ONLY public.question
    ADD CONSTRAINT question_fk FOREIGN KEY (test_id) REFERENCES public.quiz(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

INSERT INTO public.quiz (id, title)
    VALUES (0, 'Изначальный тест - только для корня!')
    ON CONFLICT DO NOTHING;

INSERT INTO public.question (question_id, question_in_test_id, test_id, text)
    VALUES (0, 0, 0, 'Выбери тест')
    ON CONFLICT DO NOTHING;

