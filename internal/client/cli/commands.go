package cli
import (
	"fmt"
	"regexp"
	"gophkeeper/internal/models"
)
type ClientInterface interface {
	Register(username, email, password string) error
	Login(username, password string) error
	AddData(dataType, title string, data []string) error
	GetData(id string) error
	DeleteData(id string) error
	SyncData() error
	ShowHistory(id string) error
	ListData() error
	GetDataList() ([]models.StoredData, error)
}
type Command interface {
	Execute(client ClientInterface) error
}
type RegisterCommand struct {
	Username string
	Email    string
	Password string
}
func (c *RegisterCommand) Execute(client ClientInterface) error {
	if c.Username == "" {
		return fmt.Errorf("username is required")
	}
	if len(c.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if len(c.Username) > 50 {
		return fmt.Errorf("username must be no more than 50 characters long")
	}
	if c.Email == "" {
		return fmt.Errorf("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(c.Email) {
		return fmt.Errorf("invalid email format")
	}
	if c.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(c.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	return client.Register(c.Username, c.Email, c.Password)
}
type LoginCommand struct {
	Username string
	Password string
}
func (c *LoginCommand) Execute(client ClientInterface) error {
	if c.Username == "" {
		return fmt.Errorf("username is required")
	}
	if c.Password == "" {
		return fmt.Errorf("password is required")
	}
	return client.Login(c.Username, c.Password)
}
type AddCommand struct {
	DataType string
	Title    string
	Data     []string
}
func (c *AddCommand) Execute(client ClientInterface) error {
	validTypes := map[string]bool{"login_password": true, "text": true, "binary": true, "bank_card": true}
	if c.DataType == "" {
		return fmt.Errorf("data type is required")
	}
	if !validTypes[c.DataType] {
		return fmt.Errorf("invalid data type: %s. Valid types: login_password, text, binary, bank_card", c.DataType)
	}
	if c.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(c.Title) > 255 {
		return fmt.Errorf("title must be no more than 255 characters long")
	}
	switch c.DataType {
	case "login_password":
		if len(c.Data) < 2 {
			return fmt.Errorf("login_password requires at least username and password")
		}
		if len(c.Data) > 4 {
			return fmt.Errorf("login_password accepts at most: username, password, url, metadata")
		}
	case "text":
		if len(c.Data) < 1 {
			return fmt.Errorf("text requires at least the text content")
		}
		if len(c.Data) > 2 {
			return fmt.Errorf("text accepts at most: content, metadata")
		}
	case "binary":
		if len(c.Data) < 1 {
			return fmt.Errorf("binary requires at least the binary data")
		}
		if len(c.Data) > 2 {
			return fmt.Errorf("binary accepts at most: data, metadata")
		}
	case "bank_card":
		if len(c.Data) < 4 {
			return fmt.Errorf("bank_card requires: number, holder, expiry, cvv")
		}
		if len(c.Data) > 5 {
			return fmt.Errorf("bank_card accepts at most: number, holder, expiry, cvv, metadata")
		}
	}
	return client.AddData(c.DataType, c.Title, c.Data)
}
type GetCommand struct{ ID string }
func (c *GetCommand) Execute(client ClientInterface) error {
	if c.ID == "" {
		return fmt.Errorf("data ID is required")
	}
	if len(c.ID) < 1 {
		return fmt.Errorf("invalid data ID")
	}
	return client.GetData(c.ID)
}
type DeleteCommand struct{ ID string }
func (c *DeleteCommand) Execute(client ClientInterface) error {
	if c.ID == "" {
		return fmt.Errorf("data ID is required")
	}
	if len(c.ID) < 1 {
		return fmt.Errorf("invalid data ID")
	}
	return client.DeleteData(c.ID)
}
type SyncCommand struct{}
func (c *SyncCommand) Execute(client ClientInterface) error {
	return client.SyncData()
}
type HistoryCommand struct{ ID string }
func (c *HistoryCommand) Execute(client ClientInterface) error {
	if c.ID == "" {
		return fmt.Errorf("data ID is required")
	}
	if len(c.ID) < 1 {
		return fmt.Errorf("invalid data ID")
	}
	return client.ShowHistory(c.ID)
}
type ListCommand struct{}
func (c *ListCommand) Execute(client ClientInterface) error {
	return client.ListData()
}
type HelpCommand struct{}
func (c *HelpCommand) Execute(client ClientInterface) error {
	ShowHelp()
	return nil
}
type VersionCommand struct{}
func (c *VersionCommand) Execute(client ClientInterface) error {
	fmt.Println("GophKeeper Client")
	return nil
}
func ParseCommand(args []string) (Command, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no command specified")
	}
	command := args[0]
	commandArgs := args[1:]
	switch command {
	case "register":
		if len(commandArgs) != 3 {
			return nil, fmt.Errorf("register command requires exactly 3 arguments: username, email, password")
		}
		return &RegisterCommand{Username: commandArgs[0], Email: commandArgs[1], Password: commandArgs[2]}, nil
	case "login":
		if len(commandArgs) != 2 {
			return nil, fmt.Errorf("login command requires exactly 2 arguments: username, password")
		}
		return &LoginCommand{Username: commandArgs[0], Password: commandArgs[1]}, nil
	case "add":
		if len(commandArgs) < 2 {
			return nil, fmt.Errorf("add command requires at least 2 arguments: type, title")
		}
		return &AddCommand{DataType: commandArgs[0], Title: commandArgs[1], Data: commandArgs[2:]}, nil
	case "get":
		if len(commandArgs) != 1 {
			return nil, fmt.Errorf("get command requires exactly 1 argument: id")
		}
		return &GetCommand{ID: commandArgs[0]}, nil
	case "delete":
		if len(commandArgs) != 1 {
			return nil, fmt.Errorf("delete command requires exactly 1 argument: id")
		}
		return &DeleteCommand{ID: commandArgs[0]}, nil
	case "sync":
		if len(commandArgs) != 0 {
			return nil, fmt.Errorf("sync command takes no arguments")
		}
		return &SyncCommand{}, nil
	case "history":
		if len(commandArgs) != 1 {
			return nil, fmt.Errorf("history command requires exactly 1 argument: id")
		}
		return &HistoryCommand{ID: commandArgs[0]}, nil
	case "list":
		if len(commandArgs) != 0 {
			return nil, fmt.Errorf("list command takes no arguments")
		}
		return &ListCommand{}, nil
	case "help":
		if len(commandArgs) != 0 {
			return nil, fmt.Errorf("help command takes no arguments")
		}
		return &HelpCommand{}, nil
	case "version":
		if len(commandArgs) != 0 {
			return nil, fmt.Errorf("version command takes no arguments")
		}
		return &VersionCommand{}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}
func ShowHelp() {
	fmt.Println("GophKeeper - Secure Password Manager")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  register <username> <email> <password>  Register a new user")
	fmt.Println("    - username: 3-50 characters")
	fmt.Println("    - email: valid email format")
	fmt.Println("    - password: minimum 6 characters")
	fmt.Println("")
	fmt.Println("  login <username> <password>             Login to your account")
	fmt.Println("")
	fmt.Println("  add <type> <title> [data...]            Add new data")
	fmt.Println("    - type: login_password, text, binary, bank_card")
	fmt.Println("    - title: up to 255 characters")
	fmt.Println("    - data: type-specific arguments")
	fmt.Println("")
	fmt.Println("  list                                    List all data")
	fmt.Println("  get <id>                                Get specific data")
	fmt.Println("  delete <id>                             Delete data")
	fmt.Println("  sync                                    Synchronize with server")
	fmt.Println("  history <id>                            Show data history")
	fmt.Println("  help                                    Show this help")
	fmt.Println("  version                                 Show version information")
	fmt.Println("")
	fmt.Println("Data types and their required arguments:")
	fmt.Println("  login_password: username, password [, url] [, metadata]")
	fmt.Println("  text: content [, metadata]")
	fmt.Println("  binary: data [, metadata]")
	fmt.Println("  bank_card: number, holder, expiry, cvv [, metadata]")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  gophkeeper register john john@example.com mypassword123")
	fmt.Println("  gophkeeper add login_password \"My Bank\" john mypass123 https://bank.com \"Banking login\"")
	fmt.Println("  gophkeeper add text \"Important Note\" \"This is my secret note\" \"Personal notes\"")
}
