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
	callback 	func(*pokeapi.PokeAPIClient, string) error
	arg0 		string
}

func commandExit(c *pokeapi.PokeAPIClient, arg0 string) error{
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *pokeapi.PokeAPIClient, arg0 string) error{
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("help: Displays a help message")
	fmt.Println("map: Displays 20  next pokeworld locations")
	fmt.Println("mapb: Displays 20 previous pokeworld locations")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

func commandMap(c *pokeapi.PokeAPIClient, arg0 string) error{
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

func commandMapb(c *pokeapi.PokeAPIClient, arg0 string) error{
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

func commandExplore(c *pokeapi.PokeAPIClient, arg0 string) error{
	if err := c.FetchEncounters(arg0); err!=nil{
		fmt.Println(err)
	}
	return nil
}

func commandCatch(c *pokeapi.PokeAPIClient, arg0 string) error{
	for _, poke := range c.Enc.PokemonEncounters{
		if poke.Pokemon.Name == arg0 {
			fmt.Printf("Throwing a pokeball at %v\n", arg0)
			return nil
		}
	}
	fmt.Println("Pokemon not found")
	return nil
}
var commands = map[string]cliCommand{
	"catch":{
        name:        "catch",
        description: "Throw a pokeball at a pokemon",
        callback:    commandCatch,
	},
	"explore":{
        name:        "explore",
        description: "Displays Pokemons presents in the area",
        callback:    commandExplore,
	},
    "exit": {
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    commandExit,
    },
    "help": {
        name:        "help",
        description: "Displays Pokedex Manual",
        callback:    commandHelp,
    },
    "map": {
        name:        "map",
        description: "Displays 20 next Pokemon world areas",
        callback:    commandMap,
    },
    "mapb": {
        name:        "mapb",
        description: "Displays 20 previous Pokemon world areas",
        callback:    commandMapb,
    },
}

func main(){
	scanner := bufio.NewScanner(os.Stdin)
	client, clErr := pokeapi.NewClient(pokeapi.BASE_URL, 10*time.Second)
	arg0 := ""
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
			cleanString :=  cleanInput(text)
			firstWord := strings.ToLower(cleanString[0])
			if _,ok := commands[firstWord]; !ok{
				fmt.Println("Error: command unknown")
				continue
			}
			if (firstWord == "explore" || firstWord == "catch") && len(cleanString) < 2{
				fmt.Println("Error: missing argument")
				continue
			}else if firstWord == "explore" || firstWord == "catch"{
				arg0 = strings.ToLower(cleanString[1])
			}
			if err := commands[firstWord].callback(client, arg0); err!=nil{
				fmt.Println("Error:",err)
			}
		}
	}
}

