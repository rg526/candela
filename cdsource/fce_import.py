import json
import mysql.connector
import config


with open("fce_data.json", "r") as f:
	fce_data = json.load(f)


conn = mysql.connector.connect(
	host = config.DB_HOST,
	user = config.DB_USER,
	password = config.DB_PWD,
	database = config.DB_NAME,
	autocommit = True
)
cursor = conn.cursor()


# fce table
for cid, fce in fce_data.items():
	sql = "INSERT INTO fce (cid, hours, teachingRate, courseRate, level, studentCount) VALUES (%s, %s, %s, %s, %s, %s)"
	val = (cid, fce["FCEHours"], fce["FCETeachingRate"], fce["FCECourseRate"], fce["FCELevel"], fce["FCEStudentCount"])
	cursor.execute(sql, val)
