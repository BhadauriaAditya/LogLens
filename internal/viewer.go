package internal

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const logDir = "./logs"

// ViewLogs handles GET /logs
func ViewLogs(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(logDir)
	if err != nil {
		http.Error(w, "Failed to read log directory", http.StatusInternalServerError)
		return
	}

	// Collect only .log files
	var logFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".log") {
			logFiles = append(logFiles, f.Name())
		}
	}

	// Sort (latest first)
	sort.Sort(sort.Reverse(sort.StringSlice(logFiles)))

	// If ?file=YYYY-MM-DD.log is provided, show file content
	selected := r.URL.Query().Get("file")
	if selected != "" {
		filePath := filepath.Join(logDir, selected)
		data, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Failed to read log file", http.StatusInternalServerError)
			return
		}

		tmpl := `
		<html>
		<head>
			<title>Log Viewer</title>
			<meta http-equiv="refresh" content="5">
			<style>
				body { background: #111; color: #0f0; font-family: monospace; }
				pre { background: #000; padding: 15px; border-radius: 8px; overflow-x: auto; }
				a { color: #08f; text-decoration: none; }
				.header { margin-bottom: 10px; }
			</style>
		</head>
		<body>
			<div class="header">
				<h2>Viewing Log: {{.File}}</h2>
				<a href="/logs">Back to log list</a>
			</div>
			<pre>{{.Content}}</pre>
		</body>
		</html>`
		t, _ := template.New("view").Parse(tmpl)
		t.Execute(w, map[string]interface{}{
			"File":    selected,
			"Content": string(data),
		})
		return
	}

	// Else show list of available logs
	tmpl := `
	<html>
	<head>
		<title>Log Viewer</title>
		<style>
			body { font-family: sans-serif; background: #f4f4f4; }
			ul { list-style: none; padding: 0; }
			li { margin: 8px 0; }
			a { color: #007bff; text-decoration: none; font-weight: bold; }
			a:hover { text-decoration: underline; }
		</style>
	</head>
	<body>
		<h2>Available Logs</h2>
		<ul>
			{{range .}}
				<li><a href="/logs?file={{.}}">{{.}}</a></li>
			{{end}}
		</ul>
	</body>
	</html>`
	t, _ := template.New("list").Parse(tmpl)
	t.Execute(w, logFiles)
}
