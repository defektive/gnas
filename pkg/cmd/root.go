package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/defektive/gnas/pkg/process"
	"github.com/reeflective/console"
	"github.com/reeflective/readline"
	"github.com/spf13/cobra"
)

var ServerListener = ":8000"
var UploadServer = "http://127.0.0.1"
var UploadToken = "you should change this at build time"

var UploadEndpoint = UploadServer + ServerListener

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gnas",
	Short: "A cross platform shell like program with mini-programs",
	Long:  `A cross platform shell like program with mini-programs`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		app := console.New("gnas")
		app.NewlineBefore = true
		app.NewlineAfter = true

		app.SetPrintLogo(func(_ *console.Console) {
			fmt.Print(`
  ▄████  ███▄    █  ▄▄▄        ██████ 
 ██▒ ▀█▒ ██ ▀█   █ ▒████▄    ▒██    ▒ 
▒██░▄▄▄░▓██  ▀█ ██▒▒██  ▀█▄  ░ ▓██▄   
░▓█  ██▓▓██▒  ▐▌██▒░██▄▄▄▄██   ▒   ██▒
░▒▓███▀▒▒██░   ▓██░ ▓█   ▓██▒▒██████▒▒
 ░▒   ▒ ░ ▒░   ▒ ▒  ▒▒   ▓▒█░▒ ▒▓▒ ▒ ░
  ░   ░ ░ ░░   ░ ▒░  ▒   ▒▒ ░░ ░▒  ░ ░
░ ░   ░    ░   ░ ░   ░   ▒   ░  ░  ░  
      ░          ░       ░  ░      ░  

GNAS is Not A Shell

`)
		})

		menu := app.ActiveMenu()

		// Set some custom prompt handlers for this menu.
		setupPrompt(menu)

		// All menus currently each have a distinct, in-memory history source.
		// Replace the main (current) menu's history with one writing to our
		// application history file. The default history is named after its menu.
		hist, _ := embeddedHistory(".gnas-history")
		menu.AddHistorySource("local history", hist)

		// We bind a special handler for this menu, which will exit the
		// application (with confirm), when the shell readline receives
		// a Ctrl-D keystroke. You can map any error to any handler.
		menu.AddInterrupt(io.EOF, exitCtrlD)

		// Make a command yielder for our main menu.
		// menu.SetCommands(makeflagsCommands(app))
		// Thanks ChatGPT for generating this for us!
		//InternalCmd.AddCommand(RootCmd.Commands())
		menu.SetCommands(func() *cobra.Command { return cmd })

		// Everything is ready for a tour.
		// Run the console and take a look around.
		app.Start()
	},
}

//// InternalCmd represents the base command when called without any subcommands
//var InternalCmd = &cobra.Command{
//	Use:   "gnas",
//	Short: "A cross platform shell like program with mini-programs",
//	Long:  `A cross platform shell like program with mini-programs`,
//	// Uncomment the following line if your bare application
//	// has an action associated with it:
//	//Run: func(cmd *cobra.Command, args []string) {
//	//
//	//
//	//},
//}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gnas.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// errorCtrlSwitchMenu is a custom interrupt handler which will
// switch back to the main menu when the current menu receives
// a CtrlD (io.EOF) error.
func errorCtrlSwitchMenu(c *console.Console) {
	fmt.Println("Switching back to main menu")
	c.SwitchMenu("")
}

// setupPrompt is a function which sets up the prompts for the main menu.
func setupPrompt(m *console.Menu) {
	p := m.Prompt()

	p.Primary = func() string {
		prompt := "\x1b[33m%s\x1b[0m [%s] <%s> in \x1b[34m%s\x1b[0m\n> "
		wd, _ := os.Getwd()

		dir, err := filepath.Rel(os.Getenv("HOME"), wd)
		if err != nil {
			dir = filepath.Base(wd)
		}

		return fmt.Sprintf(prompt, "GNAS", getTime(), getIntegrity(), dir)
	}

	p.Transient = func() string { return "\x1b[1;30m" + ">> " + "\x1b[0m" }
}

func getTime() string {
	return time.Now().Format(time.RFC3339)
}

func getIntegrity() string {
	if process.IsAdmin() {
		return "HIGH"
	}
	return "LOW"
}

// exitCtrlD is a custom interrupt handler to use when the shell
// readline receives an io.EOF error, which is returned with CtrlD.
func exitCtrlD(c *console.Console) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Confirm exit (Y/y): ")

	text, _ := reader.ReadString('\n')
	answer := strings.TrimSpace(text)

	if (answer == "Y") || (answer == "y") {
		os.Exit(0)
	}
}

func switchMenu(c *console.Console) {
	fmt.Println("Switching to client menu")
	c.SwitchMenu("client")
}

var (
	errOpenHistoryFile = errors.New("failed to open history file")
	errNegativeIndex   = errors.New("cannot use a negative index when requesting historic commands")
	errOutOfRangeIndex = errors.New("index requested greater than number of items in history")
)

type fileHistory struct {
	file  string
	lines []Item
}

type Item struct {
	Index    int
	DateTime time.Time
	Block    string
}

// NewSourceFromFile returns a new history source writing to and reading from a file.
func embeddedHistory(file string) (readline.History, error) {
	var err error

	hist := new(fileHistory)
	hist.file = file
	hist.lines, err = openHist(file)

	return hist, err
}

func openHist(filename string) (list []Item, err error) {
	//file, err := historyFile.Open(filename)
	//if err != nil {
	//	return list, fmt.Errorf("error opening history file: %s", err.Error())
	//}
	//
	//scanner := bufio.NewScanner(file)
	//for scanner.Scan() {
	//	var item Item
	//
	//	err := json.Unmarshal(scanner.Bytes(), &item)
	//	if err != nil || len(item.Block) == 0 {
	//		continue
	//	}
	//
	//	item.Index = len(list)
	//	list = append(list, item)
	//}
	//
	//file.Close()

	return list, nil
}

// Write item to history file.
func (h *fileHistory) Write(s string) (int, error) {
	block := strings.TrimSpace(s)
	if block == "" {
		return 0, nil
	}

	item := Item{
		DateTime: time.Now(),
		Block:    block,
		Index:    len(h.lines),
	}

	if len(h.lines) == 0 || h.lines[len(h.lines)-1].Block != block {
		h.lines = append(h.lines, item)
	}

	// line := struct {
	// 	DateTime time.Time `json:"datetime"`
	// 	Block    string    `json:"block"`
	// }{
	// 	Block:    block,
	// 	DateTime: item.DateTime,
	// }
	//
	// data, err := json.Marshal(line)
	// if err != nil {
	// 	return h.Len(), err
	// }
	//
	// f, err := historyFile.Open(h.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	// if err != nil {
	// 	return 0, fmt.Errorf("%w: %s", errOpenHistoryFile, err.Error())
	// }
	//
	// _, err = f.Write(append(data, '\n'))
	// f.Close()

	return h.Len(), nil
}

// GetLine returns a specific line from the history file.
func (h *fileHistory) GetLine(pos int) (string, error) {
	if pos < 0 {
		return "", errNegativeIndex
	}

	if pos < len(h.lines) {
		return h.lines[pos].Block, nil
	}

	return "", errOutOfRangeIndex
}

// Len returns the number of items in the history file.
func (h *fileHistory) Len() int {
	return len(h.lines)
}

// Dump returns the entire history file.
func (h *fileHistory) Dump() interface{} {
	return h.lines
}
