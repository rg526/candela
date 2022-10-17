USE candela;

CREATE TABLE comment (
	commentID			INT				NOT NULL AUTO_INCREMENT PRIMARY KEY,
	cid					VARCHAR(128)	NOT NULL,
	uid					VARCHAR(128)	NOT NULL,
	content				MEDIUMTEXT		NOT NULL,
	time				TEXT			NOT NULL,
	anonymous			INT				NOT NULL,
	score				INT				NOT NULL DEFAULT 0
);

CREATE TABLE comment_reply (
	replyID				INT				NOT NULL AUTO_INCREMENT PRIMARY KEY,
	commentID			INT				NOT NULL,
	uid					VARCHAR(128)	NOT NULL,
	content				MEDIUMTEXT		NOT NULL,
	time				TEXT			NOT NULL,
	anonymous			INT				NOT NULL
);

CREATE TABLE comment_response (
	responseID			INT				NOT NULL AUTO_INCREMENT PRIMARY KEY,
	commentID			INT				NOT NULL,
	uid					VARCHAR(128)	NOT NULL,
	time				TEXT			NOT NULL
);
