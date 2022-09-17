package cdmodel

type Course struct {
	CID					int
	Name				string
	Description			string
	Dept				string
	Units				float32
	Prof				string

	Prereq				string
	Coreq				string

	FCEHours			float32
	FCETeachingRate		float32
	FCECourseRate		float32
	FCELevel			string
	FCEStudentCount		int
}
