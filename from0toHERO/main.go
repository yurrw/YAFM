package main

import (
	"/pkg/http/rest/db"
	"/pkg/http/rest/router"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")
	r := router.New()
	v1 := r.Group("/api")

	//inicia o gorm
	d := db.New()
	db.AutoMigrate(d)

	// // Serving static content from web - we will populate this from within the docker container

	// dbUrl := os.Getenv("DATABASE_URL")
	// log.Printf("DB [%s]", dbUrl)
	// db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Fatalf("Error opening database: %q", err)
	// }
	// log.Println("booyah")
	// api.GET("/ping", pingFunc2(db))

	r.Run(":7777")
}

// func pingFunc2(db *sql.DB) gin.HandlerFunc {

// 	return func(c *gin.Context) {

// 		defer registerPing(db)
// 		r := db.QueryRow("SELECT occurred FROM ping_timestamp ORDER BY id DESC LIMIT 1")
// 		var lastDate pq.NullTime
// 		r.Scan(&lastDate)

// 		message := "first time!"
// 		if lastDate.Valid {
// 			message = fmt.Sprintf("%v ago", time.Now().Sub(lastDate.Time).String())
// 		}

// 		c.JSON(200, gin.H{
// 			"message": message,
// 		})
// 	}
// }
