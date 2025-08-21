package pokeapi
import(
	"github.com/pierrefoulquie/pokedexcli/internal/pokecache"
)
type PokeAPIClient struct{
	baseURL		string
	Res			LocationResponse
	Enc			Encounters
	cache		*pokecache.Cache
}


type Location struct{
	Name 	string	`json:"name"`
	URL 	string 	`json:"url"`
}

type LocationResponse struct{
	Next string			`json:"next"`
	Previous string		`json:"previous"`
	Results []Location	`json:"results"`
}

type Encounters struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"` 
}
