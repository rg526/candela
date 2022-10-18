package cdmodel

type Course struct {
	CID					string
	Name				string
	Description			string
	Dept				string
	Units				float32

	Prereq				string
	Coreq				string
}


type FCE struct {
	CID					string
	Hours				float32
	TeachingRate		float32
	CourseRate			float32
	Level				string
	StudentCount		int
}
