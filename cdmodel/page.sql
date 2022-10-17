USE candela;

CREATE TABLE page (
	pageID				INT				NOT NULL AUTO_INCREMENT PRIMARY KEY,
	cid					VARCHAR(128)	NOT NULL,
	title				TEXT			NOT NULL,
	link				TEXT			NOT NULL,
	content				MEDIUMTEXT		NOT NULL,
	priority			INT				NOT NULL
);
