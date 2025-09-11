package json

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/dcaravel/jsonfinder/pkg/config"
)

type searchContext struct {
	searchRe  *regexp.Regexp
	findings  []*JSONItem
	interests map[string]bool
}

func Search(c *config.Config) ([]*JSONItem, error) {
	dataB, err := os.ReadFile(c.FilePath)
	if err != nil {
		return nil, fmt.Errorf("reading file")
	}

	var data any
	err = json.Unmarshal(dataB, &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling file data")
	}

	searchRe, err := regexp.Compile(c.SearchTerm)
	if err != nil {
		return nil, fmt.Errorf("compiling search term")
	}

	interests := map[string]bool{}
	for _, interest := range c.Context {
		interests[interest] = true
	}

	sc := &searchContext{
		searchRe:  searchRe,
		interests: interests,
	}

	dosearch(sc, "", data, 0, &Breadcrumb{})

	return sc.findings, nil
}

func dosearch(sc *searchContext, key any, untypedVal any, depth int, bc *Breadcrumb) {
	switch val := untypedVal.(type) {
	case []any:
		for i, untypedVal := range val {
			bc.Push(i, Slice)
			dosearch(sc, i, untypedVal, depth+1, bc)
			bc.Pop()

		}
	case map[string]any:
		for k, untypedVal := range val {
			bc.Push(k, Map)
			dosearch(sc, k, untypedVal, depth+1, bc)
			bc.Pop()
		}
	case string:
		interestKey := strings.Join(bc.AllStringNoIndex(), " ")
		_, found := sc.interests[interestKey]
		if found {
			bc.AddContextToParent(fmt.Sprintf("%s", key), val)
		}
		if sc.searchRe.MatchString(val) {
			sc.findings = append(sc.findings, &JSONItem{
				Value:             val,
				Breadcrumb:        strings.Join(bc.AllWithContext(), " "),
				BreadcrumbEntries: bc.Entries(),
			})
		}
	case float64:
	case bool:
	default:
		panic(fmt.Sprintf("%v has an UNKONWN type %T", key, untypedVal))
	}
}

// JSONItem represents a single json element, including its full hierarchy.
type JSONItem struct {
	Value             string
	Breadcrumb        string
	BreadcrumbEntries []*BreadcrumbEntry
}
