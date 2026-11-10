package cmd

import (
	"encoding/json"
	"fmt"
	"testing"

	cfv0 "github.com/cloudflare/cloudflare-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDropUnmanageableResources(t *testing.T) {
	custom := map[string]interface{}{"id": "a", "kind": "zone", "phase": "http_request_firewall_custom"}
	root := map[string]interface{}{"id": "b", "kind": "root", "phase": "http_request_firewall_managed"}
	managed := map[string]interface{}{"id": "c", "kind": "managed", "phase": "ddos_l7"}

	t.Run("drops managed rulesets", func(t *testing.T) {
		response := []interface{}{custom, managed, root}
		dropUnmanageableResources(&response, "cloudflare_ruleset")
		assert.Equal(t, []interface{}{custom, root}, response)
	})

	t.Run("leaves other resource types alone", func(t *testing.T) {
		response := []interface{}{custom, managed}
		dropUnmanageableResources(&response, "cloudflare_page_rule")
		assert.Len(t, response, 2)
	})

	t.Run("keeps rows with no kind rather than dropping everything", func(t *testing.T) {
		// A positive match on "managed" means an unexpected payload shape loses
		// nothing; a negative match on "not zone" would empty the whole list.
		noKind := map[string]interface{}{"id": "d"}
		response := []interface{}{noKind, managed}
		dropUnmanageableResources(&response, "cloudflare_ruleset")
		assert.Equal(t, []interface{}{noKind}, response)
	})

	t.Run("is idempotent on already filtered data", func(t *testing.T) {
		// The v4 path filters in its own switch before reaching this call.
		response := []interface{}{custom, root}
		dropUnmanageableResources(&response, "cloudflare_ruleset")
		dropUnmanageableResources(&response, "cloudflare_ruleset")
		assert.Equal(t, []interface{}{custom, root}, response)
	})

	t.Run("handles an empty response", func(t *testing.T) {
		response := []interface{}{}
		dropUnmanageableResources(&response, "cloudflare_ruleset")
		assert.Empty(t, response)
	})
}

func TestSortResourcesForStableNaming(t *testing.T) {
	ids := func(response []interface{}) []string {
		out := make([]string, 0, len(response))
		for _, r := range response {
			out = append(out, stringAttr(r, "id"))
		}
		return out
	}

	t.Run("orders by phase", func(t *testing.T) {
		// IDs are deliberately anti-correlated with phase order: sorting by id
		// alone would give [a b c], so this fails if the phase key is ever lost
		// and the comparator silently falls through to the tiebreak.
		response := []interface{}{
			map[string]interface{}{"id": "a", "phase": "http_response_headers_transform"},
			map[string]interface{}{"id": "c", "phase": "http_request_firewall_custom"},
			map[string]interface{}{"id": "b", "phase": "http_request_late_transform"},
		}
		sortResourcesForStableNaming(&response, "cloudflare_ruleset")
		assert.Equal(t, []string{"c", "b", "a"}, ids(response))
	})

	t.Run("breaks ties on id so phase need not be unique", func(t *testing.T) {
		// Several rulesets can share a phase at account scope. The previous
		// sort.Slice was unstable, so these could be renamed between runs.
		response := []interface{}{
			map[string]interface{}{"id": "z", "phase": "http_request_firewall_custom"},
			map[string]interface{}{"id": "m", "phase": "http_request_firewall_custom"},
			map[string]interface{}{"id": "a", "phase": "http_request_firewall_custom"},
		}
		sortResourcesForStableNaming(&response, "cloudflare_ruleset")
		assert.Equal(t, []string{"a", "m", "z"}, ids(response))
	})

	t.Run("is deterministic regardless of input order", func(t *testing.T) {
		// This is the property that actually matters: `generate` and `import`
		// receive the list from different clients and must still agree.
		rows := func() []interface{} {
			return []interface{}{
				map[string]interface{}{"id": "a", "phase": "ddos_l7"},
				map[string]interface{}{"id": "b", "phase": "ddos_l7"},
				map[string]interface{}{"id": "c", "phase": "http_ratelimit"},
			}
		}
		forward := rows()
		reversed := []interface{}{rows()[2], rows()[1], rows()[0]}

		sortResourcesForStableNaming(&forward, "cloudflare_ruleset")
		sortResourcesForStableNaming(&reversed, "cloudflare_ruleset")
		assert.Equal(t, ids(forward), ids(reversed))
	})

	t.Run("leaves other resource types alone", func(t *testing.T) {
		response := []interface{}{
			map[string]interface{}{"id": "z"},
			map[string]interface{}{"id": "a"},
		}
		sortResourcesForStableNaming(&response, "cloudflare_page_rule")
		assert.Equal(t, []string{"z", "a"}, ids(response))
	})

	t.Run("tolerates rows with missing fields", func(t *testing.T) {
		response := []interface{}{
			map[string]interface{}{"id": "b", "phase": "ddos_l7"},
			"not an object",
			map[string]interface{}{"id": "a"},
		}
		assert.NotPanics(t, func() {
			sortResourcesForStableNaming(&response, "cloudflare_ruleset")
		})
		assert.Len(t, response, 3)
	})
}

