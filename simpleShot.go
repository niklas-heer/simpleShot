package main

import (
	ftp4go "code.google.com/p/ftp4go"
	gcfg "code.google.com/p/gcfg"
	"crypto/rand"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/codegangsta/cli"
	notify "github.com/lenormf/go-notify"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"sync"
	"time"
)

var (
	Debug = false
)

const (
	DELAY = 3000
)

func main() {
	app := cli.NewApp()
	app.Name = "simpleShot"
	app.Version = "0.2.0"
	app.Author = "Niklas Heer"
	app.Email = "niklas.heer@gmail.com"
	app.Usage = "Takes a screenshot, uploads it via FTP and copies the url into your clipboard!"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "select, s",
			Usage: "Select the area for the screenshot.",
		},
		cli.BoolFlag{
			Name:  "upload,u",
			Usage: "Reads the credentials under ~/.simpleShot.gcfg and uploads it.",
		},
		cli.BoolFlag{
			Name:  "quiet,q",
			Usage: "Don't notify me!",
		},
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "Enable debugging.",
		},
		cli.IntFlag{
			Name:  "name-length,nl",
			Value: 6,
			Usage: "Choose the length for the name generator.",
		},
		cli.StringFlag{
			Name:  "name-alphabet,na",
			Value: "alphanum",
			Usage: "Choose the alphabet for the name generator.",
		},
		cli.StringFlag{
			Name:  "folder,f",
			Value: "screenshots",
			Usage: "Choose the directory where the screenshots are saved. (e.g. /home/nh/screenshots)",
		},
		cli.StringFlag{
			Name:  "type,t",
			Value: "jpg",
			Usage: "Choose the format of the screenshot. (png, jpg...)",
		},
	}

	app.Action = func(c *cli.Context) {

		// Config structure
		cfg := struct {
			Ftp struct {
				Url    string
				Server string
				Port   int
				Path   string
				User   string
				Pw     string
			}
		}{}

		Debug = c.Bool("debug")

		screenPath := getHomeDir() + "/" + c.String("folder")
		configPath := getHomeDir() + "/.simpleShot.gcfg"
		fileformat := c.String("type")
		filename := randStr(c.Int("name-length"), c.String("name-alphabet")) + "." + fileformat
		filepath := screenPath + "/" + filename
		command := "import -frame "

		// Read Config file
		err := gcfg.ReadFileInto(&cfg, configPath)
		if err != nil {
			log.Fatal(err)
		}

		// make a unified file url
		fileurl := cfg.Ftp.Url + filename

		/* Set the user defined path and make the dir if it doesn't exists else only make the default dir if it doesn't exist. */
		if c.String("folder") != "screenshots" {
			screenPath = c.String("folder")
			makeDir(screenPath)
		} else {
			makeDir(screenPath)
		}

		/* Behavior for selection option */
		if c.Bool("select") {
			command = "import -frame "
		}

		// Take the screenshot
		wg := new(sync.WaitGroup)
		commands := []string{command + filepath}
		for _, str := range commands {
			wg.Add(1)
			go exe_cmd(str, wg)
		}
		wg.Wait()

		/* Upload if the option is set */
		if c.Bool("upload") {

			// Upload the file
			uploadFTP(cfg.Ftp.Port, cfg.Ftp.Server, cfg.Ftp.User, cfg.Ftp.Pw, cfg.Ftp.Path, filepath, filename)

			// Copy url to clipboard
			copyToClipboard(fileurl)
		}

		if !c.Bool("quiet") {

			//if we uploaded picture send special notification else send default one
			if c.Bool("upload") {
				sendNotification("The images url was copied to the clipboard and uploaded to: " + fileurl)
			} else {
				sendNotification("Image was saved under: " + filepath)
			}
		}
	}

	app.Run(os.Args)
}

func copyToClipboard(text string) {
	clipboard.WriteAll(text)
}

func sendNotification(text string) {

	notify.Init("simpleShot")
	info := notify.NotificationNew("simpleShot",
		text,
		"")

	if info == nil {
		fmt.Fprintf(os.Stderr, "Unable to create a new notification\n")
		return
	}
	// info.SetTimeout(3000)
	notify.NotificationSetTimeout(info, DELAY)

	// info.Show()
	if e := notify.NotificationShow(info); e != nil {
		fmt.Fprintf(os.Stderr, "%s\n", e.Message())
		return
	}

	time.Sleep(DELAY * 1000000)

	// info.Close()
	notify.NotificationClose(info)

	notify.UnInit()
}

func uploadFTP(port int, server, user, pw, serverpath, filepath, name string) {

	debugNum := 0

	if Debug {
		debugNum = 1
	}

	ftpClient := ftp4go.NewFTP(debugNum) // 1 for debugging

	//connect
	_, err := ftpClient.Connect(server, port, "")
	if err != nil {
		log.Fatal(err)
	}

	//login
	_, err = ftpClient.Login(user, pw, "")
	if err != nil {
		log.Fatal(err)
	}

	//go into the right folder
	if serverpath != "" {
		_, err = ftpClient.Cwd(serverpath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// upload
	err = ftpClient.UploadFile(name, filepath, true, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ftpClient.Quit()
}

func makeDir(dirPath string) {
	exists, err := exists(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		wg := new(sync.WaitGroup)
		commands := []string{"mkdir -p " + dirPath}
		for _, str := range commands {
			wg.Add(1)
			go exe_cmd(str, wg)
		}
		wg.Wait()
	}
}

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func randStr(strSize int, randType string) string {
	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}

	return string(bytes)
}

func exe_cmd(cmd string, wg *sync.WaitGroup) {
	if Debug {
		fmt.Println("command is ", cmd)
	}
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		log.Fatal(err)
	}
	// DEBUGGING
	if Debug {
		fmt.Printf("%s", out)
	}
	wg.Done() // Need to signal to waitgroup that this goroutine is done
}
