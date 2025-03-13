package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/M1z23R/google-cli/google"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDb() {
	err := ensureConfigDir()
	dbPath, err := getConfigPath()

	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open db")
	}

	_, err = DB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("Failed to set foreign_keys")
	}

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS profiles (id INTEGER PRIMARY KEY, email_address TEXT, last_updated_at DATE);")
	if err != nil {
		log.Fatal("Failed to create profiles table")
	}

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS tokens (id INTEGER PRIMARY KEY, access_token TEXT, refresh_token TEXT, scope TEXT, token_type TEXT, id_token TEXT, expires_in NUMBER, profile_id INTEGER, FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE);")
	if err != nil {
		log.Fatal("Failed to create tokens table")
	}
}

func InsertProfile(profile *google.GoogleProfile) error {
	r, err := DB.Exec("INSERT INTO profiles (email_address, last_updated_at) VALUES (?, ?)", profile.EmailAddress, time.Now())
	if err != nil {
		return err
	}

	lastInsertedId, err := r.LastInsertId()
	if err != nil {
		return err
	}
	_, err = DB.Exec(`
    INSERT INTO 'tokens' (
      access_token,
      refresh_token,
      scope,
      token_type,
      id_token,
      expires_in,
      profile_id
    ) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		profile.Tokens.AccessToken,
		profile.Tokens.RefreshToken,
		profile.Tokens.Scope,
		profile.Tokens.TokenType,
		profile.Tokens.IdToken,
		profile.Tokens.ExpiresIn,
		lastInsertedId)

	return err
}

func UpdateProfile(profile *google.GoogleProfile, profileId int) error {
	_, err := DB.Exec("UPDATE profiles SET email_address = ?, last_updated_at = ? where id = ?", profile.EmailAddress, time.Now(), profileId)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
    UPDATE 'tokens' SET access_token = ?, refresh_token = ?, scope = ?, token_type = ?, id_token = ?, expires_in = ? WHERE profile_id = ?`,
		profile.Tokens.AccessToken,
		profile.Tokens.RefreshToken,
		profile.Tokens.Scope,
		profile.Tokens.TokenType,
		profile.Tokens.IdToken,
		profile.Tokens.ExpiresIn,
		profileId)

	return err
}

func UpsertProfile(profile *google.GoogleProfile) error {
	r := DB.QueryRow("SELECT id from profiles where email_address = ?", profile.EmailAddress)
	profileId := 0
	r.Scan(&profileId)
	var err error
	if profileId == 0 {
		err = InsertProfile(profile)
	} else {
		err = UpdateProfile(profile, profileId)
	}

	return err
}

func GetFirstProfile(profile *google.GoogleProfile) error {
	r := DB.QueryRow("SELECT email_address, last_updated_at From profiles")
	r.Scan(&profile.EmailAddress, &profile.LastUpdatedAt)

	r = DB.QueryRow("SELECT access_token, refresh_token, scope, token_type, id_token, expires_in From 'tokens' where profile_id = ?", profile.ID)
	r.Scan(&profile.Tokens.AccessToken, &profile.Tokens.RefreshToken, &profile.Tokens.Scope, &profile.Tokens.TokenType, &profile.Tokens.IdToken, &profile.Tokens.ExpiresIn)

	return nil
}

func GetProfile(profileId string, profile *google.GoogleProfile) error {
	r := DB.QueryRow("SELECT email_address, last_updated_at From profiles where Id = ?", profileId)
	r.Scan(&profile.EmailAddress, &profile.LastUpdatedAt)

	r = DB.QueryRow("SELECT access_token, refresh_token, scope, token_type, id_token, expires_in From 'tokens' where profile_id = ?", profileId)
	r.Scan(&profile.Tokens.AccessToken, &profile.Tokens.RefreshToken, &profile.Tokens.Scope, &profile.Tokens.TokenType, &profile.Tokens.IdToken, &profile.Tokens.ExpiresIn)

	return nil
}

func GetProfiles(profiles *[]google.GoogleProfile) error {
	rows, err := DB.Query(`
    SELECT email_address, last_updated_at, access_token, refresh_token, scope, token_type, id_token, expires_in
      FROM profiles p
      JOIN tokens t on t.profile_id = p.id
    `)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		profile := google.GoogleProfile{}
		rows.Scan(&profile.EmailAddress, &profile.LastUpdatedAt, &profile.Tokens.AccessToken, &profile.Tokens.RefreshToken, &profile.Tokens.Scope, &profile.Tokens.TokenType, &profile.Tokens.IdToken, &profile.Tokens.ExpiresIn)

		*profiles = append(*profiles, profile)
	}

	return nil
}

func ClearDb() error {
	_, err := DB.Exec("DELETE FROM profiles")
	return err
}

func RemoveProfile(id int) error {
	_, err := DB.Exec("DELETE FROM profiles WHERE ID = ?", id)
	return err
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %v", err)
	}
	return filepath.Join(home, ".config", "g-cli", "db.db"), nil
}

func ensureConfigDir() error {
	dbPath, err := getConfigPath()
	if err != nil {
		return err
	}
	configDir := filepath.Dir(dbPath)

	os.MkdirAll(configDir, 0755)
	return nil
}
