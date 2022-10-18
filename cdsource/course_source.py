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
	Return course dict:
		- key CID
		- value: a dict
			- key: name, description, dept, units, prof, prereq, coreq
			- value: string/float
	Return prof dict:
		- key CID
		- value: a list containing prof name in a list
	"""
	# Read course data
	course_full = cmu_course_api.get_course_data(semester)

	# Write into result dict
	result = {}
	prof_result = {}
	for cid, course in course_full["courses"].items():
		cid = cid[:2] + cid[3:]
		cid = str(cid)

		result[cid] = {}
		result[cid]["name"] = (
				course["name"]
				if "name" in course and course["name"] is not None
				else "Unknown")
		result[cid]["description"] = (
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

		prof_list = set()
		if "lectures" in course:
			for lecture in course["lectures"]:
				if "instructors" not in lecture:
					continue
				for prof in lecture["instructors"]:
					prof_list.add(__course_prof_name(prof))

		prof_result[cid] = list(prof_list)

	return result, prof_result


course_data, prof_data = generate_course("F")
with open("course_data.json", "w") as f:
	json.dump(course_data, f)
with open("prof_data.json", "w") as f:
	json.dump(prof_data, f)
