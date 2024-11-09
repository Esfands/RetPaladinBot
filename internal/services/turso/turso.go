package turso

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/esfands/retpaladinbot/internal/db"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type SetupOptions struct {
	URL string
}

func Setup(ctx context.Context, opts SetupOptions) (Service, error) {
	svc := &tursoService{}
	var err error

	svc.db, err = sql.Open("libsql", opts.URL)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		return nil, err
	}

	slog.Info("Turso database connection opened")

	err = svc.db.Ping()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error pinging database: %v\n", err)
		return nil, err
	}

	slog.Info("Turso database connection pinged")

	svc.queries = db.NewQueries(svc.db)

	go func() {
		<-ctx.Done()
		err := svc.db.Close()
		if err != nil {
			return
		}
		slog.Info("Turso database connection closed")
	}()

	return svc, nil
}
