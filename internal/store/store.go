package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"
	_ "modernc.org/sqlite"
)

type DB struct { db *sql.DB }

type Services struct {
	ID string `json:"id"`
	ServiceName string `json:"service_name"`
	DurationMinutes int64 `json:"duration_minutes"`
	Price float64 `json:"price"`
	Description string `json:"description"`
	Active bool `json:"active"`
	CreatedAt string `json:"created_at"`
}

type Appointments struct {
	ID string `json:"id"`
	ClientName string `json:"client_name"`
	ClientEmail string `json:"client_email"`
	ClientPhone string `json:"client_phone"`
	Service string `json:"service"`
	Date string `json:"date"`
	Time string `json:"time"`
	Status string `json:"status"`
	Notes string `json:"notes"`
	CreatedAt string `json:"created_at"`
}

type Availability struct {
	ID string `json:"id"`
	DayOfWeek string `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
	Active bool `json:"active"`
	CreatedAt string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil { return nil, err }
	db, err := sql.Open("sqlite", filepath.Join(d, "booking.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil { return nil, err }
	db.SetMaxOpenConns(1)
	db.Exec(`CREATE TABLE IF NOT EXISTS services(id TEXT PRIMARY KEY, service_name TEXT NOT NULL, duration_minutes INTEGER DEFAULT 0, price REAL DEFAULT 0, description TEXT DEFAULT '', active INTEGER DEFAULT 0, created_at TEXT DEFAULT(datetime('now')))`)
	db.Exec(`CREATE TABLE IF NOT EXISTS appointments(id TEXT PRIMARY KEY, client_name TEXT NOT NULL, client_email TEXT DEFAULT '', client_phone TEXT DEFAULT '', service TEXT DEFAULT '', date TEXT NOT NULL, time TEXT NOT NULL, status TEXT DEFAULT '', notes TEXT DEFAULT '', created_at TEXT DEFAULT(datetime('now')))`)
	db.Exec(`CREATE TABLE IF NOT EXISTS availability(id TEXT PRIMARY KEY, day_of_week TEXT NOT NULL, start_time TEXT NOT NULL, end_time TEXT NOT NULL, active INTEGER DEFAULT 0, created_at TEXT DEFAULT(datetime('now')))`)
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }
func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string { return time.Now().UTC().Format(time.RFC3339) }

func (d *DB) CreateServices(e *Services) error {
	e.ID = genID(); e.CreatedAt = now()
	_, err := d.db.Exec(`INSERT INTO services(id, service_name, duration_minutes, price, description, active, created_at) VALUES(?, ?, ?, ?, ?, ?, ?)`, e.ID, e.ServiceName, e.DurationMinutes, e.Price, e.Description, e.Active, e.CreatedAt)
	return err
}

func (d *DB) GetServices(id string) *Services {
	var e Services
	if d.db.QueryRow(`SELECT id, service_name, duration_minutes, price, description, active, created_at FROM services WHERE id=?`, id).Scan(&e.ID, &e.ServiceName, &e.DurationMinutes, &e.Price, &e.Description, &e.Active, &e.CreatedAt) != nil { return nil }
	return &e
}

func (d *DB) ListServices() []Services {
	rows, _ := d.db.Query(`SELECT id, service_name, duration_minutes, price, description, active, created_at FROM services ORDER BY created_at DESC`)
	if rows == nil { return nil }; defer rows.Close()
	var o []Services
	for rows.Next() { var e Services; rows.Scan(&e.ID, &e.ServiceName, &e.DurationMinutes, &e.Price, &e.Description, &e.Active, &e.CreatedAt); o = append(o, e) }
	return o
}

func (d *DB) UpdateServices(e *Services) error {
	_, err := d.db.Exec(`UPDATE services SET service_name=?, duration_minutes=?, price=?, description=?, active=? WHERE id=?`, e.ServiceName, e.DurationMinutes, e.Price, e.Description, e.Active, e.ID)
	return err
}

func (d *DB) DeleteServices(id string) error {
	_, err := d.db.Exec(`DELETE FROM services WHERE id=?`, id)
	return err
}

func (d *DB) CountServices() int {
	var n int; d.db.QueryRow(`SELECT COUNT(*) FROM services`).Scan(&n); return n
}

func (d *DB) SearchServices(q string, filters map[string]string) []Services {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (service_name LIKE ? OR description LIKE ?)"
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
	}
	rows, _ := d.db.Query(`SELECT id, service_name, duration_minutes, price, description, active, created_at FROM services WHERE `+where+` ORDER BY created_at DESC`, args...)
	if rows == nil { return nil }; defer rows.Close()
	var o []Services
	for rows.Next() { var e Services; rows.Scan(&e.ID, &e.ServiceName, &e.DurationMinutes, &e.Price, &e.Description, &e.Active, &e.CreatedAt); o = append(o, e) }
	return o
}

func (d *DB) CreateAppointments(e *Appointments) error {
	e.ID = genID(); e.CreatedAt = now()
	_, err := d.db.Exec(`INSERT INTO appointments(id, client_name, client_email, client_phone, service, date, time, status, notes, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, e.ID, e.ClientName, e.ClientEmail, e.ClientPhone, e.Service, e.Date, e.Time, e.Status, e.Notes, e.CreatedAt)
	return err
}

