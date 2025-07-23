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
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
;

CREATE TABLE public.user_types (
    id SERIAL PRIMARY KEY,
    type_key VARCHAR(50) NOT NULL,
    permission_bitfield BIT(8) NOT NULL DEFAULT B'00000000'
)
;


/* * * * * * * * * * * * * * * * * * * * * *
 *
 *          DATA
 *
 * * * * * * * * * * * * * * * * * * * * * */

INSERT INTO public.users 
    (username, nickname, email, user_type)
 VALUES 
    ('Liam', 'L dawg', 'liam@email.com', 'UTYPE_ADMIN'),
    ('Jon', NULL, 'jon@email.com', 'UTYPE_USER'),
    ('Myles', 'Big M', 'myles@email.com', 'UTYPE_USER')
;

INSERT INTO public.user_types 
    (type_key, permission_bitfield)
 VALUES 
    ('UTYPE_USER',      B'00000000'),
    ('UTYPE_ADMIN',     B'10000000'),
    ('UTYPE_ADMIN',     B'10000000'),
    ('UTYPE_MODERATOR', B'01000000')
;
