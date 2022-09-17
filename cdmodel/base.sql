USE candela;
CREATE TABLE course (
	cid 				INTEGER		NOT NULL PRIMARY KEY,
	description			MEDIUMTEXT	NOT NULL,
	dept				TEXT		NOT NULL,
	units				DECIMAL		NOT NULL,
	prof				TEXT		NOT NULL,
	prereq				TEXT		NOT NULL,
	coreq				TEXT		NOT NULL,
	FCEHours			DECIMAL		NOT NULL,
	FCETeachingRate		DECIMAL		NOT NULL,
	FCECourseRate		DECIMAL		NOT NULL,
	FCELevel			TEXT		NOT NULL,
	FCEStudentCount		INTEGER		NOT NULL
);

CREATE TABLE professor (
	name				VARCHAR(256)	NOT NULL PRIMARY KEY,
	RMPRatingClass		TEXT			NOT NULL,
	RMPRatingOverall	DECIMAL			NOT NULL
);
