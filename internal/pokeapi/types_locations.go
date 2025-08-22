package pokeapi
import(
	"github.com/pierrefoulquie/pokedexcli/internal/pokecache"
)
type PokeAPIClient struct{
	baseURL		string
	Res			LocationResponse
	PokeRes		PokeResponse
	Enc			Encounters
	cache		*pokecache.Cache
	Pokemon		Pokemon
	Pokedex		map[string]Pokemon
}

type PageResponse struct{
	Name 	string	`json:"name"`
	URL 	string 	`json:"url"`
}

type PokeResponse struct{
	Count int				`json:"count"`
	Next string				`json:"next"`
	Previous string			`json:"previous"`
	Results []PageResponse	`json:"results"`
}

type LocationResponse struct{
	Count int				`json:"count"`
	Next string				`json:"next"`
	Previous string			`json:"previous"`
	Results []PageResponse	`json:"results"`
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
	Name 	string  `json:"name"`
	Height	int		`json:"height"`
	Weight	int		`json:"weight"`
	Stats	[]PokeStats	`json:"stats"`
	Types	[]PokeTypes	`json:"types"`
	Xp  	int 	`json:"base_experience"` 
}

type PokeStats struct{
	BaseStat int	`json:"base_stat"`
	Stat Stat	`json:"stat"`
}

type Stat struct{
	Name	string	`json:"name"`
}

type PokeTypes struct{
	Type 	Type	`json:"type"`
}

type Type struct{
	Name	string	`json:"name"`
}

