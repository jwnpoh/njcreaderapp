package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type Config struct {
	pscaleDSN    string
	cockroachDSN string
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("unable to load config: %v", err)
	}
	fmt.Println("config loaded.")

	fmt.Println("connecting databases...")
	pscaleDB, cockroachDB := openDBConnections(cfg.pscaleDSN, cfg.cockroachDSN)
	defer pscaleDB.Close()
	defer cockroachDB.Close()

	// fmt.Println("migrating articles from pscale to cockroach...")
	// err = migrateArticlesTable(pscaleDB, cockroachDB)
	// if err != nil {
	// 	fmt.Printf("unable to complete articles table migration: %v\n", err)
	// }

	// fmt.Println("migrating users from pscale to cockroach...")
	// err = migrateUsersTable(pscaleDB, cockroachDB)
	// if err != nil {
	// 	fmt.Printf("unable to complete users table migration: %v\n", err)
	// }

	// fmt.Println("migrating longs from pscale to cockroach...")
	// err = migrateLongTable(pscaleDB, cockroachDB)
	// if err != nil {
	// 	fmt.Printf("unable to complete users table migration: %v\n", err)
	// }

	// fmt.Println("migrating notes from pscale to cockroach...")
	// err = migrateNotesTable(pscaleDB, cockroachDB)
	// if err != nil {
	// 	fmt.Printf("unable to complete users table migration: %v\n", err)
	// }

	fmt.Println("migrating questions list from pscale to cockroach...")
	err = migrateQuestionListTable(pscaleDB, cockroachDB)
	if err != nil {
		fmt.Printf("unable to complete users table migration: %v\n", err)
	}
}

func loadConfig() (Config, error) {
	godotenv.Load(".env")

	cfg := Config{
		pscaleDSN:    os.Getenv("PSCALE"),
		cockroachDSN: os.Getenv("COCKROACH"),
	}

	return cfg, nil
}

func openDBConnections(pscaleDSN, cockroachDSN string) (pscaleDB, cockroachDB *sqlx.DB) {

	fmt.Println("opening pscale connection...")
	pscaleDB, err := sqlx.Open("mysql", pscaleDSN)
	if err != nil {
		fmt.Printf("ArticlesDB: unable to initialize pscale database - %v\n", err)
		return
	}
	fmt.Println("pinging pscale...")
	err = pscaleDB.Ping()
	if err != nil {
		fmt.Printf("PScaleArticles: no response from pscale database - %v\n", err)
		return
	}
	pscaleDB.SetMaxIdleConns(10)
	pscaleDB.SetMaxOpenConns(50)

	fmt.Println("opening cockroach connection...")
	cockroachDB, err = sqlx.Open("pgx", cockroachDSN)
	if err != nil {
		fmt.Printf("ArticlesDB: unable to initialize cockroach database - %v\n", err)
		return
	}
	fmt.Println("pinging cockroach...")
	err = cockroachDB.Ping()
	if err != nil {
		fmt.Printf("PScaleArticles: no response from pscale database - %v\n", err)
		return
	}
	cockroachDB.SetMaxIdleConns(10)
	cockroachDB.SetMaxOpenConns(50)

	return pscaleDB, cockroachDB
}
