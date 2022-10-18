import re
import csv
import json
import requests
import cmu_course_api


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
			try:
				result[name]["RMPRatingOverall"] = float(prof["overall_rating"])
			except ValueError:
				result[name]["RMPRatingOverall"] = -1.0;

		if page_obj["remaining"] == 0:
			break

	return result


rmp_data = generate_rmp(181)
with open("rmp_data.json", "w") as f:
	json.dump(rmp_data, f)
