import json
import mysql.connector
import config


with open("course_data.json", "r") as f:
	course_data = json.load(f)


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
	sql = "INSERT INTO course (cid, name, description, dept, units, prof, prereq, coreq) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)"
	val = (cid, course["name"], course["description"], course["dept"], course["units"], course["prof"], course["prereq"], course["coreq"])
	cursor.execute(sql, val)
