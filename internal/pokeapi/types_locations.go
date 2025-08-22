package pokeapi
import(
	"github.com/pierrefoulquie/pokedexcli/internal/pokecache"
)
type PokeAPIClient struct{
	baseURL		string
	Res			LocationResponse
	Enc			Encounters
	cache		*pokecache.Cache
	pokemon		Pokemon
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
	Pokemon PokeEncounter `json:"pokemon"`
}

type PokeEncounter struct {
	Name string `json:"name"`
	URL  string `json:"url"` 
}

type Pokemon struct{
	Name 	string `json:"name"`
	Xp  	int 	`json:"base_experience"` 
}
