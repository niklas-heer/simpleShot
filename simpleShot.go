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

/**
 * The main function of the program.
 */
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

		// config structure
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

		// initalize and declare the needed variables
		screenPath 	:= getHomeDir() + "/" + c.String("folder")
		configPath 	:= getHomeDir() + "/.simpleShot.gcfg"
		fileformat 	:= c.String("type")
		filename 	:= randStr(c.Int("name-length"), c.String("name-alphabet")) + "." + fileformat
		filepath 	:= screenPath + "/" + filename
		command 	:= "import -frame "

		// read the config file
		err := gcfg.ReadFileInto(&cfg, configPath)
		if err != nil {
			log.Fatal(err)
		}

		// make a unified file url
		fileurl := cfg.Ftp.Url + filename

		
		// set the user defined path.
		// create the directory, if it doesn't exists.
		if c.String("folder") != "screenshots" {
			screenPath = c.String("folder")
			makeDir(screenPath)
		} else {
			makeDir(screenPath)
		}

		// behavior for the selection option 
		if c.Bool("select") {
			command = "import -frame "
		}

		// take a screenshot
		wg := new(sync.WaitGroup)
		commands := []string{command + filepath}
		for _, str := range commands {
			wg.Add(1)
			go exe_cmd(str, wg)
		}
		wg.Wait()

		// upload if the option is set
		if c.Bool("upload") {

			// upload the file
			uploadFTP(cfg.Ftp.Port, cfg.Ftp.Server, cfg.Ftp.User, cfg.Ftp.Pw, cfg.Ftp.Path, filepath, filename)

			// copy url to clipboard
			copyToClipboard(fileurl)
		}

		// send notifications if the option "quiet" is NOT set
		if !c.Bool("quiet") {

			// if we uploaded the picture send a special notification, else send the default one
			if c.Bool("upload") {
				sendNotification("The images url was copied to the clipboard and uploaded to: " + fileurl)
			} else {
				sendNotification("Image was saved under: " + filepath)
			}
		}
	}

	app.Run(os.Args)
}

/**
 * Copies the given text into the clipboard
 * 
 * @param text to be copied into the clipboard
 */
func copyToClipboard(text string) {
	clipboard.WriteAll(text)
}

/**
 * Sends the given text as desktop notification
 * 
 * @param text which will be send as notification
 */
func sendNotification(text string) {

	notify.Init("simpleShot")
	info := notify.NotificationNew("simpleShot", text, "")

	// send error message if needed
	if info == nil {
		fmt.Fprintf(os.Stderr, "Unable to create a new notification\n")
		return
	}

	// set the timeout for the notification
	notify.NotificationSetTimeout(info, DELAY)

	// show the notification
	if e := notify.NotificationShow(info); e != nil {
		fmt.Fprintf(os.Stderr, "%s\n", e.Message())
		return
	}

	time.Sleep(DELAY * 1000000)

	// close the notification
	notify.NotificationClose(info)

	notify.UnInit()
}

/**
 * Uploads a given file via FTP to a given server.
 * 
 * @param port - the port of the FTP server.
 * @param server - the adress of the FTP server. (IP or domain name)
 * @param user - the username to log into the FTP server.
 * @param pw - the password to log into the FTP server.
 * @param serverpath - the path on the FTP server, where the file will be stored.
 * @param filepath - the local path to the file, which will be uploaded.
 * @param name - the name of the file which will be uploaded.
 */
func uploadFTP(port int, server, user, pw, serverpath, filepath, name string) {

	debugNum := 0

	// enable debugging if the option is set
	if Debug {
		debugNum = 1
	}

	// 1 for debugging
	ftpClient := ftp4go.NewFTP(debugNum) 

	// connect
	_, err := ftpClient.Connect(server, port, "")
	if err != nil {
		log.Fatal(err)
	}

	// login
	_, err = ftpClient.Login(user, pw, "")
	if err != nil {
		log.Fatal(err)
	}

	// go into the right folder
	if serverpath != "" {
		_, err = ftpClient.Cwd(serverpath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// upload the file
	err = ftpClient.UploadFile(name, filepath, true, nil)
	if err != nil {
		log.Fatal(err)
	}

	// close the connection
	defer ftpClient.Quit()
}

/**
 * Checks if a given file or directory already exists.
 * 
 * @param path - the path to the file or directory which will be checked.
 *
 * @returns whether the given file or directory exists or not. (true or false)
 */
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

/**
 * Creates a directory in the specified path, if it doesn't exist already.
 * 
 * @param dirPath - the path to the directory, which will be created.
 *
 * @note It will also create all parent directories if thez don't exist.
 */
func makeDir(dirPath string) {
	exists, err := exists(dirPath)

	if err != nil {
		log.Fatal(err)
	}

	// only make the directory if it doesn't exists already
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

/**
 * Gets the full path to the home directory of the current user.
 *
 * @returns the full path to the home directory of the current user as string.
 */
func getHomeDir() string {
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir
}

/**
 * Generates a random string.
 * 
 * @param strSize - the length of the string which will be generated.
 * @param randType - the dictionary which conatins the character, which will be used. (alpha, number and alphanum)
 *
 * @returns the the randomly genarted string.
 */
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

/**
 * Executes a given command.
 * 
 * @param cmd - the command which will be executed. (e.g. "ls -al")
 * @param wg - the Waitgroup, which is needed to execute the goroutine. (wg := new(sync.WaitGroup))
 *
 * @note atm it's only desinged to execute GNU/Linux commands. Although Unix commands might also work.
 */
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

	// need to signal to the waitgroup that this goroutine is done
	wg.Done() 
}
