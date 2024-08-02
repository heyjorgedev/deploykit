package view

import (
	"io"

	"github.com/heyjorgedev/deploykit"
)

type TeamListData struct {
	SelectedTeam *deploykit.Team
	Teams        []*deploykit.Team
}

func RenderTeamList(w io.Writer, data TeamListData) error {
	return htmlTemplates.ExecuteTemplate(w, "team_list.html", data)
}

func RenderTeamListItem(w io.Writer, data *deploykit.Team) error {
	return htmlTemplates.ExecuteTemplate(w, "team_list_item", data)
}
