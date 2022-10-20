"""
tag.json should contain the following information:
- CID: string, course ID
- Content: string
- Priority: int, optional, priority of the link
"""
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
cid = tag_data["CID"]
content = tag_data["Content"]
priority = tag_data.get("Priority", 0)

sql = "INSERT INTO tag (cid, content, priority) VALUES (%s, %s, %s)"
val = (cid, content, priority)
cursor.execute(sql, val)
