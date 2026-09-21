package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDropUnmanageableResources(t *testing.T) {
	custom := map[string]interface{}{"id": "a", "kind": "zone"}
	root := map[string]interface{}{"id": "b", "kind": "root"}
	managed := map[string]interface{}{"id": "c", "kind": "managed"}

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

	t.Run("handles an empty response", func(t *testing.T) {
		response := []interface{}{}
		dropUnmanageableResources(&response, "cloudflare_ruleset")
		assert.Empty(t, response)
	})
}
