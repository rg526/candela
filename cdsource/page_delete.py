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

# Read page data
with open("page.json", "r") as f:
	page_data = json.load(f)

# Fetch article
# Delete column
cid = page_data["CID"]
link = page_data["URL"]


sql = "DELETE FROM page WHERE cid = %s AND link = %s"
val = (cid, link)
cursor.execute(sql, val)
