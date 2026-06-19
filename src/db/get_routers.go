package db

import (
	"database/sql"
	"mikctl/src/config"
	"mikctl/src/models"
	"strings"
)

func GetRouters(mainOnly bool) ([]models.Router, error) {
	var where []string
	var args []any
	if mainOnly {
		where = append(where, "rf.fMain = 1")
	} else {
		where = append(where, "rf.fMain = ?")
		args = append(args, config.Main)
	}

	query := `
		SELECT rf.id,
			rf.name,
			rf.ip,
			rf.ssh_user,
			rf.site,
			rf.device_id,
			rf.fMain
		FROM routers_full rf
		`
	if config.Group != "" {
		query += `
			JOIN router_groups rg ON rg.router_id = rf.id
		`
		where = append(where, "rg.name = ?")
		args = append(args, config.Group)
	}

	if config.Router != "" {
		where = append(where, "rf.name = ?")
		args = append(args, config.Router)
	}

	if config.Site != "" {
		where = append(where, "rf.site = ?")
		args = append(args, config.Site)
	}

	if config.ROS != "" {
		where = append(where, "rf.ros_version LIKE ?")
		args = append(args, config.ROS+".%")
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY rf.name"

	rows, err := DB.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var routers []models.Router
	for rows.Next() {
		var r models.Router
		err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.IP,
			&r.SSHUser,
			&r.Site,
			&r.DeviceID,
			&r.Main,
		)
		if err != nil {
			return nil, err
		}
		routers = append(
			routers,
			r,
		)
	}
	return routers, nil

}

func IsMainRouter(name string) (bool, error) {
	var main int

	err := DB.QueryRow(`
		SELECT fMain
		FROM routers
		WHERE name = ?
		LIMIT 1
	`, name).Scan(&main)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return main == 1, nil
}
