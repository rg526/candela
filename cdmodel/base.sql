USE candela;
CREATE TABLE course (
	cid 				VARCHAR(128)	NOT NULL PRIMARY KEY,
	name				TEXT			NOT NULL,
	description			MEDIUMTEXT		NOT NULL,
	dept				TEXT			NOT NULL,
	units				DECIMAL			NOT NULL,
	prereq				TEXT			NOT NULL,
	coreq				TEXT			NOT NULL
);

CREATE TABLE fce (
	cid					VARCHAR(128)	NOT NULL PRIMARY KEY,
	hours				DECIMAL			NOT NULL,
	teachingRate		DECIMAL			NOT NULL,
	courseRate			DECIMAL			NOT NULL,
	level				TEXT			NOT NULL,
	studentCount		INTEGER			NOT NULL
);

CREATE TABLE prof (
	cid					VARCHAR(128)	NOT NULL,
	name				VARCHAR(128)	NOT NULL,
	PRIMARY KEY(cid, name)
);

CREATE TABLE rmp (
	name				VARCHAR(128)	NOT NULL PRIMARY KEY,
	ratingClass			TEXT			NOT NULL,
	ratingOverall		DECIMAL			NOT NULL
);
