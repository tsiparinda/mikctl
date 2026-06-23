package db

import (
	"fmt"
	"mikctl/src/config"
	"mikctl/src/models"
	"strings"
)

func SaveRouterInfo(router *models.Router) error {

	router.Serial = strings.TrimSpace(router.Serial)

	var serialValue any
	if router.Serial == "" {
		serialValue = nil
	} else {
		serialValue = router.Serial
	}

	config.Verbosef("DEVICE model=%q serial=%q version=%q\n", router.Model, router.Serial, router.ROSVersion)

	if serialValue != nil {
		_, err := DB.Exec(`
		INSERT INTO devices(
			model,
			serial_number,
			ros_version
		)
		VALUES (?, ?, ?)
		ON CONFLICT(serial_number)
		DO UPDATE SET
			model         = excluded.model,
			ros_version   = excluded.ros_version
		`,
			router.Model,
			serialValue,
			router.ROSVersion,
		)
		if err != nil {
			return fmt.Errorf("devices %v", err)
		}
		err = DB.QueryRow(`
        SELECT id
        FROM devices
        WHERE serial_number = ?
    `,
			serialValue,
		).Scan(&router.DeviceID)

		if err != nil {
			return err
		}

	}

	var parentID any
	var deviceID any
	if router.ParentID == 0 {
		parentID = nil
	} else {
		parentID = router.ParentID
	}
	if router.DeviceID == 0 {
		deviceID = nil
	} else {
		deviceID = router.DeviceID
	}

	config.Verbosef("DEVICE device_id=%q name=%q fMain=%q parent_router_id=%q\n", deviceID, strings.TrimSpace(router.Name), router.Main, parentID)

	_, err := DB.Exec(`
		INSERT INTO routers(
			device_id,
			name,
			fMain,
			parent_router_id,
			ssh_user,
			password_id,
			last_seen_at
		)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(name)
		DO UPDATE SET
			name = excluded.name,
			fMain = excluded.fMain,
			parent_router_id = excluded.parent_router_id,
			last_seen_at = CURRENT_TIMESTAMP,
			device_id=COALESCE(excluded.device_id, routers.device_id)
		`,
		deviceID,
		strings.TrimSpace(router.Name),
		router.Main,
		parentID,
		router.SSHUser,
		router.PasswordID,
	)
	if err != nil {
		return fmt.Errorf("routers %v", err)
	}
	//get ID
	err = DB.QueryRow(`
    SELECT id
    FROM routers
    WHERE name = ?
`,
		strings.TrimSpace(router.Name),
	).Scan(&router.ID)

	if err != nil {
		return err
	}
	config.Verbosef("router.ID=%d name=%s\n", router.ID, router.Name)
	//update/insert address
	_, err = DB.Exec(`
    UPDATE routers_addresses
    SET fMain = 0
    WHERE router_id = ?
`,
		router.ID,
	)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
    INSERT INTO routers_addresses(
        router_id,
        ip,
        fMain
    )
    VALUES (?, ?, 1)
    ON CONFLICT(router_id, ip)
    DO UPDATE SET
        fMain = 1
`,
		router.ID,
		router.IP,
	)
	if err != nil {
		return err
	}
	config.Verbosef("%s updated\n", router.Name)
	return nil
}
