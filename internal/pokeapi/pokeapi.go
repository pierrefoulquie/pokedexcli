package pokeapi
import (
	"time"
	"fmt"
	"strconv"
	"encoding/json"
	"net/http"
	"net/url"
	"io"
	"github.com/pierrefoulquie/pokedexcli/internal/pokecache"
)

const BASE_URL = "https://pokeapi.co/api/v2/location-area/"
const POKE_URL = "https://pokeapi.co/api/v2/pokemon/"

func NewClient(base_url string, interval time.Duration) (*PokeAPIClient, error){
	client := PokeAPIClient{}
	client.baseURL = base_url
	client.Enc = Encounters{}
	client.cache =  pokecache.NewCache(interval)
	if err:= client.FetchBaseLocationArea(); err!=nil{
		return &client, err
	}
	return &client, nil
}

func (c *PokeAPIClient) FetchPokemon(poke string) error{
	endpoint := POKE_URL+poke
	fmt.Println(endpoint)

	if val, ok:= c.cache.Get(endpoint); ok{
		if err := json.Unmarshal(val, &c.pokemon); err!=nil{
			return err
		}
	//if not, http request and unmarshal
	}else{
		//get data through http request
		res, err := http.Get(endpoint)
		if err!=nil{
			return err
		}
		defer res.Body.Close()
		val,err := io.ReadAll(res.Body)
		if err!=nil{
			return err
		}
		//cache the data
		c.cache.Add(endpoint, val)

		//unmarshal the data
		marshErr := json.Unmarshal(val, &c.pokemon)
		if marshErr!=nil{
			return marshErr
		}
		fmt.Printf("Name : %v\n",c.pokemon.Name)
		fmt.Printf("Level : %v\n",c.pokemon.Xp)
	}
	return nil
}

func (c *PokeAPIClient) FetchEncounters(area string) error{
	fmt.Println("Exploring "+area+"...")
	endpoint := BASE_URL+area

	//if the data is already cached, unmarshal from the cache
	if val, ok:= c.cache.Get(endpoint); ok{
		if err := json.Unmarshal(val, &c.Enc); err!=nil{
			return err
		}
	//if not, http request and unmarshal
	}else{
		//get data through http request
		res, err := http.Get(endpoint)
		if err!=nil{
			return err
		}
		defer res.Body.Close()
		val,err := io.ReadAll(res.Body)
		if err!=nil{
			return err
		}
		//cache the data
		c.cache.Add(endpoint, val)

		//unmarshal the data
		marshErr := json.Unmarshal(val, &c.Enc)
		if marshErr!=nil{
			return marshErr
		}
	}

	//print result
	numPokemon := len(c.Enc.PokemonEncounters)
	if numPokemon == 0{
		fmt.Println("No Pokemon found")
		return  nil
	}
	for _, poke := range c.Enc.PokemonEncounters{
		fmt.Printf(" - %v\n",poke.Pokemon.Name)
	}
	return nil
}

func (c *PokeAPIClient) FetchLocationAreas(url string) error{
	if val, ok:= c.cache.Get(url); ok{
		if err := json.Unmarshal(val, &c.Res); err!=nil{
			return err
		}
		return nil
	}

	res, err := http.Get(url)
	if err!=nil{
		return err
	}
	defer res.Body.Close()
	val,err := io.ReadAll(res.Body)
	if err!=nil{
		return err
	}
	marshErr := json.Unmarshal(val, &c.Res)
	if marshErr!=nil{
		return marshErr
	}
	c.cache.Add(url, val)
	return nil
}

func (c *PokeAPIClient) FetchBaseLocationArea() error{
	if err := c.FetchLocationAreas(BASE_URL); err!=nil{
		return err
	}
	c.Res.Previous = BASE_URL + "?offset=0&limit=20"
	c.Res.Next = BASE_URL + "?offset=0&limit=20"
	return nil
}

func (c *PokeAPIClient) FetchPreviousLocationArea() error{
	if err := c.FetchLocationAreas(c.Res.Previous); err!=nil{
		return err
	}
	return nil
}

func (c *PokeAPIClient) FetchNextLocationArea() error{
	if err := c.FetchLocationAreas(c.Res.Next); err!=nil{
		return err
	}
	return nil
}

func (c *PokeAPIClient) DetectFirstPage() (bool, error){
	//if next page offset is 20, we are on the first page
	u, err := url.Parse(c.Res.Next)
	if err!=nil{
		return false, err
	}
	if u.Query().Get("offset")=="20"{
		return true, err
	}
	return false, err
}

func (c* PokeAPIClient) DetectLastPage() (bool, error){
	//if previous page limit is less than 20, we are on the last page
	u, err := url.Parse(c.Res.Previous)
	if err!=nil{
		return false, err
	}
	if u.Query().Get("limit")!="20"{
		return true, nil
	}
	return false, nil 
}

func (c* PokeAPIClient) CorrectPrevious() error{
	u, err := url.Parse(c.Res.Next)
	if err != nil{
		return err
	}

	if num, err := strconv.Atoi(u.Query().Get("offset")); err != nil{
		fmt.Println("### CORRECTION FAILED ###")
	}else{
		newOffset := num - 20
		newLimit := "20"
		c.Res.Previous = fmt.Sprintf("%s?offset=%v&limit=%s",BASE_URL,newOffset,newLimit)
	}
	return err
}
