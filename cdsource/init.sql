USE candela
CREATE TABLE course (
	cid 				INTEGER		NOT NULL PRIMARY KEY,
	description			MEDIUMTEXT	NOT NULL,
	dept				TEXT		NOT NULL,
	units				DECIMAL		NOT NULL,
	prof				TEXT		NULL,
	prereq				TEXT		NULL,
	coreq				TEXT		NULL,
	FCEHours			DECIMAL		NULL,
	FCETeachingRate		DECIMAL		NULL,
	FCECourseRate		DECIMAL		NULL,
	FCELevel			TEXT		NULL,
	FCEStudentCount		INTEGER		NULL
);

CREATE TABLE professor (
	name				VARCHAR(256)	NOT NULL PRIMARY KEY,
	RMPRatingClass		TEXT			NULL,
	RMPRatingOverall	DECIMAL			NULL
);
