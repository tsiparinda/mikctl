package db

import (
	"mikctl/src/models"
	"time"
)

func SaveScriptResult(
	r models.ScriptResult,
) error {

	_, err := DB.Exec(`
        INSERT INTO script_runs(
            router_id,
            script_name,
            started_at,
            finished_at,
            status,
            message,
			tag
        )
        VALUES(?,?,?,?,?,?,?)
    `,
		r.Router.ID,
		r.Script,
		r.StartedAt.Format(time.DateTime),
		r.FinishedAt.Format(time.DateTime),
		boolToStatus(r.Success),
		r.Output,
		r.Tag,
	)

	return err
}

func boolToStatus(status bool) string {
	if status {
		return "OK"
	}
	return "ERROR"
}
