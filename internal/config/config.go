package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"wakey/internal/common/wol"

	_ "modernc.org/sqlite"
)

type Device struct {
	ID          string
	DeviceName  string
	Description string
	MacAddress  string
	IPAddress   string
	State       string
}

type Group struct {
	ID        string
	GroupName string
	Devices   []string // contains IDs of devices
}

type Config struct {
	Devices []Device
	Groups  []Group
}

var (
	HomeDir, HomeDirErr = os.UserHomeDir()
	DBPath              string
	db                  *sql.DB
)

func init() {
	if HomeDirErr == nil {
		DBPath = filepath.Join(HomeDir, ".wakey.db")
	}
}

func getDB() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	if HomeDirErr != nil {
		return nil, fmt.Errorf("error getting home directory: %v", HomeDirErr)
	}
	conn, err := sql.Open("sqlite", DBPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	conn.SetMaxOpenConns(1)
	db = conn
	return db, nil
}

// ResetDB closes and clears the cached DB connection. Used in tests.
func ResetDB() {
	if db != nil {
		db.Close()
		db = nil
	}
}

func createTables(conn *sql.DB) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS devices (
			id          TEXT PRIMARY KEY,
			device_name TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			mac_address TEXT NOT NULL DEFAULT '',
			ip_address  TEXT NOT NULL DEFAULT '',
			state       TEXT NOT NULL DEFAULT 'Offline'
		);
		CREATE TABLE IF NOT EXISTS groups (
			id         TEXT PRIMARY KEY,
			group_name TEXT NOT NULL DEFAULT ''
		);
		CREATE TABLE IF NOT EXISTS group_devices (
			group_id  TEXT NOT NULL,
			device_id TEXT NOT NULL,
			PRIMARY KEY (group_id, device_id)
		);
	`)
	return err
}

// CreateConfig initializes the SQLite database
func CreateConfig() error {
	if HomeDirErr != nil {
		return fmt.Errorf("error getting home directory: %v", HomeDirErr)
	}

	conn, err := getDB()
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	if err := createTables(conn); err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	return fmt.Errorf("database ready at: %v", DBPath)
}

// ReadConfig reads all devices and groups from the database.
func ReadConfig() Config {
	conn, err := getDB()
	if err != nil {
		fmt.Println("error opening database:", err)
		return Config{}
	}

	devices, err := readDevices(conn)
	if err != nil {
		fmt.Println("error reading devices:", err)
		return Config{}
	}

	groups, err := readGroups(conn)
	if err != nil {
		fmt.Println("error reading groups:", err)
		return Config{Devices: devices}
	}

	return Config{Devices: devices, Groups: groups}
}

func readDevices(conn *sql.DB) ([]Device, error) {
	rows, err := conn.Query(`SELECT id, device_name, description, mac_address, ip_address, state FROM devices`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var d Device
		if err := rows.Scan(&d.ID, &d.DeviceName, &d.Description, &d.MacAddress, &d.IPAddress, &d.State); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	if devices == nil {
		devices = []Device{}
	}
	return devices, nil
}

func readGroups(conn *sql.DB) ([]Group, error) {
	rows, err := conn.Query(`
		SELECT g.id, g.group_name, gd.device_id
		FROM groups g
		LEFT JOIN group_devices gd ON g.id = gd.group_id
		ORDER BY g.rowid
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groupMap := make(map[string]*Group)
	var order []string
	for rows.Next() {
		var groupID, groupName string
		var deviceID sql.NullString
		if err := rows.Scan(&groupID, &groupName, &deviceID); err != nil {
			return nil, err
		}
		if _, exists := groupMap[groupID]; !exists {
			groupMap[groupID] = &Group{ID: groupID, GroupName: groupName, Devices: []string{}}
			order = append(order, groupID)
		}
		if deviceID.Valid {
			groupMap[groupID].Devices = append(groupMap[groupID].Devices, deviceID.String)
		}
	}

	groups := make([]Group, len(order))
	for i, id := range order {
		groups[i] = *groupMap[id]
	}
	if groups == nil {
		groups = []Group{}
	}
	return groups, nil
}

// WriteConfig replaces all devices and groups in the database with the given config.
func WriteConfig(config Config) {
	conn, err := getDB()
	if err != nil {
		fmt.Println("error opening database:", err)
		return
	}

	tx, err := conn.Begin()
	if err != nil {
		fmt.Println("error starting transaction:", err)
		return
	}
	defer tx.Rollback()

	tx.Exec(`DELETE FROM group_devices`)
	tx.Exec(`DELETE FROM groups`)
	tx.Exec(`DELETE FROM devices`)

	for _, d := range config.Devices {
		if _, err := tx.Exec(
			`INSERT INTO devices (id, device_name, description, mac_address, ip_address, state) VALUES (?, ?, ?, ?, ?, ?)`,
			d.ID, d.DeviceName, d.Description, d.MacAddress, d.IPAddress, d.State,
		); err != nil {
			fmt.Println("error inserting device:", err)
			return
		}
	}

	for _, g := range config.Groups {
		if _, err := tx.Exec(`INSERT INTO groups (id, group_name) VALUES (?, ?)`, g.ID, g.GroupName); err != nil {
			fmt.Println("error inserting group:", err)
			return
		}
		for _, deviceID := range g.Devices {
			if _, err := tx.Exec(`INSERT INTO group_devices (group_id, device_id) VALUES (?, ?)`, g.ID, deviceID); err != nil {
				fmt.Println("error inserting group device:", err)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("error committing transaction:", err)
	}
}

// GetUpdateState pings each device and updates its state in the database.
func GetUpdateState() Config {
	cfg := ReadConfig()
	devices := cfg.Devices

	for i, device := range devices {
		if wol.IsOnline(device.IPAddress) {
			devices[i].State = "Online"
		} else {
			devices[i].State = "Offline"
		}
	}

	WriteConfig(Config{Devices: devices, Groups: cfg.Groups})
	return Config{Devices: devices, Groups: cfg.Groups}
}
