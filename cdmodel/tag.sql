USE candela;

CREATE TABLE tag (
	tagID				INT				NOT NULL AUTO_INCREMENT PRIMARY KEY,
	cid					VARCHAR(128)	NOT NULL,
	content				TEXT			NOT NULL,
	priority			INT				NOT NULL DEFAULT 0,
	uid                 VARCHAR(128)    NOT NULL
);
