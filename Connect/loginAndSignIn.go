package connect

import (
	"database/sql"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func UserCookies(w http.ResponseWriter, id int) {
	cookie := http.Cookie{
		Name:  "Id",
		Value: strconv.Itoa(id),
		Path:  "/",
	}
	http.SetCookie(w, &cookie)
}

func Islogin(r *http.Request) (bool, int) {
	cookie, err := r.Cookie("Id")
	if err != nil {
		return false, 0
	}
	id, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return false, 0
	}
	return true, id
}

func HashPassword(password string) (string, error) {
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hasedPassword), nil
}

func QueryNameUsers(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT Username FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		usernames = append(usernames, name)
	}
	return usernames, nil
}

func QueryEmailUsers(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT Email FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		Emails = append(Emails, email)
	}
	return Emails, nil
}

func QueryPasswordUsers(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT Password FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var Passwords []string
	for rows.Next() {
		var password string
		if err := rows.Scan(&password); err != nil {
			return nil, err
		}
		Passwords = append(Passwords, password)
	}
	return Passwords, nil
}

func QueryUserId(usernameOrEmail string) (int, error) {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	query := "Select UserID FROM Users WHERE Username = ? OR Email = ?"
	var id int
	err = db.QueryRow(query, usernameOrEmail, usernameOrEmail).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		}
		return 0, err
	}
	return id, nil
}

type User struct {
	email    string
	password string
}

func GetAllUsersDetails(db *sql.DB) (map[string]User, error) {
	rows, err := db.Query("SELECT Username, Email, Password FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make(map[string]User)
	for rows.Next() {
		var datas struct{ Username, Email, Password string }
		if err := rows.Scan(&datas.Username, &datas.Email, &datas.Password); err != nil {
			return nil, err
		}
		users[datas.Username] = User{datas.Email, datas.Password}
	}
	return users, nil
}

func IsMatch(usernameOrEmail, password string, db *sql.DB) (bool, error) {
	users, err := GetAllUsersDetails(db)
	if err != nil {
		return false, err
	}
	for userName, userValue := range users {
		if usernameOrEmail == userName || usernameOrEmail == userValue.email {
			err := bcrypt.CompareHashAndPassword([]byte(userValue.password), []byte(password))
			if err == nil {
				return true, nil
			} else {
				return false, err
			}
		}
	}
	return false, nil
}