// TestSortOperatesOnMarshalledRulesetShape pins the contract between the
// cloudflare-go Ruleset struct and sortResourcesForStableNaming.
//
// `generate` sorts jsonStructData after a marshal/unmarshal round-trip of
// []cfv0.Ruleset, whereas `import` sorts the raw API JSON. The helper is only
// genuinely shared if the sort keys survive both routes under the same names:
// if the SDK ever renames or restructures these fields, stringAttr would return
// "" for every row, the generate-side sort would silently become a no-op, and
// the two subcommands would drift apart again with no test failing.
//
// No cassette currently contains more than one non-managed ruleset, so this is
// the only coverage of ordering against the real payload shape.
func TestSortOperatesOnMarshalledRulesetShape(t *testing.T) {
	ids := func(response []interface{}) []string {
		out := make([]string, 0, len(response))
		for _, r := range response {
			out = append(out, stringAttr(r, "id"))
		}
		return out
	}

	// IDs are deliberately anti-correlated with phase order. Sorting by id
	// alone yields [aaa bbb ccc], so if the phase key stopped resolving the
	// comparator would fall through to the tiebreak and produce a different
	// result than the assertion below. Without this the test would pass
	// vacuously on a broken sort.
	//
	// The generate path: typed SDK structs through a JSON round-trip.
	typed := []cfv0.Ruleset{
		{ID: "aaa", Phase: "http_request_late_transform", Kind: "zone"},
		{ID: "ccc", Phase: "http_request_firewall_custom", Kind: "zone"},
		{ID: "bbb", Phase: "http_request_firewall_custom", Kind: "zone"},
	}
	var generateSide []interface{}
	m, err := json.Marshal(typed)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(m, &generateSide))

	// The keys the helper reads must exist after the round-trip. Without this
	// the ordering assertion below could pass vacuously on an unsorted list.
	first, ok := generateSide[0].(map[string]interface{})
	require.True(t, ok)
	require.Contains(t, first, "id", "sort key must survive the SDK round-trip")
	require.Contains(t, first, "phase", "sort key must survive the SDK round-trip")

	// The import path: the same objects as raw API JSON, in a different order.
	var importSide []interface{}
	require.NoError(t, json.Unmarshal([]byte(`[
		{"id":"bbb","kind":"zone","phase":"http_request_firewall_custom"},
		{"id":"aaa","kind":"zone","phase":"http_request_late_transform"},
		{"id":"ccc","kind":"zone","phase":"http_request_firewall_custom"}
	]`), &importSide))

	sortResourcesForStableNaming(&generateSide, "cloudflare_ruleset")
	sortResourcesForStableNaming(&importSide, "cloudflare_ruleset")

	// Phase primary, id tiebreak - and crucially, identical across both shapes.
	assert.Equal(t, []string{"bbb", "ccc", "aaa"}, ids(generateSide))
	assert.Equal(t, ids(generateSide), ids(importSide),
		"generate and import must order the same rulesets identically despite arriving via different clients")
}

// TestImportGenerateParity is the regression test for the reported bug: the two
// subcommands are handed the same zone listing in different orders, and must
// derive the same set of resource addresses from it.
func TestImportGenerateParity(t *testing.T) {
	zoneListing := func() []interface{} {
		return []interface{}{
			map[string]interface{}{"id": "r1", "kind": "zone", "phase": "http_request_firewall_custom"},
			map[string]interface{}{"id": "m1", "kind": "managed", "phase": "ddos_l7"},
			map[string]interface{}{"id": "r2", "kind": "zone", "phase": "http_request_late_transform"},
			map[string]interface{}{"id": "m2", "kind": "managed", "phase": "http_request_firewall_managed"},
			map[string]interface{}{"id": "r3", "kind": "root", "phase": "http_request_origin"},
		}
	}

	addresses := func(response []interface{}) []string {
		dropUnmanageableResources(&response, "cloudflare_ruleset")
		sortResourcesForStableNaming(&response, "cloudflare_ruleset")

		out := make([]string, 0, len(response))
		for i, r := range response {
			out = append(out, fmtAddress(stringAttr(r, "id"), i))
		}
		return out
	}

	generateOrder := zoneListing()

	importOrder := zoneListing()
	// The import path receives the raw API listing, which is not guaranteed to
	// arrive in the same order.
	importOrder[0], importOrder[4] = importOrder[4], importOrder[0]

	gen := addresses(generateOrder)
	imp := addresses(importOrder)

	assert.Equal(t, gen, imp, "import must derive the same addresses as generate")
	assert.Len(t, gen, 3, "managed rulesets must not produce an address")
	assert.NotContains(t, gen, fmtAddress("m1", 0))
}

// fmtAddress mirrors the resource address both subcommands build
// (generate.go and import.go).
func fmtAddress(id string, i int) string {
	return fmt.Sprintf("%s_%s_%d", terraformResourceNamePrefix, id, i)
}
