/*
    This file runs as part of the docker build process. If you're making 
    changes to this file, you'll have to rebuild the docker image.
*/


/* * * * * * * * * * * * * * * * * * * * * *
 *
 *          SCHEMA
 *
 * * * * * * * * * * * * * * * * * * * * * */

CREATE SCHEMA IF NOT EXISTS public
;

CREATE TABLE public.users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    nickname VARCHAR(50),
    user_type VARCHAR(50) NOT NULL DEFAULT 'UTYPE_USER',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    message_count INT DEFAULT 0
)
;

CREATE TABLE public.user_types (
    id SERIAL PRIMARY KEY,
    type_key VARCHAR(50) NOT NULL,
    permission_bitfield BIT(8) NOT NULL DEFAULT B'00000000'
)
;

/*
    Uncertain on the best way to handle this, but for now, we'll just create a table that stores messages
    and link them to users via a foreign key. This will allow us to easily query messages by user and 
    also enforce referential integrity. If a user is deleted, all their messages will be deleted as well.
    This is a simple design, but it should work for the purposes of this challenge.
    NOTE: Assumed that messages are text-based and do not require complex formatting and that we do not need to store
    the recipient of the message.
*/
CREATE TABLE public.messages (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
;


/* * * * * * * * * * * * * * * * * * * * * *
 *
 *          DATA
 *
 * * * * * * * * * * * * * * * * * * * * * */

INSERT INTO public.users 
    (username, nickname, email, user_type, message_count)
 VALUES 
    ('Liam', 'L dawg', 'liam@email.com', 'UTYPE_ADMIN', 1),
    ('Jon', NULL, 'jon@email.com', 'UTYPE_USER', 2),
    ('Myles', 'Big M', 'myles@email.com', 'UTYPE_USER', 2)
;

INSERT INTO public.user_types 
    (type_key, permission_bitfield)
 VALUES 
    ('UTYPE_USER',      B'00000000'),
    ('UTYPE_ADMIN',     B'10000000'),
    ('UTYPE_ADMIN',     B'10000000'),
    ('UTYPE_MODERATOR', B'01000000')
;

INSERT INTO public.messages 
    (user_id, content)
 VALUES 
    (1, 'Hello, this is Liam!'),
    (2, 'Hi, I am Jon.'),
    (3, 'Myles here!'),
    (2, 'Another message from Jon.'),
    (3, 'Myles again with another message.')
;