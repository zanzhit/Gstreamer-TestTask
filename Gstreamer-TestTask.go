package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func readUrls(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	read := bufio.NewScanner(file)

	for read.Scan() {
		urls = append(urls, read.Text())
	}

	if err := read.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func recordStream(urls []string, done chan struct{}) {
	var commands []*exec.Cmd

	for i, url := range urls {
		timestamp := time.Now().Format("20060102_150405")
		outputFileName := fmt.Sprintf("location=stream_%d_%s.mkv", i+1, timestamp)
		urlArg := fmt.Sprintf("location=%s", url)

		recordFn := func(urlArg string, outputFileName string) *exec.Cmd {
			cmd := exec.Command("gst-launch-1.0", "rtspsrc", urlArg, "!", "rtph264depay", "!", "h264parse", "!", "matroskamux", "!", "filesink", outputFileName)
			commands = append(commands, cmd)
			return cmd
		}

		cmd := recordFn(urlArg, outputFileName)

		fmt.Printf("Запись потока %d началась.\n", i+1)

		err := cmd.Start()
		if err != nil {
			panic(err)
		}

		// Проверка на доступ к видеопотоку
		go func(i int, timestamp string, outputFileName string) {
			checkFileName := fmt.Sprintf("stream_%d_%s.mkv", i+1, timestamp)

			time.Sleep(time.Second * 4)

			fileInfo, err := os.Stat(checkFileName)
			if err != nil {
				fmt.Println("Ошибка при записи файла.", err)
			}
			select {
			case <-done:
				cmd.Process.Kill()
			default:
				for fileInfo.Size() == 0 {
					cmd := recordFn(urlArg, outputFileName)
					fmt.Printf("Повторная запись потока %d\n", i+1)

					err := cmd.Start()
					if err != nil {
						panic(err)
					}

					time.Sleep(time.Second * 4)

					fileInfo, err = os.Stat(checkFileName)
					if err != nil {
						fmt.Println("Ошибка при записи файла.", err)
					}
				}
			}
		}(i, timestamp, outputFileName)

		go func() {
			<-done
			cmd.Process.Kill()
		}()
	}

	for i, cmd := range commands {
		err := cmd.Wait()
		if err != nil {
			fmt.Printf("Запись потока %d прервана.\n", i+1)
		}
	}
}

func main() {
	filePath := "urls.txt"
	urls, err := readUrls(filePath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла: ", err)
		return
	}

	done := make(chan struct{})

	go func() {
		var input string
		fmt.Scanln(&input)
		close(done)
	}()

	recordStream(urls, done)

	fmt.Println("Программа завершена")
}
