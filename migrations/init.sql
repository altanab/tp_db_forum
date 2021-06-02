CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS Users (
    nickname CITEXT PRIMARY KEY,
    fullname VARCHAR(128) NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    about TEXT
);

CREATE TABLE IF NOT EXISTS Forums (
    title TEXT NOT NULL,
    username CITEXT NOT NULL REFERENCES Users (nickname),
    slug CITEXT PRIMARY KEY,
    posts INTEGER DEFAULT 0,
    threads INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS Threads (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author CITEXT NOT NULL REFERENCES Users (nickname),
    forum CITEXT NOT NULL REFERENCES Forums (slug),
    message TEXT NOT NULL,
    votes INTEGER DEFAULT 0,
    slug CITEXT UNIQUE,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS Posts (
    id SERIAL PRIMARY KEY,
    parent INTEGER DEFAULT NULL REFERENCES Posts (id),
    author CITEXT NOT NULL REFERENCES Users (nickname),
    message TEXT NOT NULL,
    is_edited BOOLEAN DEFAULT FALSE,
    forum CITEXT NOT NULL REFERENCES Forums (slug),
    thread INTEGER NOT NULL REFERENCES Threads (id),
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    m_path INTEGER[] DEFAULT ARRAY []::INTEGER []
);

CREATE TABLE IF NOT EXISTS Votes (
    nickname CITEXT NOT NULL REFERENCES Users (nickname),
    voice INTEGER NOT NULL,
    thread INTEGER NOT NULL REFERENCES Threads (id),
    UNIQUE (nickname, thread)
);

CREATE TABLE IF NOT EXISTS Forum_users (
    forum_user CITEXT NOT NULL REFERENCES Users (nickname),
    forum CITEXT NOT NULL REFERENCES Forums (slug),
    UNIQUE (forum_user, forum)
);

CREATE OR REPLACE FUNCTION count_forum_threads()
    RETURNS TRIGGER
    AS $count_forum_threads$
BEGIN
    UPDATE forums SET threads=(threads+1) WHERE LOWER(slug)=LOWER(NEW.forum);
    RETURN NEW;
END;
$count_forum_threads$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_vote()
    RETURNS TRIGGER
    AS $insert_vote$
BEGIN
    UPDATE threads SET votes=(votes+NEW.voice) WHERE id=NEW.thread;
    RETURN NEW;
END;
$insert_vote$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_vote()
    RETURNS TRIGGER
    AS $update_vote$
BEGIN
    IF OLD.voice <> NEW.voice THEN
        UPDATE threads SET votes=(votes+NEW.voice*2) WHERE id=NEW.thread;
    END IF;
    RETURN NEW;
END;
$update_vote$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_forum_users()
    RETURNS TRIGGER
    AS $update_forum_users$
BEGIN
    INSERT INTO forum_users(forum_user, forum)
    VALUES(NEW.author, NEW.forum)
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$update_forum_users$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_m_path()
    RETURNS TRIGGER
    AS $update_m_path$
DECLARE
    parent_path INTEGER[];
    first_parent_thread INT;
BEGIN
    IF NEW.parent is NULL THEN
        NEW.m_path := array_append(NEW.m_path, NEW.id);
    ELSE
        SELECT m_path FROM posts WHERE id=NEW.parent INTO parent_path;
        SELECT thread FROM posts WHERE id = parent_path[1] INTO first_parent_thread;
        IF NOT FOUND OR first_parent_thread != NEW.thread THEN
            RAISE EXCEPTION 'parent not found';
        END IF;
        NEW.m_path := NEW.m_path || parent_path || new.id;
    END IF;
    UPDATE forums SET posts=(posts+1) WHERE LOWER(slug)=LOWER(NEW.forum);
    RETURN NEW;
END;
$update_m_path$ LANGUAGE plpgsql;


CREATE TRIGGER new_vote
    AFTER INSERT
    ON votes
    FOR EACH ROW
    EXECUTE PROCEDURE insert_vote();

CREATE TRIGGER upd_vote
    AFTER UPDATE
    ON votes
    FOR EACH ROW
    EXECUTE PROCEDURE update_vote();

CREATE TRIGGER new_thread
    AFTER INSERT
    ON threads
    FOR EACH ROW
    EXECUTE PROCEDURE update_forum_users();

CREATE TRIGGER new_post
    AFTER INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE update_forum_users();

CREATE TRIGGER add_thread
    AFTER INSERT
    ON threads
    FOR EACH ROW
    EXECUTE PROCEDURE count_forum_threads();

CREATE TRIGGER path_update
    BEFORE INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE update_m_path();

CREATE INDEX IF NOT EXISTS users_nickname_email_index ON users (LOWER(nickname), LOWER(email));

CREATE INDEX IF NOT EXISTS thread_slug_index ON threads (LOWER(slug));
CREATE INDEX IF NOT EXISTS forum_threads_index ON threads (LOWER(forum), created);

CREATE INDEX IF NOT EXISTS vote_index ON votes (LOWER(nickname), thread);

CREATE UNIQUE INDEX forum_users_unique ON forum_users (forum_user, forum);

CREATE INDEX IF NOT EXISTS posts_id_path1_index ON posts (id, (m_path[1]))
CREATE INDEX IF NOT EXISTS posts_thread_id_index ON posts (thread, id)
CREATE INDEX IF NOT EXISTS posts_thread_path_id_index ON posts(thread, m_path, id)
CREATE INDEX IF NOT EXISTS posts_thread_index ON posts (thread)
CREATE INDEX IF NOT EXISTS posts_path1_index ON posts (m_path[1])
