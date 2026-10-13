package registry

import (
	"testing"
)

func TestDetectChangedProviders_KimiAliases(t *testing.T) {
	oldData := &staticModelsJSON{
		Kimi: []*ModelInfo{{ID: "kimi-k2"}},
	}
	newData := &staticModelsJSON{
		Kimi: []*ModelInfo{{ID: "kimi-k2"}, {ID: "kimi-k3"}},
	}

	changed := detectChangedProviders(oldData, newData)
	expected := map[string]bool{
		"kimi":     false,
		"kimi-ai":  false,
		"kimi.ai":  false,
		"kimi.com": false,
	}

	for _, p := range changed {
		if _, ok := expected[p]; ok {
			expected[p] = true
		}
	}

	for p, found := range expected {
		if !found {
			t.Errorf("expected changed provider %q to be reported, got %v", p, changed)
		}
	}
}

func TestApplyPinnedModels_RestoresDroppedHaiku55(t *testing.T) {
	data := &staticModelsJSON{Claude: []*ModelInfo{{ID: "claude-haiku-4-5-20251001"}}}

	applyPinnedModels(data)

	var pinned *ModelInfo
	for _, model := range data.Claude {
		if model.ID == "claude-haiku-5-5" {
			pinned = model
		}
	}
	if pinned == nil {
		t.Fatal("claude-haiku-5-5 missing after applyPinnedModels")
	}
	if pinned.ContextLength != 1000000 || pinned.MaxCompletionTokens != 128000 {
		t.Fatalf("pinned limits = %d/%d, want 1000000/128000", pinned.ContextLength, pinned.MaxCompletionTokens)
	}
	if len(data.Claude) != 2 {
		t.Fatalf("len(Claude) = %d, want 2", len(data.Claude))
	}
}

func TestApplyPinnedModels_CatalogDefinitionWins(t *testing.T) {
	data := &staticModelsJSON{Claude: []*ModelInfo{{ID: "claude-haiku-5-5", DisplayName: "upstream"}}}

	applyPinnedModels(data)

	if len(data.Claude) != 1 || data.Claude[0].DisplayName != "upstream" {
		t.Fatalf("Claude = %+v, want the single upstream entry", data.Claude)
	}
}

func TestEmbeddedCatalogIncludesPinnedModels(t *testing.T) {
	if err := loadModelsFromBytes(embeddedModelsJSON, "embed"); err != nil {
		t.Fatal(err)
	}
	for _, model := range GetClaudeModels() {
		if model.ID == "claude-haiku-5-5" {
			return
		}
	}
	t.Fatal("embedded catalog does not register claude-haiku-5-5")
}
