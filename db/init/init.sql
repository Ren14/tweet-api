-- Enable the UUID extension if it's not already enabled
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the Users table
CREATE TABLE IF NOT EXISTS users
(
    id
    UUID
    PRIMARY
    KEY
    DEFAULT
    uuid_generate_v4
(
),
    username VARCHAR
(
    255
) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW
(
),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW
(
)
    );

-- Create the Tweets table
CREATE TABLE IF NOT EXISTS tweets
(
    id
    UUID
    PRIMARY
    KEY
    DEFAULT
    uuid_generate_v4
(
),
    user_id UUID NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW
(
),

    -- Foreign key constraint to link tweets to users
    CONSTRAINT fk_user
    FOREIGN KEY
(
    user_id
)
    REFERENCES users
(
    id
)
    ON DELETE CASCADE -- If a user is deleted, their tweets are also deleted
    );

-- Create the Follows table for the many-to-many relationship
CREATE TABLE IF NOT EXISTS follows
(
    follower_id
    UUID
    NOT
    NULL,
    following_id
    UUID
    NOT
    NULL,
    created_at
    TIMESTAMPTZ
    NOT
    NULL
    DEFAULT
    NOW
(
),

    -- Foreign key for the user who is following
    CONSTRAINT fk_follower
    FOREIGN KEY
(
    follower_id
)
    REFERENCES users
(
    id
)
    ON DELETE CASCADE,

    -- Foreign key for the user who is being followed
    CONSTRAINT fk_following
    FOREIGN KEY
(
    following_id
)
    REFERENCES users
(
    id
)
    ON DELETE CASCADE,

    -- Composite primary key to ensure a user can't follow the same person more than once
    PRIMARY KEY
(
    follower_id,
    following_id
)
    );

-- Create indexes for faster lookups on foreign keys
CREATE INDEX IF NOT EXISTS idx_tweets_user_id ON tweets(user_id);
CREATE INDEX IF NOT EXISTS idx_follows_follower_id ON follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_follows_following_id ON follows(following_id);

-- =================================================================
-- Seed Data for Testing
-- =================================================================

-- Insert some users with pre-defined UUIDs for easy reference
INSERT INTO users (id, username)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'nachito'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'agus'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'diegores'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'messi'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'renzito'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'tinchomanik'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a17', 'maurini'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a18', 'togencio'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a19', 'tucuboss'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a20', 'semidios'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a21', 'subzero'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'tucudev'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a23', 'danielito'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a24', 'fran');

-- Insert some follow relationships
INSERT INTO follows (follower_id, following_id)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12');
INSERT INTO follows (follower_id, following_id)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');
INSERT INTO follows (follower_id, following_id)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');

-- Insert some tweets from our users
INSERT INTO tweets (user_id, content)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Hello world! This is my first tweet.'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Just setting up my twttr.'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Loving this new platform! The performance is amazing.'),
       ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Is anyone out there? #golang');