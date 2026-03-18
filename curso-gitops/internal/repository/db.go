package repository

import (
	"curso-gitops/internal/models"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func ConnectDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	return DB.Ping()
}

func GetUserRole(username, password string) (string, error) {
	var role, hashedPassword string
	err := DB.QueryRow("SELECT password, role FROM users WHERE username=?", username).Scan(&hashedPassword, &role)
	if err != nil {
		return "", fmt.Errorf("credenciales inválidas")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return "", fmt.Errorf("credenciales inválidas")
	}
	return role, nil
}

func CreateUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al procesar contraseña")
	}
	_, err = DB.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, 'student')", username, string(hashedPassword))
	return err
}

func GetModules(isAdmin bool) ([]models.Module, error) {
	query := "SELECT id, title, is_hidden FROM modules"
	if !isAdmin {
		query += " WHERE is_hidden=FALSE"
	}
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []models.Module
	for rows.Next() {
		var m models.Module
		if err := rows.Scan(&m.ID, &m.Title, &m.IsHidden); err != nil {
			continue
		}

		epQuery := "SELECT id, title, video_url, is_hidden FROM episodes WHERE module_id=?"
		if !isAdmin {
			epQuery += " AND is_hidden=FALSE"
		}
		epRows, err := DB.Query(epQuery, m.ID)
		if err != nil {
			modules = append(modules, m)
			continue
		}
		for epRows.Next() {
			var ep models.Episode
			if err := epRows.Scan(&ep.ID, &ep.Title, &ep.VideoURL, &ep.IsHidden); err != nil {
				continue
			}
			m.Episodes = append(m.Episodes, ep)
		}
		epRows.Close()
		modules = append(modules, m)
	}
	return modules, nil
}

func AddModule(title string) error {
	_, err := DB.Exec("INSERT INTO modules (title) VALUES (?)", title)
	return err
}

func DeleteModule(id string) error {
	_, err := DB.Exec("DELETE FROM modules WHERE id=?", id)
	return err
}

func ToggleModule(id string) error {
	_, err := DB.Exec("UPDATE modules SET is_hidden = NOT is_hidden WHERE id=?", id)
	return err
}

func UpdateModule(id string, title string) error {
	_, err := DB.Exec("UPDATE modules SET title=? WHERE id=?", title, id)
	return err
}

func AddEpisode(modID int, title, url string) error {
	_, err := DB.Exec("INSERT INTO episodes (module_id, title, video_url) VALUES (?, ?, ?)", modID, title, url)
	return err
}

func DeleteEpisode(id string) error {
	_, err := DB.Exec("DELETE FROM episodes WHERE id=?", id)
	return err
}

func UpdateEpisode(id string, title, url string) error {
	_, err := DB.Exec("UPDATE episodes SET title=?, video_url=? WHERE id=?", title, url, id)
	return err
}
