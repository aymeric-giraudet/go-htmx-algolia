package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

type Hit struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

//go:embed index.html
var f embed.FS

func main() {
	var err error
	testTemplate, err := template.New("index.html").Funcs(template.FuncMap{
		"contains": func(s []string, e string) bool {
			for _, a := range s {
				if a == e {
					return true
				}
			}
			return false
		},
	}).ParseFS(f, "index.html")
	if err != nil {
		panic(err)
	}

	client := search.NewClient("latency", "6be0576ff61c053d5f9a3225e2a90f76")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		isHtmxRequest := r.Header.Get("Hx-Request") == "true"
		query := r.URL.Query().Get("query")
		page := r.URL.Query().Get("page")
		if page == "" {
			page = "0"
		}
		pageNumber, _ := strconv.Atoi(page)
		brands := r.URL.Query()["brand"]

		facetFilters := make([]interface{}, len(brands))
		for i, brand := range brands {
			facetFilters[i] = opt.FacetFilter("brand:" + brand)
		}
		queries := []search.IndexedQuery{
			search.NewIndexedQuery("instant_search", opt.Query(query), opt.Page(pageNumber), opt.HitsPerPage(5), opt.Facets("brand"), opt.MaxValuesPerFacet(10), opt.FacetFilterOr(facetFilters...)),
			search.NewIndexedQuery("instant_search", opt.Query(query), opt.HitsPerPage(0), opt.Facets("brand"), opt.MaxValuesPerFacet(10)),
		}
		results, _ := client.MultipleQueries(queries, "")
		var hits []Hit
		results.Results[0].UnmarshalHits(&hits)

		// Doing this to keep the form state when selected facets are not displayed
		brandsToDisplay := results.Results[1].Facets["brand"]
		hiddenFacets := make([]string, 0, len(brands))
		for _, brand := range brands {
			if _, ok := brandsToDisplay[brand]; !ok {
				hiddenFacets = append(hiddenFacets, brand)
			}
		}

		nbPages := results.Results[0].NbPages
		props := map[string]interface{}{
			"Hits":         hits,
			"Facets":       brandsToDisplay,
			"Query":        query,
			"Brands":       brands,
			"HiddenFacets": hiddenFacets,
			"Pages":        paginator(pageNumber, nbPages, 3),
			"IsFirstPage":  pageNumber == 0,
			"IsLastPage":   pageNumber == nbPages-1,
			"LastPage":     nbPages - 1,
		}

		// If it's an htmx request, we only render a partial
		if isHtmxRequest {
			testTemplate.ExecuteTemplate(w, "search-form", props)
		} else {
			testTemplate.Execute(w, props)
		}
	})

	// Workaround for https://github.com/bigskysoftware/htmx/issues/1784
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
