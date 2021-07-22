package couch

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kivik/couchdb/v3"
	kivik "github.com/go-kivik/kivik/v3"
)

type Documents struct {
	ID string `json:"_id"`
}

func NewClient(ctx context.Context, user string, pass string, url string) (*kivik.Client, error) {
	client, err := kivik.New("couch", url)
	if err != nil {
		errs := fmt.Sprintf("Error creating Couch Client: %w\n", err.Error())
		return nil, errs
	}

	err = client.Authenticate(ctx, couchdb.BasicAuth(user, pass))
	if err != nil {
		fmt.Printf("Error Authenticating against couch: %w\n", err.Error())
		return nil, err
	}

	return client, nil
}

func Devices(ctx context.Context, client *kivik.Client) (devices, error) {
	var documents []string

	devices := client.DB(ctx, "devices") // connect to the devices database on couch

	fmt.Printf("Devices Output: %s\n", devices.Name())

	// Building Query for finding all the VIAs in the database
	query := map[string]interface{}{
		"fields": []string{"_id"},
		"limit":  10,
		"selector": map[string]interface{}{
			"_id": map[string]interface{}{
				"$regex": "VIA",
			},
		},
	}

	stats, err := devices.Stats(ctx)
	if err != nil {
		fmt.Printf("Stats Failed: %w\n", err.Error())
	}

	fmt.Printf("Running Stats: %v\n", stats.DocCount)

	rows, err := devices.Find(ctx, query)
	if err != nil {
		err_out := fmt.Sprintf("Sorry, I failed: %w", err.Error())
		return nil, err_out
	}
	fmt.Printf("What is the length: %v\n", rows.TotalRows)
	for rows.Next() {
		var docs Documents
		if err := rows.ScanDoc(&docs); err != nil {
			fmt.Printf("Error Scanning: %w", err.Error())
		}
		fmt.Printf("All Docs: %v\n", docs)
		fmt.Printf("It's type: %T\n", docs)
		documents = append(documents, docs.ID)
	}

	fmt.Printf("Output of Documents: %v", documents)

	return documents, nil
}
