package couch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kivik/couchdb/v3"
	kivik "github.com/go-kivik/kivik/v3"
)

// type QueryString (QS) creating map so that a query parameter can be passed
type QS map[string]interface{}

type Documents struct {
	ID string `json:"_id"`
}

func NewClient(ctx context.Context, user string, pass string, url string) (*kivik.Client, error) {
	shortDuration := 10 * time.Second
	d := time.Now().Add(shortDuration)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	client, err := kivik.New("couch", url)
	if err != nil {
		errs := fmt.Sprintf("Error creating Couch Client: %w\n", err.Error())
		return nil, errors.New(errs)
	}

	err = client.Authenticate(ctx, couchdb.BasicAuth(user, pass))
	if err != nil {
		fmt.Printf("Error Authenticating against couch: %w\n", err.Error())
		return nil, err
	}

	return client, nil
}

func CouchQuery(ctx context.Context, client *kivik.Client, query QS, database string) ([]string, error) {
	var documents []string

	devices := client.DB(ctx, database) // connect to the database of your choice on couch

	fmt.Printf("Devices Output: %s\n", devices.Name())

	stats, err := devices.Stats(ctx)
	if err != nil {
		fmt.Printf("Stats Failed: %w\n", err.Error())
	}

	fmt.Printf("Running Stats: %v\n", stats.DocCount)

	rows, err := devices.Find(ctx, query)
	if err != nil {
		err_out := fmt.Sprintf("Sorry, I failed: %w", err.Error())
		return nil, errors.New(err_out)
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

	fmt.Printf("Output of Documents: %v\n", documents)

	return documents, nil
}
