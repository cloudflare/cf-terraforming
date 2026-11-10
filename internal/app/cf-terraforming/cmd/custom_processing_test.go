package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
		response := []interface{}{
			map[string]interface{}{"id": "c", "phase": "http_response_headers_transform"},
			map[string]interface{}{"id": "a", "phase": "http_request_firewall_custom"},
			map[string]interface{}{"id": "b", "phase": "http_request_late_transform"},
		}
		sortResourcesForStableNaming(&response, "cloudflare_ruleset")
		assert.Equal(t, []string{"a", "b", "c"}, ids(response))
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
