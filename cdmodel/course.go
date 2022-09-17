package cdmodel

type Course struct {
	cid					int
	description			string
	dept				string
	units				float32
	prof				string

	prereq				string
	coreq				string

	FCEHours			float32
	FCETeachingRate		float32
	FCECourseRate		float32
	FCELevel			string
	FCEStudentCount		int
}
