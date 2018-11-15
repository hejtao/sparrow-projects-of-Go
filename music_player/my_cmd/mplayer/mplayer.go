package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"my_pkg/mplayer/library"
	"my_pkg/mplayer/mp"
)

var lib *library.MusicManager

func handleLibCommands(tokens []string) {
	switch tokens[1] {
	case "list":
		for i := 0; i < lib.Len(); i++ {
			e, _ := lib.Get(i)
			fmt.Println(i+1, ":", e.Name, e.Artist, e.Source, e.Type)
		}

	case "add":
		if len(tokens) == 6 {
			lib.Add(&library.MusicEntry{tokens[2], tokens[3], tokens[4], tokens[5]})
		} else {
			fmt.Println("USAGE: lib add <name><artist><source><type>")
		}

	case "remove":
		if len(tokens) == 3 {
			lib.RemoveByName(tokens[2])
		} else {
			fmt.Println("USAGE: lib remove <name>")
		}

	default:
		fmt.Println("Unrecognized lib command :", tokens[1])
	}
}

func handlePlayCommands(tokens []string) {
	if len(tokens) != 2 {
		fmt.Println("USAGE: play <name>")
		return
	}

	e := lib.Find(tokens[1])
	if e == nil {
		fmt.Println("The music", tokens[1], "does not exist.")
		return
	}

	mp.Play(e.Source, e.Type)
}

func main() {
	fmt.Println(`		
			Enter following cmds to control the player:
			lib list -- view the existing music lib
			lib add <name><artist><source><type> -- Add a music to the music lib
			lib remove <name> -- Remove the specified music from the lib
			play <name> -- play the specified music
		`) // 多行打印，使用反引号
	lib = library.NewMusicManager()

	r := bufio.NewReader(os.Stdin) //?????

	for {
		fmt.Print("Enter cmd ->")

		rawLine, _, _ := r.ReadLine()
		line := string(rawLine)

		if line == "q" || line == "e" {
			break
		}

		tokens := strings.Split(line, " ")

		if tokens[0] == "lib" {
			handleLibCommands(tokens)
		} else if tokens[0] == "play" {
			handlePlayCommands(tokens)
		} else {
			fmt.Println("Unrecognized cmd:", tokens[0])
		}
	}
}
