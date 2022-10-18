import json
import mysql.connector
import config


with open("course_data.json", "r") as f:
	course_data = json.load(f)
with open("prof_data.json", "r") as f:
	prof_data = json.load(f)


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
	sql = "INSERT INTO course (cid, name, description, dept, units, prereq, coreq) VALUES (%s, %s, %s, %s, %s, %s, %s)"
	val = (cid, course["name"], course["description"], course["dept"], course["units"], course["prereq"], course["coreq"])
	cursor.execute(sql, val)

# prof table
for cid, prof_list in prof_data.items():
	for prof in prof_list:
		sql = "INSERT INTO prof (cid, name) VALUES (%s, %s)"
		val = (cid, prof)
		cursor.execute(sql, val)
