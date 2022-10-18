import json
import mysql.connector
import config


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


# professor table
for name, prof in rmp_data.items():
	sql = "INSERT INTO rmp (name, ratingClass, ratingOverall) VALUES (%s, %s, %s)"
	val = (name, prof["RMPRatingClass"], prof["RMPRatingOverall"])
	cursor.execute(sql, val)
