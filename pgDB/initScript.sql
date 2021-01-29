CREATE SCHEMA goapp;

CREATE TABLE goapp.todolist
(
    id character(36) NOT NULL,
    title character varying(50) NOT NULL,
    todo_description character varying(200),
    due_date timestamp without time zone NOT NULL DEFAULT CURRENT_DATE,
    todo_priority character varying(10),
    iscomplete boolean DEFAULT false,
    CONSTRAINT todo_pkey PRIMARY KEY (title),
    CONSTRAINT check_value CHECK (todo_priority = 'high' or todo_priority = 'low' or todo_priority = 'medium')
);

CREATE UNIQUE INDEX uniq_todo_title ON goapp.todolist USING btree(title);


CREATE TABLE goapp.auth
(
    username character varying(50) NOT NULL,
    password_hash character varying(1024) NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (username)
);
