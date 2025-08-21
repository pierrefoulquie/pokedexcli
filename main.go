package main
import (
	"log"
	"time"
	"fmt"
	"strings"
	"bufio"
	"os"
	"github.com/pierrefoulquie/pokedexcli/internal/pokeapi"
)

func cleanInput(text string) []string{
	output := strings.Fields(text)
	return output
}

type cliCommand struct{
	name 		string
	description string
	callback 	func(*pokeapi.PokeAPIClient) error
}

func commandExit(c *pokeapi.PokeAPIClient) error{
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *pokeapi.PokeAPIClient) error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("help: Displays a help message")
	fmt.Println("map: Displays 20  next pokeworld locations")
	fmt.Println("mapb: Displays 20 previous pokeworld locations")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

func commandMap(c *pokeapi.PokeAPIClient) error{
	c.FetchNextLocationArea()
	for i:=0;i<len(c.Res.Results);i++{
		fmt.Println(c.Res.Results[i].Name)
	}
	last, errL := c.DetectLastPage()
	if errL!=nil{
		return errL
	}
	if last{
		fmt.Println("You're on the last page")
	}
	return nil
}

func commandMapb(c *pokeapi.PokeAPIClient) error{
	first, err := c.DetectFirstPage()
	if err!=nil{
		return err
	}
	if first{
		fmt.Println("You're on the first page")
		//if not first page, check if last and print result
	}else{
		//detect last page to correct previous limit
		last, errL := c.DetectLastPage()
		if errL!=nil{
			return errL
		}
		if last{
			if errP := c.CorrectPrevious(); errP!=nil{
				return errP
			}
		}
		c.FetchPreviousLocationArea()
		for i:=0;i<len(c.Res.Results);i++{
			fmt.Println(c.Res.Results[i].Name)
		}
	}
	return nil
}

var commands = map[string]cliCommand{
    "exit": {
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    commandExit,
    },
    "help": {
        name:        "help",
        description: "Display Pokedex Manual",
        callback:    commandHelp,
    },
    "map": {
        name:        "map",
        description: "displays 20 next Pokemon world cities",
        callback:    commandMap,
    },
    "mapb": {
        name:        "mapb",
        description: "displays 20 previous Pokemon world cities",
        callback:    commandMapb,
    },
}

func main(){
	scanner := bufio.NewScanner(os.Stdin)
	client, clErr := pokeapi.NewClient(pokeapi.BASE_URL, 10*time.Second)
	if clErr!=nil{
		log.Fatal()
	}

	for{
		fmt.Printf("Pokedex > ")
		more := scanner.Scan()
		if more {
			text := scanner.Text()
			if text == ""{
				fmt.Println("Error: enpty command")
				continue
			}
			clean_string :=  cleanInput(text)
			first_word := strings.ToLower(clean_string[0])
			if _,ok := commands[first_word]; !ok{
				fmt.Println("Error: command unknown")
				continue
			}
			if err := commands[first_word].callback(client); err!=nil{
				fmt.Println("Error:",err)
			}
		}
	}
}