func (d *DB) GetAppointments(id string) *Appointments {
	var e Appointments
	if d.db.QueryRow(`SELECT id, client_name, client_email, client_phone, service, date, time, status, notes, created_at FROM appointments WHERE id=?`, id).Scan(&e.ID, &e.ClientName, &e.ClientEmail, &e.ClientPhone, &e.Service, &e.Date, &e.Time, &e.Status, &e.Notes, &e.CreatedAt) != nil { return nil }
	return &e
}

func (d *DB) ListAppointments() []Appointments {
	rows, _ := d.db.Query(`SELECT id, client_name, client_email, client_phone, service, date, time, status, notes, created_at FROM appointments ORDER BY created_at DESC`)
	if rows == nil { return nil }; defer rows.Close()
	var o []Appointments
	for rows.Next() { var e Appointments; rows.Scan(&e.ID, &e.ClientName, &e.ClientEmail, &e.ClientPhone, &e.Service, &e.Date, &e.Time, &e.Status, &e.Notes, &e.CreatedAt); o = append(o, e) }
	return o
}

func (d *DB) UpdateAppointments(e *Appointments) error {
	_, err := d.db.Exec(`UPDATE appointments SET client_name=?, client_email=?, client_phone=?, service=?, date=?, time=?, status=?, notes=? WHERE id=?`, e.ClientName, e.ClientEmail, e.ClientPhone, e.Service, e.Date, e.Time, e.Status, e.Notes, e.ID)
	return err
}

func (d *DB) DeleteAppointments(id string) error {
	_, err := d.db.Exec(`DELETE FROM appointments WHERE id=?`, id)
	return err
}

func (d *DB) CountAppointments() int {
	var n int; d.db.QueryRow(`SELECT COUNT(*) FROM appointments`).Scan(&n); return n
}

func (d *DB) SearchAppointments(q string, filters map[string]string) []Appointments {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (client_name LIKE ? OR client_email LIKE ? OR client_phone LIKE ? OR service LIKE ? OR time LIKE ? OR notes LIKE ?)"
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
	}
	if v, ok := filters["status"]; ok && v != "" { where += " AND status=?"; args = append(args, v) }
	rows, _ := d.db.Query(`SELECT id, client_name, client_email, client_phone, service, date, time, status, notes, created_at FROM appointments WHERE `+where+` ORDER BY created_at DESC`, args...)
	if rows == nil { return nil }; defer rows.Close()
	var o []Appointments
	for rows.Next() { var e Appointments; rows.Scan(&e.ID, &e.ClientName, &e.ClientEmail, &e.ClientPhone, &e.Service, &e.Date, &e.Time, &e.Status, &e.Notes, &e.CreatedAt); o = append(o, e) }
	return o
}

func (d *DB) CreateAvailability(e *Availability) error {
	e.ID = genID(); e.CreatedAt = now()
	_, err := d.db.Exec(`INSERT INTO availability(id, day_of_week, start_time, end_time, active, created_at) VALUES(?, ?, ?, ?, ?, ?)`, e.ID, e.DayOfWeek, e.StartTime, e.EndTime, e.Active, e.CreatedAt)
	return err
}

func (d *DB) GetAvailability(id string) *Availability {
	var e Availability
	if d.db.QueryRow(`SELECT id, day_of_week, start_time, end_time, active, created_at FROM availability WHERE id=?`, id).Scan(&e.ID, &e.DayOfWeek, &e.StartTime, &e.EndTime, &e.Active, &e.CreatedAt) != nil { return nil }
	return &e
}

func (d *DB) ListAvailability() []Availability {
	rows, _ := d.db.Query(`SELECT id, day_of_week, start_time, end_time, active, created_at FROM availability ORDER BY created_at DESC`)
	if rows == nil { return nil }; defer rows.Close()
	var o []Availability
	for rows.Next() { var e Availability; rows.Scan(&e.ID, &e.DayOfWeek, &e.StartTime, &e.EndTime, &e.Active, &e.CreatedAt); o = append(o, e) }
	return o
}

func (d *DB) UpdateAvailability(e *Availability) error {
	_, err := d.db.Exec(`UPDATE availability SET day_of_week=?, start_time=?, end_time=?, active=? WHERE id=?`, e.DayOfWeek, e.StartTime, e.EndTime, e.Active, e.ID)
	return err
}

func (d *DB) DeleteAvailability(id string) error {
	_, err := d.db.Exec(`DELETE FROM availability WHERE id=?`, id)
	return err
}

func (d *DB) CountAvailability() int {
	var n int; d.db.QueryRow(`SELECT COUNT(*) FROM availability`).Scan(&n); return n
}

func (d *DB) SearchAvailability(q string, filters map[string]string) []Availability {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (start_time LIKE ? OR end_time LIKE ?)"
		args = append(args, "%"+q+"%")
		args = append(args, "%"+q+"%")
	}
	if v, ok := filters["day_of_week"]; ok && v != "" { where += " AND day_of_week=?"; args = append(args, v) }
	rows, _ := d.db.Query(`SELECT id, day_of_week, start_time, end_time, active, created_at FROM availability WHERE `+where+` ORDER BY created_at DESC`, args...)
	if rows == nil { return nil }; defer rows.Close()
	var o []Availability
	for rows.Next() { var e Availability; rows.Scan(&e.ID, &e.DayOfWeek, &e.StartTime, &e.EndTime, &e.Active, &e.CreatedAt); o = append(o, e) }
	return o
}
