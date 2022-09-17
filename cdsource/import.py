import json
import mysql.connector
import config


if __name__ == "__main__":
	with open("course_data.json", "r") as f:
		course_data = json.load(f)

	with open("fce_data.json", "r") as f:
		fce_data = json.load(f)

	with open("rmp_data.json", "r") as f:
		rmp_data = json.load(f)
	

	conn = mysql.connector.connect(
		host = config.DB_HOST,
		user = config.DB_USER,
		password = config.DB_PWD,
		database = config.DB_NAME,
		autocommit = True
	)
	cursor = conn.cursor()


	# course table
	for cid, course in course_data.items():
		sql = "INSERT INTO course (cid, name, description, dept, units, prof, prereq, coreq, FCEHours, FCETeachingRate, FCECourseRate, FCELevel, FCEStudentCount) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"
		if cid not in fce_data:
			val = (cid, course["name"], course["description"], course["dept"], course["units"], course["prof"], course["prereq"], course["coreq"], -1.0, -1.0, -1.0, "Unknown", -1)
		else:
			fce = fce_data[cid]
			val = (cid, course["name"], course["description"], course["dept"], course["units"], course["prof"], course["prereq"], course["coreq"], fce["FCEHours"], fce["FCETeachingRate"], fce["FCECourseRate"], fce["FCELevel"], fce["FCEStudentCount"])
		cursor.execute(sql, val)


	# professor table
	for name, prof in rmp_data.items():
		sql = "INSERT INTO professor (name, RMPRatingClass, RMPRatingOverall) VALUES (%s, %s, %s)"
		val = (name, prof["RMPRatingClass"], prof["RMPRatingOverall"])
		cursor.execute(sql, val)
