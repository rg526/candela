import re
import csv
import json
import requests
import cmu_course_api


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
			- key: FCEHours, FCETeachingRate, FCECourseRate, FCELevel, FCEStudentCount
			- value: float/string/int
	"""
	# Read CSV file
	fce_full = {}
	f = open(filepath, newline="")
	reader = csv.DictReader(f)
	for row in reader:
		cid = row["Num"]
		cid = str(cid)

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


fce_data = generate_fce("fce.csv")
with open("fce_data.json", "w") as f:
	json.dump(fce_data, f)
