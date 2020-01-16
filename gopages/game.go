package gopages

import (
	"bytes"
	"database/sql"
	"html/template"
	"net/http"
)

func createGameInfo(db *sql.DB) *gameInfo {
	var (
		id          int
		displayName string
		name        string
		names       []string
		events      []string
	)
	// get user info
	rows, err := db.Query("select id, name, display_name from players;")
	if err != nil {
		// do nothing
	} else {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name, &displayName)
			if err != nil {
				displayName = ""
			}
			names = append(names, "\""+displayName+" ("+name+")\",")
		}
		names = append(names, "\"New Player\"")
	}

	// get event type info
	rows, err = db.Query("select * from event_types;")
	if err != nil {
		// do nothing
	} else {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name)
			if err == nil {
				events = append(events, "\""+name+"\",")
			} else {
				// load page with error message
			}
		}
		events = append(events, "\"New Type\"")
	}
	return &gameInfo{Players: names, GoalTypes: events}
}

type gameInfo struct {
	Players   []string
	GoalTypes []string
}

// RenderGamePage : renders the game input form page with correct data
func RenderGamePage(db *sql.DB, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	var g = createGameInfo(db)

	t, err := template.ParseFiles("webpage/game_template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return template.HTML(""), err
	}
	var buff bytes.Buffer
	if err = t.Execute(&buff, g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return template.HTML(""), err
	}
	return template.HTML(buff.String()), nil
}