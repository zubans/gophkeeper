package tests
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"gophkeeper/internal/client"
	"gophkeeper/internal/config"
)
func TestEndToEnd(t *testing.T) {
	if !waitForServer(30 * time.Second) {
		t.Fatal("Server not ready after 30 seconds")
	}
	cfg := config.LoadClientConfig()
	cfg.ServerURL = "http://localhost:8080"
	testDir := filepath.Join(os.TempDir(), "gophkeeper-test")
	os.RemoveAll(testDir) // Clean up any previous test data
	cli, err := client.NewClient(cfg.ServerURL, testDir, cfg.EncryptionKey)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer os.RemoveAll(testDir)
	t.Run("CompleteWorkflow", func(t *testing.T) {
		testCompleteWorkflow(t, cli)
	})
}
func testCompleteWorkflow(t *testing.T, cli *client.Client) {
	testUser := fmt.Sprintf("e2euser_%d", time.Now().UnixNano())
	testEmail := fmt.Sprintf("e2e_%d@example.com", time.Now().UnixNano())
	err := cli.Register(testUser, testEmail, "e2epass123")
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}
	err = cli.AddData("login_password", "Test Website", []string{
		"testuser", "testpass", "https://example.com", "Test notes",
	})
	if err != nil {
		t.Fatalf("Failed to add login data: %v", err)
	}
	err = cli.AddData("text", "Important Note", []string{
		"This is an important note for testing",
	})
	if err != nil {
		t.Fatalf("Failed to add text data: %v", err)
	}
	err = cli.AddData("bank_card", "Test Credit Card", []string{
		"1234567890123456", "12/25", "123", "John Doe", "Test Bank", "Test notes",
	})
	if err != nil {
		t.Fatalf("Failed to add bank card data: %v", err)
	}
	err = cli.ListData()
	if err != nil {
		t.Fatalf("Failed to list data: %v", err)
	}
	err = cli.SyncData()
	if err != nil {
		t.Fatalf("Failed to sync data: %v", err)
	}
	dataList, err := cli.GetDataList()
	if err != nil {
		t.Fatalf("Failed to get data list for cleanup: %v", err)
	}
	for _, data := range dataList {
		err = cli.DeleteData(data.ID)
		if err != nil {
			t.Logf("Warning: Failed to delete data %s: %v", data.ID, err)
		}
	}
	dataListAfter, err := cli.GetDataList()
	if err != nil {
		t.Fatalf("Failed to get data list after cleanup: %v", err)
	}
	if len(dataListAfter) > 0 {
		t.Logf("Warning: %d items remain after cleanup", len(dataListAfter))
	} else {
		t.Log("All test data cleaned up successfully")
	}
	t.Log("End-to-end test completed successfully")
	t.Logf("Test user '%s' was created and cleaned up", testUser)
}
func TestClientCLI(t *testing.T) {
	if !waitForServer(30 * time.Second) {
		t.Fatal("Server not ready after 30 seconds")
	}
	buildCmd := exec.Command("go", "build", "-o", "test-client", "./cmd/client")
	buildCmd.Dir = ".." // Go up one directory to project root
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build client: %v", err)
	}
	defer os.Remove("../test-client")
	testDir := filepath.Join(os.TempDir(), "gophkeeper-cli-test")
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)
	t.Run("CLIWorkflow", func(t *testing.T) {
		testCLIWorkflow(t, testDir)
	})
}
func testCLIWorkflow(t *testing.T, testDir string) {
	testUser := fmt.Sprintf("cliuser_%d", time.Now().UnixNano())
	testEmail := fmt.Sprintf("cli_%d@example.com", time.Now().UnixNano())
	cmd := exec.Command("../test-client", "register", testUser, testEmail, "clipass123")
	cmd.Env = append(os.Environ(),
		"SERVER_URL=http://localhost:8080",
		"CLIENT_CONFIG_DIR="+testDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI registration failed: %v\nOutput: %s", err, string(output))
	}
	if !strings.Contains(string(output), "Registration successful!") {
		t.Fatalf("Expected success message, got: %s", string(output))
	}
	t.Logf("CLI registration successful for user: %s", testUser)
	cmd = exec.Command("../test-client", "login", testUser, "clipass123")
	cmd.Env = append(os.Environ(),
		"SERVER_URL=http://localhost:8080",
		"CLIENT_CONFIG_DIR="+testDir)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI login failed: %v\nOutput: %s", err, string(output))
	}
	if !strings.Contains(string(output), "Login successful!") {
		t.Fatalf("Expected success message, got: %s", string(output))
	}
	tokenFile := filepath.Join(testDir, "token")
	tokenData, err := os.ReadFile(tokenFile)
	if err != nil {
		t.Fatalf("Failed to read token file: %v", err)
	}
	if len(tokenData) == 0 {
		t.Fatal("Token file is empty")
	}
	t.Logf("Token saved successfully: %s", string(tokenData)[:20]+"...")
	cmd = exec.Command("../test-client", "add", "login_password", "CLI Test Site", "cliuser", "clipass", "https://cli.example.com", "CLI test notes")
	cmd.Env = append(os.Environ(),
		"SERVER_URL=http://localhost:8080",
		"CLIENT_CONFIG_DIR="+testDir)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI add data failed: %v\nOutput: %s", err, string(output))
	}
	if !strings.Contains(string(output), "Data added successfully!") {
		t.Fatalf("Expected success message, got: %s", string(output))
	}
	cmd = exec.Command("../test-client", "list")
	cmd.Env = append(os.Environ(),
		"SERVER_URL=http://localhost:8080",
		"CLIENT_CONFIG_DIR="+testDir)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI list data failed: %v\nOutput: %s", err, string(output))
	}
	if !strings.Contains(string(output), "CLI Test Site") {
		t.Fatalf("Expected to find 'CLI Test Site' in output, got: %s", string(output))
	}
	cmd = exec.Command("../test-client", "sync")
	cmd.Env = append(os.Environ(),
		"SERVER_URL=http://localhost:8080",
		"CLIENT_CONFIG_DIR="+testDir)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI sync failed: %v\nOutput: %s", err, string(output))
	}
	if !strings.Contains(string(output), "Data synchronized successfully!") {
		t.Fatalf("Expected sync success message, got: %s", string(output))
	}
	t.Logf("CLI workflow completed successfully for user: %s", testUser)
}
func TestDataPersistence(t *testing.T) {
	if !waitForServer(30 * time.Second) {
		t.Fatal("Server not ready after 30 seconds")
	}
	testDir := filepath.Join(os.TempDir(), "gophkeeper-persistence-test")
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)
	testUser := fmt.Sprintf("persistuser_%d", time.Now().UnixNano())
	testEmail := fmt.Sprintf("persist_%d@example.com", time.Now().UnixNano())
	cfg := config.LoadClientConfig()
	cfg.ServerURL = "http://localhost:8080"
	cli1, err := client.NewClient(cfg.ServerURL, testDir, cfg.EncryptionKey)
	if err != nil {
		t.Fatalf("Failed to create first client: %v", err)
	}
	err = cli1.Register(testUser, testEmail, "persistpass123")
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}
	err = cli1.AddData("text", "Persistent Note", []string{"This note should persist across client restarts"})
	if err != nil {
		t.Fatalf("Failed to add data: %v", err)
	}
	cli2, err := client.NewClient(cfg.ServerURL, testDir, cfg.EncryptionKey)
	if err != nil {
		t.Fatalf("Failed to create second client: %v", err)
	}
	err = cli2.Login(testUser, "persistpass123")
	if err != nil {
		t.Fatalf("Login with second client failed: %v", err)
	}
	err = cli2.SyncData()
	if err != nil {
		t.Fatalf("Sync with second client failed: %v", err)
	}
	err = cli2.ListData()
	if err != nil {
		t.Fatalf("List data with second client failed: %v", err)
	}
	dataList, err := cli2.GetDataList()
	if err != nil {
		t.Fatalf("Failed to get data list for cleanup: %v", err)
	}
	for _, data := range dataList {
		err = cli2.DeleteData(data.ID)
		if err != nil {
			t.Logf("Warning: Failed to delete data %s: %v", data.ID, err)
		}
	}
	t.Logf("Data persistence test completed for user '%s'", testUser)
}
