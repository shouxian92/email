package driving

// ToHTML converts an array of Timeslots into a html table
func ToHTML(metadata []map[string]string) string {
	rows := ""
	for _, t := range metadata {
		rows += "<tr><td>" + t["StartTime"] + "</td><td>" + t["Date"] + "</td></tr>"
	}

	return `<table border="1"><tr><th>Time</th><th>Date</th></tr>` + rows + "</table>"
}
