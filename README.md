# Project Candela

## Introduction

[Project Candela](https://cmucourselist.org) (PROJCD) is a centralized knowledgebase for courses at CMU. You can search for courses, view/post comments on a course, and find related pages (such as course syllabus) on this site.

Project Candela contains the following modules:

- CDSource: Data sources
- CDModel: Data structure definitions, including database table schemes and Go structs
- CDAPI: Interface to CDEngine (JSON)
- CDEngine: Engine for PROJCD - "logic element" for PROJCD
- CDSite: A website interface to CDEngine (https://cmucourselist.org)
- CDFrontend: Frontend things

PROJCD is designed to be modular. You can write your own client to CDEngine (Android, command-line, GUI client, browser extension, etc.) using CDAPI. 

## Getting Started

To setup PROJCD,

1. Install MySQL, Golang, Python

2. Setup a database and a user

3. Go to cdmodel and execute all sql scripts

4. In cdsource, create `config.py` 
    ```
    DB_HOST = ...
    DB_USER = ...
    DB_PWD = ...
    DB_NAME = ...
    ```

5. Download `fce.csv` from FCE

6. Run `course_source.py`, `fce_source.py` and `rmp_source.py` to genertae `course_data.json`, `fce_data.json` and `rmp_data.json`

7. Run `course_import.py`, `fce_import.py` and `rmp_import.py`

8. In cdengine/route, create `config.json`
    ```json
    {
    	"Host": ...,
    	"Port": ...,
    	"DBUser": ...,
    	"DBPwd": ...,
    	"DBName": ...,
    	"DBHost": ...,
    	"OAuth2ClientID": ...,
    	"OAuth2ClientSecret": ...,
    	"OAuth2Scope": ...,
    	"OAuth2RedirectURI": ...,
    	"MaxSearchResult": ...
    }
    ```

    You may need to go to Google Cloud Console to setup OAuth2

9. Run `route.go` (you may need to `go get`)

10. In cdsite/route, create `config.json`
    ```json
    {
    	"Host": ...,
    	"Port": ...,
    	"CDAPIUrl": ...,
    	"CookieSecret": ...,
    	"OAuth2ClientID": ...,
    	"OAuth2ClientSecret": ...,
    	"OAuth2Scope": ...,
    	"OAuth2RedirectURI": ...
    }
    ```

11. Run `route.go` (you may need to `go get`)

12. To add/delete pages, in cdsource, create `page.json`
    ```json
    {
    	"CID": ...,
    	"URL": ...,
    	"Priority": ...
    }
    ```

    where priority is optional - it specifies the priority of that page
    and then run `page_add.py` or `page_delete.py`

13. Setup your webserver to point to the CDSite address/port

## Credits

This project uses the following third-party source code:

- [Bootswatch](https://bootswatch.com/) (both [Litera](https://bootswatch.com/litera/) and [Darkly](https://bootswatch.com/darkly/) themes)
