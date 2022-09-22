USE candela;
CREATE TABLE user (
	uid					VARCHAR(128)	NOT NULL PRIMARY KEY,
	name				TEXT			NOT NULL,
	givenName			TEXT			NOT NULL,
	familyName			TEXT			NOT NULL,
	Email				TEXT			NOT NULL
);

CREATE TABLE token (
	token				VARCHAR(128)	NOT NULL PRIMARY KEY,
	uid					VARCHAR(128)	NOT NULL,
	time				TEXT			NOT NULL
);
