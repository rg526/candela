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
cid = page_data["CID"]
link = page_data["URL"]
article = g.extract(url = link)
priority = page_data.get("Priority", 0)

sql = "INSERT INTO page (cid, title, link, content, priority) VALUES (%s, %s, %s, %s, %s)"
val = (cid, article.title, link, article.cleaned_text, priority)
cursor.execute(sql, val)
