package main

import (
	"context"
	"log"
	"receipt-processor/api/routes"
	"receipt-processor/ent"
)

func main() {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("failed closing connection to sqlite: %v", err)
		}
	}(client)

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	router := routes.SetupRoutes(client)
	err = router.Run(":8081")
	if err != nil {
		return
	}
}
