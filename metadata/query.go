package metadata

import (
	"context"
	"fmt"

	"github.com/Jguer/aur"
	"github.com/itchyny/gojq"
	"github.com/ohler55/ojg/oj"
)

const joiner = " or "

type AURQuery struct {
	Needles  []string
	By       aur.By
	Contains bool // if true, search for packages containing the needle, not exact matches
}

// Get returns a list of packages that provide the given search term.
func (a *Client) Get(ctx context.Context, query *AURQuery) ([]*aur.Pkg, error) {
	found := make([]*aur.Pkg, 0, len(query.Needles))
	if len(query.Needles) == 0 {
		return found, nil
	}

	iterFound, errNeedle := a.gojqGetBatch(ctx, query)
	if errNeedle != nil {
		return nil, errNeedle
	}

	found = append(found, iterFound...)

	return found, nil
}

func (a *Client) gojqGetBatch(ctx context.Context, query *AURQuery) ([]*aur.Pkg, error) {
	pattern := ".[] | select("

	for i, searchTerm := range query.Needles {
		if i != 0 {
			pattern += joiner
		}

		bys := toSearchBy(query.By)
		for j, by := range bys {
			if query.Contains {
				pattern += fmt.Sprintf("(.%s // empty | test(%q))", by, searchTerm)
			} else {
				pattern += fmt.Sprintf("(.%s == %q)", by, searchTerm)
			}

			if j != len(bys)-1 {
				pattern += joiner
			}
		}
	}

	pattern += ")"

	if a.debugLoggerFn != nil {
		a.debugLoggerFn("AUR metadata query", pattern)
	}

	parsed, err := gojq.Parse(pattern)
	if err != nil {
		return nil, fmt.Errorf("unable to parse query: %w", err)
	}

	unmarshalledCache, errCache := a.cache(ctx)
	if errCache != nil {
		return nil, errCache
	}

	final := make([]*aur.Pkg, 0, len(query.Needles))
	iter := parsed.RunWithContext(ctx, unmarshalledCache) // or query.RunWithContext
	dedup := make(map[string]bool)

	for pkgMap, ok := iter.Next(); ok; pkgMap, ok = iter.Next() {
		if err, ok := pkgMap.(error); ok {
			return nil, err
		}

		name := pkgMap.(map[string]interface{})["Name"].(string)
		if dedup[name] {
			continue
		}

		dedup[name] = true

		pkg := new(aur.Pkg)

		bValue, err := gojq.Marshal(pkgMap)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal aur package: %w", err)
		}

		errU := oj.Unmarshal(bValue, pkg)
		if errU != nil {
			return nil, fmt.Errorf("unable to unmarshal aur package: %w", errU)
		}

		final = append(final, pkg)
	}

	if a.debugLoggerFn != nil {
		a.debugLoggerFn("AUR metadata query found", len(final))
	}

	return final, nil
}

func toSearchBy(by aur.By) []string {
	switch by {
	case aur.Name:
		return []string{"Name"}
	case aur.NameDesc:
		return []string{"Name", "Description"}
	case aur.None:
		return []string{"Name", "Provides[]?"}
	case aur.Provides:
		return []string{"Provides[]?"}
	case aur.Maintainer:
		return []string{"Maintainer"}
	case aur.Depends:
		return []string{"Depends[]?"}
	case aur.MakeDepends:
		return []string{"MakeDepends[]?"}
	case aur.OptDepends:
		return []string{"OptDepends[]?"}
	case aur.CheckDepends:
		return []string{"CheckDepends[]?"}
	default:
		panic("invalid By")
	}
}
