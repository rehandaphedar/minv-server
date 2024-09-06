PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS channels (
	   channelname text NOT NULL UNIQUE PRIMARY KEY,
	   created text NOT NULL DEFAULT CURRENT_TIMESTAMP,
	   password text NOT NULL
);

CREATE TABLE IF NOT EXISTS videos (
	   slug text NOT NULL UNIQUE PRIMARY KEY,
	   title text NOT NULL,
	   description text NOT NULL,
	   uploader text references channels(channelname) ON DELETE CASCADE NOT NULL ,
	   uploaded text NOT NULL DEFAULT CURRENT_TIMESTAMP,
	   processed integer NOT NULL DEFAULT 0
)
