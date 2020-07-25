package postgresql

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

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
