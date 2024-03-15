package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/amp-labs/connectors"
	"github.com/amp-labs/connectors/providers"
	"github.com/amp-labs/connectors/proxy"
	"github.com/amp-labs/connectors/salesforce"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// Set the appropriate environment variables in a .env file, then run:
// go run test/salesforce.go

const TimeoutSeconds = 30

func main() {
	os.Exit(mainFn())
}

func mainFn() int { //nolint:funlen
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)

		return 1
	}

	salesforceWorkspace := os.Getenv("SALESFORCE_WORKSPACE")
	clientId := os.Getenv("SALESFORCE_CLIENT_ID")
	clientSecret := os.Getenv("SALESFORCE_CLIENT_SECRET")
	accessToken := os.Getenv("SALESFORCE_ACCESS_TOKEN")
	refreshToken := os.Getenv("SALESFORCE_REFRESH_TOKEN")

	workspace := flag.String("workspace", salesforceWorkspace, "Salesforce workspace")
	flag.Parse()

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	cfg := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   fmt.Sprintf("https://%s.my.salesforce.com/services/oauth2/authorize", *workspace),
			TokenURL:  fmt.Sprintf("https://%s.my.salesforce.com/services/oauth2/token", *workspace),
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	tok := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
		Expiry:       time.Now().Add(-1 * time.Hour), // just pretend it's expired already, whatever, it'll fetch a new one.
	}

	ctx := context.Background()

	// Create a new Salesforce connector, with a token provider that uses the sfdx CLI to fetch an access token.
	proxyConn, err := connectors.NewProxyConnector(
		providers.Salesforce,
		proxy.WithClient(ctx, http.DefaultClient, cfg, tok),
		proxy.WithCatalogSubstitutions(map[string]string{
			salesforce.PlaceholderWorkspace: *workspace,
		}),
	)

	sfc, err := salesforce.NewConnector(salesforce.WithProxyConnector(proxyConn))
	if err != nil {
		slog.Error("Error creating Salesforce connector", "error", err)

		return 1
	}

	if err := testConnector(ctx, sfc); err != nil {
		slog.Error("Error testing", "connector", sfc, "error", err)

		return 1
	}

	// IMPORTANT: every time this test is run, it will create a new Account
	// in SFDC instance. Will need to delete those out at later date.
	writtenRecordId, err := testSalesforceValidCreate(ctx, sfc)
	if err != nil {
		slog.Error("Error creating record in Salesforce", "error", err)

		return 1
	}
	// IMPORTANT: will fail if specific recordId does not already exist in instance
	if err := testSalesforceValidUpdate(ctx, sfc, writtenRecordId); err != nil {
		slog.Error("Error updating record in Salesforce", "error", err)

		return 1
	}

	return 0
}

func testConnector(ctx context.Context, conn connectors.ReadConnector) error {
	// Create a context with a timeout
	ctx, done := context.WithTimeout(ctx, TimeoutSeconds*time.Second)
	defer done()

	// Read some data from Salesforce
	res, err := conn.Read(ctx, connectors.ReadParams{
		ObjectName: "Account",
		Fields:     []string{"Id", "Name", "BillingCity", "IsDeleted"},
	})
	if err != nil {
		return fmt.Errorf("error reading from Salesforce: %w", err)
	}

	// Print the results
	jsonStr, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	_, _ = os.Stdout.Write(jsonStr)
	_, _ = os.Stdout.WriteString("\n")

	return nil
}

const accountNumber = 123

// testSalesforceValidCreate will create a valid record in Salesforce.
func testSalesforceValidCreate(ctx context.Context, conn connectors.WriteConnector) (string, error) {
	writeRes, err := conn.Write(ctx, connectors.WriteParams{
		ObjectName: "Account",
		RecordData: map[string]interface{}{
			"Name":          "TEST ACCOUNT - [TO DELETE]",
			"AccountNumber": accountNumber,
		},
	})
	if err != nil {
		return "", fmt.Errorf("error writing to Salesforce: %w", err)
	}

	// Print the results
	jsonStr, err := json.MarshalIndent(writeRes, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %w", err)
	}

	_, _ = os.Stdout.Write(jsonStr)
	_, _ = os.Stdout.WriteString("\n")

	return writeRes.RecordId, nil
}

const accountNumber2 = 456

// testSalesforceValidUpdate will update existing record in Salesforce.
func testSalesforceValidUpdate(ctx context.Context, conn connectors.WriteConnector, writtenRecordId string) error {
	writeRes, err := conn.Write(ctx, connectors.WriteParams{
		ObjectName: "Account",
		RecordData: map[string]interface{}{
			"Name":          "OKADA TEST ACCOUNT",
			"AccountNumber": accountNumber2,
		},
		RecordId: writtenRecordId,
	})
	if err != nil {
		return fmt.Errorf("error writing to Salesforce: %w", err)
	}

	if !writeRes.Success {
		return fmt.Errorf("write to %s failed when it should have succeeded", writtenRecordId) //nolint:goerr113
	}

	// Print the results
	jsonStr, err := json.MarshalIndent(writeRes, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	_, _ = os.Stdout.Write(jsonStr)
	_, _ = os.Stdout.WriteString("\n")

	return nil
}
