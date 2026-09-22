package registry

import "testing"

func TestLookupCodexGPT6Fallbacks(t *testing.T) {
	for _, id := range []string{"gpt-6-sol", "gpt-6-luna"} {
		model := LookupStaticModelInfo(id)
		if model == nil || model.ID != id || model.Thinking == nil {
			t.Fatalf("static lookup lost capabilities for %s: %+v", id, model)
		}
	}
}

func TestCodexGPT6Fallbacks(t *testing.T) {
	old := &ModelInfo{ID: "gpt-5.6-luna", ContextLength: 372000}
	remote := &ModelInfo{ID: "gpt-6-sol", ContextLength: 999999}
	models := WithCodexBuiltins([]*ModelInfo{old, remote, {ID: "gpt-5.6-sol"}})
	counts := map[string]int{}
	for _, model := range models {
		counts[model.ID]++
		if model.ID == "gpt-6-sol" && model != remote {
			t.Fatal("replaced authoritative remote definition")
		}
		if model.ID == "gpt-6-luna" {
			if model.ContextLength != 1050000 || model.MaxCompletionTokens != 128000 || len(model.Thinking.Levels) != 6 {
				t.Fatalf("incorrect launch metadata: %+v", model)
			}
		}
	}
	if counts["gpt-6-luna"] != 1 || counts["gpt-6-sol"] != 1 || old.ContextLength != 372000 {
		t.Fatal("missing/duplicate models or mutated predecessor")
	}
	for _, model := range WithCodexBuiltins([]*ModelInfo{old}) {
		if model.ID == "gpt-6-sol" {
			t.Fatal("added Sol to a tier without its predecessor")
		}
	}
	if got := WithCodexBuiltins(models); len(got) != len(models) {
		t.Fatal("repeated refresh added duplicate models")
	}
}
