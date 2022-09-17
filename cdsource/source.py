import re
import csv
import json
import requests
import cmu_course_api


def __course_prof_name(name):
	regex0 = r"^([A-Za-z\-' ]+), ([A-Za-z\-'\.]+)$"
	regex1 = r"^([A-Za-z\-' ]+), ([A-Za-z\-'\.]+) ([A-Za-z\-'\.]+)$"

	match0 = re.search(regex0, name)
	if match0 is not None:
		last = match0.group(1)
		first = match0.group(2)
		return first + " " + last

	match1 = re.search(regex1, name)
	if match1 is not None:
		last = match1.group(1)
		first = match1.group(3)
		return first + " " + last

	return name


def generate_course(semester):
	"""
	Return a dict
		- key CID
		- value: a dict
			- key: name, desc, dept, units, prof, prereq, coreq
			- value: string/float
	"""
	# Read course data
	course_full = cmu_course_api.get_course_data(semester)

	# Write into result dict
	result = {}
	for cid, course in course_full["courses"].items():
		result[cid] = {}
		result[cid]["name"] = (
				course["name"]
				if "name" in course and course["name"] is not None
				else "Unknown")
		result[cid]["desc"] = (
				course["desc"]
				if "desc" in course and course["desc"] is not None
				else "Unknown")
		result[cid]["dept"] = (
				course["department"]
				if "department" in course and course["department"] is not None
				else "Unknown")
		result[cid]["prereq"] = (
				course["prereqs"]
				if "prereqs" in course and course["prereqs"] is not None
				else "")
		result[cid]["coreq"] = (
				course["coreqs"]
				if "coreqs" in course and course["coreqs"] is not None
				else "")
		result[cid]["units"] = (
				float(course["units"])
				if "units" in course
					and course["units"] is not None
					and course["units"] != ""
				else 0.0)

		prof_list = []
		if "lectures" in course:
			for lecture in course["lectures"]:
				if "instructors" not in lecture:
					continue
				for prof in lecture["instructors"]:
					prof_list.append(__course_prof_name(prof))

		result[cid]["prof"] = ";".join(prof_list)

	return result


def __fce_float_average(arr):
	f = []
	for elem in arr:
		if elem != "":
			f.append(float(elem))

	if len(f) == 0:
		return -1.0
	else:
		return sum(f) / len(f)


def generate_fce(filepath):
	"""
	Return a dict
		- key CID
		- value: a dict
			- key: FCEHours, FCETeachingRate, FCECourseRate, FCELevel
			- value: float/string/int
	"""
	# Read CSV file
	fce_full = {}
	f = open(filepath, newline="")
	reader = csv.DictReader(f)
	for row in reader:
		cid = row["Num"]
		cid = cid[:2] + "-" + cid[2:]

		if cid not in fce_full:
			fce_full[cid] = []
		fce_full[cid].append(row)
	f.close()

	# Write into result dict
	result = {}
	for cid in fce_full:
		result[cid] = {}
		FCEHours = []
		FCETeachingRate = []
		FCECourseRate = []
		FCELevel = "Unknown"
		FCEStudentCount = 0

		for elem in fce_full[cid]:
			FCEHours.append(elem["Hrs Per Week"])
			FCETeachingRate.append(elem["Overall teaching rate"])
			FCECourseRate.append(elem["Overall course rate"])
			FCELevel = elem["Course Level"]
			FCEStudentCount = int(elem["Total # Students"]) if elem["Total # Students"] != "" else 0

		result[cid]["FCEHours"] = __fce_float_average(FCEHours)
		result[cid]["FCETeachingRate"] = __fce_float_average(FCETeachingRate)
		result[cid]["FCECourseRate"] = __fce_float_average(FCECourseRate)
		result[cid]["FCELevel"] = FCELevel
		result[cid]["FCEStudentCount"] = FCEStudentCount

	return result


def __rmp_prof_name(name):
	regex0 = r"^([A-Za-z\-' ]+), ([A-Za-z\-'\.]+)$"



def generate_rmp(school_id):
	"""
	Return a dict
		- key: prof name
		- value: a dict
			- key: RMPRatingClass, RMPRatingOverall
			- value: string/float
	"""
	RMP_URL = "https://www.ratemyprofessors.com/filter/professor?page=%d&queryoption=TEACHER&queryBy=schoolId&sid=%d"
	result = {}
	page_id = 0
	while True:
		page_id += 1
		page = requests.get(RMP_URL % (page_id, school_id))
		page_obj = json.loads(page.content)

		for prof in page_obj["professors"]:
			first = prof["tFname"].strip()
			last = prof["tLname"].strip()
			name = first + " " + last

			result[name] = {}
			result[name]["RMPRatingClass"] = prof["rating_class"]
			result[name]["RMPRatingOverall"] = prof["overall_rating"]

		if page_obj["remaining"] == 0:
			break

	return result


if __name__ == "__main__":
	course_data = generate_course("F")
	fce_data = generate_fce("fce.csv")
	rmp_data = generate_rmp(181)

	with open("course_data.json", "w") as f:
		json.dump(course_data, f)

	with open("fce_data.json", "w") as f:
		json.dump(fce_data, f)

	with open("rmp_data.json", "w") as f:
		json.dump(rmp_data, f)
