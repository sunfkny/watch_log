package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

var logDir = flag.String("d", "logs", "log directory")
var highLightMap = map[*regexp.Regexp]string{
	regexp.MustCompile(`INFO`):                              color.New(color.FgGreen).Sprint("$0"),
	regexp.MustCompile(`WARNING`):                           color.New(color.FgYellow).Sprint("$0"),
	regexp.MustCompile(`ERROR`):                             color.New(color.FgRed).Sprint("$0"),
	regexp.MustCompile(`HTTP/[1-2]\.[0-2] [2-3][0-9][0-9]`): color.New(color.FgGreen).Sprint("$0"),
	regexp.MustCompile(`HTTP/[1-2]\.[0-2] 4[0-9][0-9]`):     color.New(color.FgYellow).Sprint("$0"),
	regexp.MustCompile(`HTTP/[1-2]\.[0-2] 5[0-9][0-9]`):     color.New(color.FgRed).Sprint("$0"),
}

func watch(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", fileName, err)
	}
	defer file.Close()

	// 读取文件末尾后继续读取新增内容
	file.Seek(0, 2)
	reader := bufio.NewReader(file)

	// 循环读取文件内容
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// 高亮显示关键字
		for r, c := range highLightMap {
			line = r.ReplaceAllString(line, c)
		}

		// 输出高亮显示的内容
		fmt.Print(line)
	}
}

func main() {
	flag.Parse()

	glob := filepath.Join(*logDir, "*.log")
	fmt.Printf("Glob: %s\n", glob)
	files, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		fmt.Println("No logs found")
		os.Exit(1)
	}

	prompt := promptui.Select{
		Label: "Select a log file",
		Items: files,
	}
	_, file, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("Watching log file: %s\n", file)
	watch(file)

}
