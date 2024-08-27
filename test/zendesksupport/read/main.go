package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/providers/zendesksupport"
	"github.com/amp-labs/connectors/test/utils"
	connTest "github.com/amp-labs/connectors/test/zendesksupport"
)

var objectName = "users" // nolint: gochecknoglobals

func main() {
	// Handle Ctrl-C gracefully.
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	// Set up slog logging.
	utils.SetupLogging()

	conn := connTest.GetZendeskSupportConnector(ctx)
	defer utils.Close(conn)

	res, err := conn.Read(ctx, common.ReadParams{
		ObjectName: objectName,
		Fields: []string{
			"name", "time_zone", "role",
		},
	})
	if err != nil {
		utils.Fail("error reading from Zendesk Support", "error", err)
	}

	fmt.Println("Reading users..")
	utils.DumpJSON(res, os.Stdout)

	if res.Rows > zendesksupport.DefaultPageSize {
		utils.Fail(fmt.Sprintf("expected max %v rows", zendesksupport.DefaultPageSize))
	}
}
