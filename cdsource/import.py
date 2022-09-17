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

	for cid, course in course_data.items():
		if cid not in fce_data:
			sql = "INSERT INTO course (cid, description, dept, units, prof, prereq, coreq) VALUES (%s, %s, %s, %s, %s, %s, %s)"
			val = (cid, course["description"], course["dept"], course["units"], course["prof"], course["prereq"], course["coreq"])
			cursor.execute(sql, val)
		else:
			fce = fce_data[cid]
			sql = "INSERT INTO course (cid, description, dept, units, prof, prereq, coreq, FCEHours, FCETeachingRate, FCECourseRate, FCELevel, FCEStudentCount) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"
			val = (cid, course["description"], course["dept"], course["units"], course["prof"], course["prereq"], course["coreq"], fce["FCEHours"], fce["FCETeachingRate"], fce["FCECourseRate"], fce["FCELevel"], fce["FCEStudentCount"])
			cursor.execute(sql, val)
