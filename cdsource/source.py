import csv
import cmu_course_api


def generate_course():
	return cmu_course_api.get_course_data("F")


def __fce_float_average(arr):
	f = []
	for elem in arr:
		if elem != "":
			f.append(float(elem))

	if len(f) == 0:
		return -1
	else:
		return sum(f) / len(f)


def generate_fce(filepath):
	"""
	Return a dict
		- key CID
		- value: a dict
			- key: hours, FCETeachingRate, FCECourseRate
			- value: float
	"""
	fce_full = {}
	f = open(filepath, newline="")
	reader = csv.DictReader(f)
	for row in reader:
		cid = row["Num"]
		cid = cid[:2] + "-" + cid[2:]

		if cid not in fce_full:
			fce_full[cid] = []

		fce_full[cid].append(row)

	result = {}
	for cid in fce_full:
		result[cid] = {}
		hours = []
		FCETeachingRate = []
		FCECourseRate = []

		for elem in fce_full[cid]:
			hours.append(elem["Hrs Per Week"])
			FCETeachingRate.append(elem["Overall teaching rate"])
			FCECourseRate.append(elem["Overall course rate"])

		result[cid]["hours"] = __fce_float_average(hours)
		result[cid]["FCETeachingRate"] = __fce_float_average(FCETeachingRate)
		result[cid]["FCECourseRate"] = __fce_float_average(FCECourseRate)


if __name__ == "__main__":
	#course_data = generate_course()
	fce_data = generate_fce("fce.csv")
