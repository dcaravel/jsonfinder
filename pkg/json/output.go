package json

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/dcaravel/jsonfinder/pkg/config"
)

func PrintAsTable(c *config.Config, items []*JSONItem) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	for _, f := range items {
		// sb := strings.Builder{}
		// for _, e := range f.BreadcrumbEntries {
		// 	sb.WriteString(fmt.Sprintf("%v(%v) ", e.Val, e.Typ))
		// }
		// fmt.Fprintf(tw, "%s\t%s\n", sb.String(), f.Value)
		fmt.Fprintf(tw, "%s\t%s\n", f.Breadcrumb, f.Value)
	}
	tw.Flush()
}

func PrintAsJson(c *config.Config, items []*JSONItem) {
	data := map[string]any{}
	prevIs := map[string]int{}
	for ii, f := range items {
		_ = ii
		var prev any
		prev = data
		path := ""
		for bci := 0; bci < len(f.BreadcrumbEntries)-1; bci++ {
			bc := f.BreadcrumbEntries[bci]
			bcnext := f.BreadcrumbEntries[bci+1]

			switch pt := prev.(type) {
			case map[string]any:
				switch bcnext.Typ {
				case Map:
					key := bc.Val.(string)
					if _, ok := pt[key]; !ok {
						pt[key] = map[string]any{}
					}
					prev = pt[key]
					if len(bc.Context) != 0 {
						for k, v := range bc.Context {
							prev.(map[string]any)[k] = v
						}
					}
					path += fmt.Sprintf("%v ", key)
				case Slice:
					key := bc.Val.(string)
					if _, ok := pt[key]; !ok {
						pt[key] = &[]any{}
					}
					prev = pt[key]
					if len(bc.Context) != 0 {
						for k, v := range bc.Context {
							prev.(map[string]any)[k] = v
						}
					}
					path += fmt.Sprintf("%v ", key)
				}
			case *[]any:
				// The previous element was a slice, so this current one must be 'appended' to the slice
				originalIndex := bc.Val.(int)
				path += fmt.Sprintf("%v ", originalIndex)
				tpath := strings.TrimSpace(path)
				reuse := false
				if val, ok := prevIs[tpath]; ok && val == originalIndex {
					reuse = true
				}

				switch bcnext.Typ {
				case Map:
					if !reuse {
						// The next wants to be placed in a map, so this last item should be map
						t := map[string]any{}
						*pt = append(*pt, t)
						prev = t
						if c.AddIndexes {
							t["_oindex"] = originalIndex
						}
						if len(bc.Context) != 0 {
							for k, v := range bc.Context {
								prev.(map[string]any)[k] = v
							}
						}
					} else {
						// The next item wants to be placed in a map
						prev = (*pt)[len(*pt)-1]
					}
					prevIs[tpath] = originalIndex
					// TODO: is this case possible? A slice whose next item is a slice?
					// may be possible in normal JSON, but not invex
					// case Slice:
					// May not be possible to have slice followed immediately by slice?
					// pt = append(pt, []any{})
				}
			}
		}

		bc := f.BreadcrumbEntries[len(f.BreadcrumbEntries)-1]
		switch pt := prev.(type) {
		case map[string]any:
			pt[bc.Val.(string)] = f.Value
		case *[]any:
			*pt = append(*pt, f.Value)
		}
	}

	dataB, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", dataB)
}
