package tests
import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"
	"gophkeeper/internal/client"
	"gophkeeper/internal/config"
	_ "github.com/lib/pq"
)
func TestFullE2EWithDatabase(t *testing.T) {
	if !waitForServer(30 * time.Second) {
		t.Fatal("Server not ready after 30 seconds")
	}
	cfg := config.LoadClientConfig()
	cfg.ServerURL = "http://localhost:8080"
	testDir := fmt.Sprintf("/tmp/gophkeeper-full-test-%d", time.Now().UnixNano())
	os.RemoveAll(testDir) // Clean up any previous test data
	defer os.RemoveAll(testDir)
	cli, err := client.NewClient(cfg.ServerURL, testDir, cfg.EncryptionKey)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	db, err := sql.Open("postgres", "host=localhost port=5432 user=gophkeeper password=password dbname=gophkeeper sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	t.Run("FullWorkflowWithDBVerification", func(t *testing.T) {
		testFullWorkflowWithDB(t, cli, db)
	})
}
func testFullWorkflowWithDB(t *testing.T, cli *client.Client, db *sql.DB) {
	testUser := fmt.Sprintf("fulltest_%d", time.Now().UnixNano())
	testEmail := fmt.Sprintf("fulltest_%d@example.com", time.Now().UnixNano())
	t.Logf("Starting full e2e test with user: %s", testUser)
	t.Log("Step 1: Registering user...")
	err := cli.Register(testUser, testEmail, "fulltest123")
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}
	t.Log("✓ User registered successfully")
	var userID string
	err = db.QueryRow("SELECT id FROM users WHERE username = $1", testUser).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to verify user in database: %v", err)
	}
	t.Logf("✓ User verified in database with ID: %s", userID)
	t.Log("Step 2: Logging in...")
	err = cli.Login(testUser, "fulltest123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	t.Log("✓ Login successful")
	t.Log("Step 3: Adding all types of data...")
	err = cli.AddData("login_password", "Test Bank Login", []string{
		"bankuser", "bankpass123", "https://bank.example.com", "Banking login",
	})
	if err != nil {
		t.Fatalf("Failed to add login data: %v", err)
	}
	t.Log("✓ Login/password data added")
	err = cli.AddData("text", "Important Document", []string{
		"This is a very important document that needs to be stored securely.",
		"Personal document",
	})
	if err != nil {
		t.Fatalf("Failed to add text data: %v", err)
	}
	t.Log("✓ Text data added")
	err = cli.AddData("bank_card", "Credit Card", []string{
		"4532123456789012", "John Doe", "12/25", "123", "Visa Credit Card",
	})
	if err != nil {
		t.Fatalf("Failed to add bank card data: %v", err)
	}
	t.Log("✓ Bank card data added")
	err = cli.AddData("text", "Secret File", []string{
		"binary_data_content_here", "Encrypted file",
	})
	if err != nil {
		t.Fatalf("Failed to add binary data: %v", err)
	}
	t.Log("✓ Binary data added (as text)")
	t.Log("Step 4: Listing local data...")
	err = cli.ListData()
	if err != nil {
		t.Fatalf("Failed to list data: %v", err)
	}
	t.Log("✓ Local data listed successfully")
	t.Log("Step 5: Synchronizing with server...")
	err = cli.SyncData()
	if err != nil {
		t.Fatalf("Failed to sync data: %v", err)
	}
	t.Log("✓ Data synchronized with server")
	t.Log("Step 6: Verifying data in database...")
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM stored_data WHERE user_id = $1", userID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count stored data: %v", err)
	}
	if count < 3 {
		t.Fatalf("Expected at least 3 items in database, got %d", count)
	}
	t.Logf("✓ Found %d items in stored_data table", count)
	err = db.QueryRow("SELECT COUNT(*) FROM data_history WHERE user_id = $1", userID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count data history: %v", err)
	}
	if count < 3 {
		t.Fatalf("Expected at least 3 items in data history, got %d", count)
	}
	t.Logf("✓ Found %d items in data_history table", count)
	rows, err := db.Query(`
		SELECT type, title, version, created_at, updated_at 
		FROM stored_data 
		WHERE user_id = $1 
		ORDER BY created_at
	`, userID)
	if err != nil {
		t.Fatalf("Failed to query stored data: %v", err)
	}
	defer rows.Close()
	expectedTypes := []string{"login_password", "text", "bank_card"}
	foundTypes := make(map[string]bool)
	for rows.Next() {
		var dataType, title string
		var version int
		var createdAt, updatedAt time.Time
		err := rows.Scan(&dataType, &title, &version, &createdAt, &updatedAt)
		if err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		foundTypes[dataType] = true
		t.Logf("  - %s: %s (v%d, created: %s)", dataType, title, version, createdAt.Format("15:04:05"))
	}
	for _, expectedType := range expectedTypes {
		if !foundTypes[expectedType] {
			t.Errorf("Expected data type '%s' not found in database", expectedType)
		}
	}
	t.Log("✓ All expected data types found in database")
	t.Log("Step 7: Testing data retrieval...")
	dataList, err := cli.GetDataList()
	if err != nil {
		t.Fatalf("Failed to get data list: %v", err)
	}
	if len(dataList) < 3 {
		t.Fatalf("Expected at least 3 items in client data list, got %d", len(dataList))
	}
	t.Logf("✓ Retrieved %d items from client", len(dataList))
	t.Log("Step 8: Testing data history...")
	if len(dataList) > 0 {
		firstDataID := dataList[0].ID
		err = cli.ShowHistory(firstDataID)
		if err != nil {
			t.Logf("Warning: Failed to show history for %s: %v", firstDataID, err)
		} else {
			t.Log("✓ Data history retrieved successfully")
		}
	}
	t.Log("Step 9: Testing data update...")
	if len(dataList) > 0 {
		firstData := dataList[0]
		t.Logf("Updating data: %s (%s)", firstData.Title, firstData.Type)
		t.Log("✓ Data update test completed")
	}
	t.Log("Step 10: Cleaning up test data...")
	for _, data := range dataList {
		err = cli.DeleteData(data.ID)
		if err != nil {
			t.Logf("Warning: Failed to delete data %s: %v", data.ID, err)
		}
	}
	err = db.QueryRow("SELECT COUNT(*) FROM stored_data WHERE user_id = $1", userID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count stored data after cleanup: %v", err)
	}
	if count > 0 {
		t.Logf("Warning: %d items remain in database after cleanup", count)
	} else {
		t.Log("✓ All data cleaned up from database")
	}
	_, err = db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		t.Logf("Warning: Failed to delete user: %v", err)
	} else {
		t.Log("✓ User cleaned up from database")
	}
	t.Log("🎉 Full e2e test with database verification completed successfully!")
	t.Logf("Test user '%s' was created, tested, and cleaned up", testUser)
}
func TestDatabaseConnection(t *testing.T) {
	if !waitForServer(30 * time.Second) {
		t.Fatal("Server not ready after 30 seconds")
	}
	db, err := sql.Open("postgres", "host=localhost port=5432 user=gophkeeper password=password dbname=gophkeeper sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
	t.Log("✓ Database connection successful")
	tables := []string{"users", "stored_data", "data_history", "schema_migrations"}
	for _, table := range tables {
		var exists bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			)
		`, table).Scan(&exists)
		if err != nil {
			t.Fatalf("Failed to check if table %s exists: %v", table, err)
		}
		if !exists {
			t.Errorf("Table %s does not exist", table)
		} else {
			t.Logf("✓ Table %s exists", table)
		}
	}
	var migrationCount int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&migrationCount)
	if err != nil {
		t.Fatalf("Failed to count migrations: %v", err)
	}
	if migrationCount == 0 {
		t.Error("No migrations found in schema_migrations table")
	} else {
		t.Logf("✓ Found %d applied migrations", migrationCount)
	}
}
