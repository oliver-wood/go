package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "time"
)

func main() {
  
  file, err := os.Open("/Users/oliver/Dev/DogsTrust/messages/20180115/messages.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    var msgstrings = strings.Split(scanner.Text(), "\t")
    fmt.Println(msgstrings[0])
    var p = "/Users/oliver/Dev/DogsTrust/messages/20180115/" + msgstrings[0] + ".txt"
    createFile(p)

    msg := "Universe<!>Create<!>WG<!>391507<!>$date$<!>$msg$<!>CxIgnore"
    msg = strings.Replace(msg, "$date$", time.Now().Format("2006-01-02T15:04:05.000000+00:00)"), -1)
    msg = strings.Replace(msg, "$msg$", msgstrings[1], -1)

    writeFile(p, msg)
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }
}

func createFile(path string) {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if isError(err) { return }
		defer file.Close()
	}

	fmt.Println("==> done creating file", path)
}

func writeFile(path string, val string) {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) { return }
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(val)
	if isError(err) { return }

	// save changes
	err = file.Sync()
	if isError(err) { return }

	fmt.Println("==> done writing to file")
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}
