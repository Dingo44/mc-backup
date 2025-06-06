package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var failCount int = 0

func main() {
	fmt.Println("Starting...")

	setup()
}

func setup() {

	var filename string

	// get necessary files and check they exist
	fmt.Println("World name:")
	_, err := fmt.Scanln(&filename)

	if err != nil {
		log.Fatal(err)
	}

	if _, err = os.Stat(filename); err != nil {
		// log.Print(err)
		log.Fatal(err)
	}

	// make sure that the backups dir exists
	if _, err = os.Stat("./backups/"); err != nil {
		log.Print(err)
		log.Println("creating backups dir...")

		os.Mkdir("backups/", 0755)
	}

	fmt.Println("Backup interval (minutes):")
	var repeatInterval int
	_, err = fmt.Scanln(&repeatInterval)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Backup amount:")
	var backupAmount int
	_, err = fmt.Scanln(&backupAmount)

	if err != nil {
		log.Fatal(err)
	}

	// parse file into command string
	compressFile(filename, backupAmount)

	repeatCompress(repeatInterval, filename, backupAmount)
}

func compressFile(filename string, backupAmount int) {

	log.Println("Compressing: " + filename)
	fname := strings.TrimPrefix(filename, "../")
	cmd := exec.Command("tar", "-czvf", "./backups/"+fname+"."+string(time.Now().Format(time.UnixDate))+".tar.gz", filename)

	err := cmd.Run()

	if err != nil {

		// allow for 5 fails and if tar still fails then crash the program
		failCount++

		log.Println("Failed to compress: ", err, "\nTrying again...")
		compressFile(filename, backupAmount)

		if failCount > 5 {
			log.Fatal("Failed to compress: ", err)
		}
	} else {
		failCount = 0
	}

	log.Println("Compressed: " + filename)

	//

	checkOldFiles(backupAmount)
}

func checkOldFiles(backupAmount int) {
	files, _ := os.ReadDir("./backups/")
	// fmt.Println(len(files))

	if len(files) > backupAmount {
		for i := 0; i < len(files)-(backupAmount-1); i++ {
			log.Println(files[0])
			removeOldestFile(files[0])
			files, _ = os.ReadDir("./backups/")
		}
	}
}

func removeOldestFile(files fs.DirEntry) {
	// fmt.Println(files)
	err := os.Remove("./backups/" + files.Name())

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Removed old backup")
}

func repeatCompress(interval int, filename string, backupAmount int) {

	time.Sleep(time.Duration(interval) * time.Minute)
	// time.Sleep(time.Duration(interval) * time.Second)

	compressFile(filename, backupAmount)

	repeatCompress(interval, filename, backupAmount)
}
