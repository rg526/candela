"""
page.json should contain the following information:
- CID: string, course ID
- URL: string, link
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

# Read page data
with open("page.json", "r") as f:
	page_data = json.load(f)

# Fetch article
cid = page_data["CID"]
link = page_data["URL"]
article = g.extract(url = link)
priority = page_data.get("Priority", 0)

# Reinsert newlines
content_arr = str(article.cleaned_text).split("\n")
content_arr = [x for x in content_arr if x != "" and not x.isspace()]
content = '<br>'.join(content_arr)

sql = "INSERT INTO page (cid, title, link, content, priority) VALUES (%s, %s, %s, %s, %s)"
val = (cid, article.title, link, content, priority)
cursor.execute(sql, val)
