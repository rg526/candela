package cdengine

import (
	"strings"
	"regexp"
	"net/http"
	"github.com/gin-gonic/gin"

	"candela/cdmodel"
)

// Endpoint "/search"
// Search for courses
func GetSearch(ctx *gin.Context, ectx *Context) {
	// Verify token
	_, isAuth := VerifyTokenFromCtx(ctx, ectx)
	if !isAuth {
		return
	}

	// Find course ID
	var courseArr []cdmodel.Course

	query := ctx.Query("query")

	// Prepare query string
	queryStr := `SELECT DISTINCT
					course.cid, course.name, course.description,
					course.dept, course.units,
					course.prereq, course.coreq
			FROM course
			LEFT JOIN fce
			ON course.cid = fce.cid
			LEFT JOIN page
			ON course.cid = page.cid
			LEFT JOIN prof
			ON course.cid = prof.cid
			WHERE
				(LOCATE(?, course.cid) > 0 OR LOCATE(?, course.name) > 0 OR
					LOCATE(?, course.description) > 0 OR LOCATE(?, course.dept) > 0 OR
					LOCATE(?, page.content) > 0 OR
					LOCATE(?, prof.name) > 0) `
	var args []interface{}
	for i := 0;i < 6;i++ {
		args = append(args, query)
	}
	if ctx.Query("is_advanced") == "true" {
		// Level
		if ctx.Query("level") != "" {
			levelArr := strings.Split(ctx.Query("level"), ";")
			queryStr += " AND ( "
			for index, level := range levelArr {
				if index != 0 {
					queryStr += " OR "
				}
				queryStr += " fce.level = ? "
				args = append(args, level)
			}
			queryStr += " ) "
		}
		// Dept
		if ctx.Query("dept") != "" {
			deptArr := strings.Split(ctx.Query("dept"), ";")
			queryStr += " AND ( "
			for index, dept := range deptArr {
				if index != 0 {
					queryStr += " OR "
				}
				queryStr += " course.dept = ? "
				args = append(args, dept)
			}
			queryStr += " ) "
		}
		// Units
		if ctx.Query("units") != "" {
			unitsArr := strings.Split(ctx.Query("units"), ";")
			queryStr += " AND ( "
			for index, units := range unitsArr {
				if index != 0 {
					queryStr += " OR "
				}
				match2 := regexp.MustCompile(`^([0-9]+)\-([0-9]+)$`).
					FindSubmatch([]byte(units))
				if match2 != nil {
					queryStr += ` (course.units >= ? AND
						course.units <= ?) `
					args = append(args, match2[1])
					args = append(args, match2[2])
					continue
				}
				match1 := regexp.MustCompile(`^([0-9]+)\+$`).
					FindSubmatch([]byte(units))
				if match1 != nil {
					queryStr += ` (course.units >= ?) `
					args = append(args, match1[1])
					continue
				}
				ReportErrorFromString(ctx, http.StatusBadRequest,
					"Invalid units: " + units)
				return
			}
			queryStr += " ) "
		}
	}
	queryStr += ` LIMIT ? `
	args = append(args, ectx.Conf.MaxSearchResult)

	// Query DB
	rows, err := ectx.DB.Query(queryStr, args...)
	if err != nil {
		ReportError(ctx, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	// Append course to array
	for rows.Next() {
		var course cdmodel.Course
		err = rows.Scan(&course.CID, &course.Name, &course.Description, &course.Dept, &course.Units, &course.Prereq, &course.Coreq)
		if err != nil {
			ReportError(ctx, http.StatusInternalServerError, err)
			return
		}
		courseArr = append(courseArr, course)
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"Data": courseArr})
}
