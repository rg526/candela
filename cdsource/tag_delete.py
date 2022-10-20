import json
import mysql.connector
import config
import sys
from goose3 import Goose

# Init goose and mysql
g = Goose()
conn = mysql.connector.connect(
	host = config.DB_HOST,
	user = config.DB_USER,
	password = config.DB_PWD,
	database = config.DB_NAME,
	autocommit = True
)
cursor = conn.cursor()

# Read tag data
with open("tag.json", "r") as f:
	tag_data = json.load(f)

# Fetch article
# Delete column
cid = tag_data["CID"]
content = tag_data["Content"]

sql = "DELETE FROM tag WHERE cid = %s AND content = %s"
val = (cid, content)
cursor.execute(sql, val)
