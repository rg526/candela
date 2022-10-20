var isAdvanced = false;
function disappear(){
	if (isAdvanced) {
    	document.getElementById("advance-element").style.display="none";
		document.getElementById("advance-mode").innerHTML = "Advanced Search"
		document.getElementById("advance-mode").style = "text-decoration:none;";
		isAdvanced = false;
	} else {
    	document.getElementById("advance-element").style.display = "block";
		document.getElementById("advance-mode").innerHTML = "Hide";
		document.getElementById("advance-mode").style = "text-decoration:underline;";
		isAdvanced = true;
  	}
}

// Construct search options
var searchOpt = {
	level: new Set(),
	dept: new Set(),
	units: new Set()
};


// Recover last query status
function setLastQueryStatus() {
	const params = (new URL(document.location)).searchParams;
	// is_advanced
	if (params.get("is_advanced") == "true") {
		disappear();
	}
	// query
	if (params.get("query") !== null) {
		document.getElementById("searchbox").value = params.get("query");
	}
	// level
	if (params.get("level") !== null) {
		const arr = params.get("level").split(";");
		const elemList = document.getElementsByClassName("level-item");
		for (let i = 0;i < elemList.length; i++) {
			if (arr.includes(elemList[i].innerHTML)) {
				searchOptActivate(elemList[i].id);
			}
		}
	}
	// dept
	if (params.get("dept") !== null) {
		const arr = params.get("dept").split(";");
		const elemList = document.getElementsByClassName("dept-item");
		for (let i = 0;i < elemList.length; i++) {
			if (arr.includes(elemList[i].innerHTML.replace(/&amp;/g, "&"))) {
				searchOptActivate(elemList[i].id);
			}
		}
	}
	// units
	if (params.get("units") !== null) {
		const arr = params.get("units").split(";");
		const elemList = document.getElementsByClassName("units-item");
		for (let i = 0;i < elemList.length; i++) {
			if (arr.includes(elemList[i].innerHTML)) {
				searchOptActivate(elemList[i].id);
			}
		}
	}
}
setLastQueryStatus()

function toggle(id){
	const elem = document.getElementById(id);
	if (elem.classList.contains("active")) {
		searchOptDeactivate(id);
	} else {
		searchOptActivate(id);
	}
}

function searchOptActivate(id) {
	const elem = document.getElementById(id);
	const name = elem.innerHTML;
	if (id.startsWith("level")) searchOpt.level.add(name);
	if (id.startsWith("dept")) searchOpt.dept.add(name.replace(/&amp;/g, "&"));
	if (id.startsWith("units")) searchOpt.units.add(name);
	elem.classList.add("active");
}
function searchOptDeactivate(id) {
	const elem = document.getElementById(id);
	const name = elem.innerHTML;
	if (id.startsWith("level")) searchOpt.level.delete(name);
	if (id.startsWith("dept")) searchOpt.dept.delete(name.replace(/&amp;/g, "&"));
	if (id.startsWith("units")) searchOpt.units.delete(name);
	elem.classList.remove("active");
}
function searchOptPreset(className) {
	const elemList = document.getElementsByClassName(className)
	for (let i = 0;i < elemList.length; i++) {
		searchOptActivate(elemList[i].id);
	}
}
function searchOptReset(className) {
	const elemList = document.getElementsByClassName(className)
	for (let i = 0;i < elemList.length; i++) {
		searchOptDeactivate(elemList[i].id);
	}
}

// Execute search
function doSearch() {
	const queryText = document.getElementById("searchbox").value;

	// Detect if searchbox contains a single course number
	if (!isAdvanced) {
		if (/^[0-9]{5}$/.test(queryText)) {
			window.location.assign("/course/" + queryText);
			return;
		} else if (/^([0-9]{2})\-([0-9]{3})$/.test(queryText)) {
			const courseID = queryText.slice(0, 2) + queryText.slice(3);
			window.location.assign("/course/" + courseID);
			return;
		}
	}

	// Fallback to searching
	searchPath = "/search?exec=true";
	const searchParams = {
		query: queryText,
		is_advanced: false
	};
	if (isAdvanced) {
		searchParams.is_advanced = true;
		Object.entries(searchOpt).forEach(([key, value]) => {
			searchParams[key] = Array.from(value).join(";");
		});
	}

	Object.entries(searchParams).forEach(([key, value]) => {
		searchPath += "&" + key + "=" + encodeURIComponent(value);
	});
	window.location.replace(searchPath);
}
