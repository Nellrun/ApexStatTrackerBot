package postgresql

import (
	"database/sql"
	"database/sql/driver"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Subscription represent row from subscriptions
type Subscription struct {
	Username string
	Platform string
	ChatID   int64
}

// Subscribe add row from subscriptions base
func Subscribe(username string, chatId int64) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO subscriptions (username, chat_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;`

	_, err = db.Exec(query, username, chatId)
	if err != nil {
		return err
	}

	return nil
}

// Unsubscribe delete row from subscriptions base
func Unsubscribe(username string, chatID int64) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM subscriptions WHERE username = $1 and chat_id = $2;`

	_, err = db.Exec(query, username, chatID)
	if err != nil {
		return err
	}

	return nil
}

// AddImage upload image
func AddImage(imageTag string, image string) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO images (image_tag, image) VALUES ($1, $2) ON CONFLICT (image_tag) DO UPDATE SET image = $2;`

	_, err = db.Exec(query, imageTag, image)
	if err != nil {
		return err
	}

	return nil
}

// GetImage query
func GetImage(imageTag string) (*string, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT image from images WHERE image_tag = $1;`

	row := db.QueryRow(query, imageTag)

	image := new(string)

	err = row.Scan(&image)

	if err != nil {
		return nil, err
	}

	return image, nil
}

// GetSubscriptionsToSend get users to send statistic
func GetSubscriptionsToSend() ([]Subscription, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT username, platform, chat_id from subscriptions LIMIT 30;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Subscription
	for rows.Next() {
		subscription := new(Subscription)
		err := rows.Scan(&subscription.Username, &subscription.Platform, &subscription.ChatID)
		if err != nil {
			return nil, err
		}
		result = append(result, *subscription)
	}

	return result, nil
}

// GetSubscriptions func
func GetSubscriptions() ([]Subscription, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT DISTINCT(username), platform, chat_id from subscriptions LIMIT 30;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Subscription
	for rows.Next() {
		subscription := new(Subscription)
		err := rows.Scan(&subscription.Username, &subscription.Platform, &subscription.ChatID)
		if err != nil {
			return nil, err
		}
		result = append(result, *subscription)
	}

	return result, nil
}

// AddLog insert into log
func AddLog(username string, stats driver.Value, created time.Time) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO log (username, stats, created) VALUES ($1, $2, $3) ON CONFLICT (username, created) DO UPDATE SET stats = $2, updated = NOW();`

	_, err = db.Exec(query, username, stats, created)
	if err != nil {
		return err
	}

	return nil
}

// GetLog select user stats
func GetLog(username string, date time.Time) (string, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return "", err
	}
	defer db.Close()

	query := `SELECT stats from log WHERE username = $1 and created = $2;`

	row := db.QueryRow(query, username, date)

	var stats string

	err = row.Scan(&stats)

	if err != nil {
		return "", err
	}

	return stats, nil
}

// DeleteImage delete row from images
func DeleteImage(imageTag string) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM images WHERE image_tag = $1;`

	_, err = db.Exec(query, imageTag)
	if err != nil {
		return err
	}

	return nil
}

// CreateTables Initialization
func CreateTables() error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS subscriptions (id SERIAL PRIMARY KEY, username TEXT, platform TEXT DEFAULT 'psn', chat_id INT);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS s_username ON subscriptions (username, chat_id);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS images (image_tag TEXT PRIMARY KEY, image TEXT);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS log (
			id 		 SERIAL 		PRIMARY KEY,
			username TEXT			NOT NULL,
			stats 	 jsonb			NOT NULL,
			updated  TIMESTAMPTZ 	DEFAULT NOW(),
			created  TIMESTAMPTZ 	NOT NULL
		);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS log_username_created ON log (username, created);`)
	if err != nil {
		return err
	}

	return nil
}
