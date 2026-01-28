package modes

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/mod"
)

func Import(input string) {
	if config.Interactive {
		defer func() {
			fmt.Print("Press any key ")
			fmt.Scanln()
		}()
	}

	err := os.MkdirAll(config.ModPath, 0o755)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Open(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	reader := bufio.NewReader(f)

	line, _, err := reader.ReadLine()
	for ; err == nil; line, _, err = reader.ReadLine() {
		update, err := mod.UpdateFromString(string(line))
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Downloading %s@%s - ", update.Name, update.Version.String())
		err = update.Download()
		if err != nil {
			fmt.Println("FAIL")
			fmt.Println(err)
			continue
		}
		fmt.Println("SUCCESS")
	}

	if err != io.EOF {
		fmt.Println(err)
		return
	}
	fmt.Println("Finished import")
}
