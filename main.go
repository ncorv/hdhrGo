package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {

	hdhrconf := exec.Command("hdhomerun_config", "discover")
	out, err := hdhrconf.CombinedOutput()
	configCommandOutput := strings.Split(string(out), " ")
	hdhrIP := ""
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		for i := 0; i <= len(configCommandOutput)-1; i++ {
			if i == len(configCommandOutput)-1 {
				fmt.Print(configCommandOutput[i])
			} else {
				fmt.Print(configCommandOutput[i] + " ")
			}
		}
	}

	//check command output
	// if 4th element is "found" and 6th is an ip address
	if configCommandOutput[3] == "found" && net.ParseIP(strings.Trim(configCommandOutput[5], "\n")) != nil {
		hdhrIP = strings.Trim(configCommandOutput[5], "\n")
		fmt.Println("hdhomerun found at: " + configCommandOutput[5])
	} else {
		fmt.Println("hdhomerun device not found.")
	}

	// do we have ffmpeg
	ffmpegVersion := exec.Command("ffmpeg", "-version")
	out, err = ffmpegVersion.CombinedOutput()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	} else {

		fmt.Println("Starting ffmpeg stream to: " + os.Args[1] + " with quality: " + os.Args[2] + " on channel: " + os.Args[3])
		ffmpeg := exec.Command("ffmpeg", "-y", "-i", "http://"+hdhrIP+":5004/auto/v"+os.Args[3], "-r", "30", "-s", "hd720", "-threads", "4", "-vcodec", "libx264", "-crf", os.Args[2], "-async", "1", "-acodec", "aac", "-f", "flv", os.Args[1])

		ffmpeg.Stderr = os.Stdout
		ffmpeg.Stdin = os.Stdin
		if err := ffmpeg.Run(); err != nil {
			log.Fatalln(err)
		}

	}
}
