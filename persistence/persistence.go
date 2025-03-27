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

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS google_secrets (id INTEGER PRIMARY KEY, client_id TEXT, client_secret TEXT, redirect_uri TEXT);")
	if err != nil {
		log.Fatal("Failed to create google_secrets table")
	}

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS profiles (id INTEGER PRIMARY KEY, email_address TEXT, secrets_id INTEGER, last_updated_at DATE, FOREIGN KEY (secrets_id) REFERENCES google_secrets(id));")
	if err != nil {
		log.Fatal("Failed to create profiles table")
	}

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS tokens (id INTEGER PRIMARY KEY, access_token TEXT, refresh_token TEXT, scope TEXT, token_type TEXT, id_token TEXT, expires_in NUMBER, expires TIMESTAMP, last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  profile_id INTEGER, FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE);")
	if err != nil {
		log.Fatal("Failed to create tokens table")
	}
}

func InsertProfile(profile *google.GoogleProfile) error {
	r, err := DB.Exec("INSERT INTO profiles (email_address, secrets_id, last_updated_at) VALUES (?, ?, ?)", profile.EmailAddress, profile.Secrets.ID, time.Now())
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
      expires,
      last_updated,
      profile_id
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		profile.Tokens.AccessToken,
		profile.Tokens.RefreshToken,
		profile.Tokens.Scope,
		profile.Tokens.TokenType,
		profile.Tokens.IdToken,
		profile.Tokens.ExpiresIn,
		time.Now().Add(time.Duration(profile.Tokens.ExpiresIn)*time.Second),
		time.Now(),
		lastInsertedId)

	return err
}

func UpdateProfile(profile *google.GoogleProfile, profileId int) error {
	_, err := DB.Exec("UPDATE profiles SET email_address = ?, secrets_id = ?, last_updated_at = ? where id = ?", profile.EmailAddress, profile.Secrets.ID, time.Now(), profileId)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
    UPDATE 'tokens' SET access_token = ?, refresh_token = ?, scope = ?, token_type = ?, id_token = ?, expires_in = ?, expires = ?, last_updated = ? WHERE profile_id = ?`,
		profile.Tokens.AccessToken,
		profile.Tokens.RefreshToken,
		profile.Tokens.Scope,
		profile.Tokens.TokenType,
		profile.Tokens.IdToken,
		profile.Tokens.ExpiresIn,
		time.Now().Add(time.Duration(profile.Tokens.ExpiresIn)*time.Second),
		time.Now(),
		profileId)

	if err != nil {
		return err
	}

	_, err = DB.Exec(`
    UPDATE 'google_secrets' SET client_id = ?, client_secret = ?, redirect_uri = ? WHERE id = ?`,
		profile.Secrets.ClientId,
		profile.Secrets.ClientSecret,
		profile.Secrets.RedirectUri,
		profile.Secrets.ID,
	)

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

func GetProfile(profileId int, profile *google.GoogleProfile) error {
	var r *sql.Row
	if profileId > 0 {
		r = DB.QueryRow("SELECT email_address, secrets_id, last_updated_at From profiles where Id = ?", profileId)
	} else {
		r = DB.QueryRow("SELECT email_address, secrets_id, last_updated_at From profiles")
	}

	err := r.Scan(&profile.EmailAddress, &profile.Secrets.ID, &profile.LastUpdatedAt)
	if err != nil {
		return err
	}

	r = DB.QueryRow("SELECT access_token, refresh_token, scope, token_type, id_token, expires_in, expires, last_updated From 'tokens' where profile_id = ?", profileId)
	err = r.Scan(&profile.Tokens.AccessToken, &profile.Tokens.RefreshToken, &profile.Tokens.Scope, &profile.Tokens.TokenType, &profile.Tokens.IdToken, &profile.Tokens.ExpiresIn, &profile.Tokens.Expires, &profile.Tokens.LastUpdated)
	if err != nil {
		return err
	}

	r = DB.QueryRow("SELECT client_id, client_secret, redirect_uri From google_secrets where id = ?", profile.Secrets.ID)
	err = r.Scan(&profile.Secrets.ClientId, &profile.Secrets.ClientSecret, &profile.Secrets.RedirectUri)

	return err
}

func GetProfiles(profiles *[]google.GoogleProfile) error {
	rows, err := DB.Query(`
    SELECT email_address, last_updated_at, access_token, refresh_token, scope, token_type, id_token, expires_in, expires, last_updated, clientId, clientSecret, redirect_uri 
      FROM profiles p
      JOIN tokens t on t.profile_id = p.id
      JOIN google_secrets g on p.secrets_id = g.id
    `)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		profile := google.GoogleProfile{}
		rows.Scan(&profile.EmailAddress, &profile.LastUpdatedAt, &profile.Tokens.AccessToken, &profile.Tokens.RefreshToken, &profile.Tokens.Scope, &profile.Tokens.TokenType, &profile.Tokens.IdToken, &profile.Tokens.ExpiresIn, &profile.Tokens.Expires, &profile.Tokens.LastUpdated, &profile.Secrets.ClientId, &profile.Secrets.ClientSecret, &profile.Secrets.RedirectUri)

		*profiles = append(*profiles, profile)
	}

	return nil
}

func RemoveProfile(id int) error {
	_, err := DB.Exec("DELETE FROM profiles WHERE ID = ?", id)
	return err
}

func UpsertGoogleSecrets(secret *google.GoogleSecret) error {
	r := DB.QueryRow("SELECT id from google_secrets")
	secretId := 0
	r.Scan(&secretId)
	var err error
	if secretId == 0 {
		err = InsertSecret(secret)
	} else {
		err = UpdateSecret(secretId, secret)
	}

	return err
}

func GetSecret(secretId int, secret *google.GoogleSecret) error {
	r := DB.QueryRow("SELECT id, client_id, client_secret, redirect_uri From google_secrets WHERE ID = ?", secretId)
	r.Scan(&secret.ID, &secret.ClientId, &secret.ClientSecret, &secret.RedirectUri)

	return nil
}

func InsertSecret(secret *google.GoogleSecret) error {
	_, err := DB.Exec("INSERT INTO google_secrets (client_id, client_secret, redirect_uri) VALUES (?, ?, ?)", secret.ClientId, secret.ClientSecret, secret.RedirectUri)
	return err
}

func UpdateSecret(secretId int, secret *google.GoogleSecret) error {
	_, err := DB.Exec("UPDATE google_secrets SET client_id = ?, client_secret = ?, redirect_uri = ? where id = ?", secret.ClientId, secret.ClientSecret, secret.RedirectUri, secretId)
	return err
}

func ClearDb() {
	DB.Exec("DELETE FROM profiles")
	DB.Exec("DELETE FROM google_secrets")
	DB.Exec("DELETE FROM tokens")
	DB.Exec("DROP TABLE profiles")
	DB.Exec("DROP TABLE google_secrets")
	DB.Exec("DROP TABLE tokens")
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
