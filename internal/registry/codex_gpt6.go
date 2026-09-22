package registry

// Keep launch-day models available while the remote catalog catches up. A remote
// definition takes precedence, and each replacement follows its predecessor's tier.
func withCodexGPT6Fallbacks(models []*ModelInfo) []*ModelInfo {
	for _, tier := range []string{"sol", "luna"} {
		id := "gpt-6-" + tier
		var predecessor *ModelInfo
		found := false
		for _, model := range models {
			if model == nil {
				continue
			}
			if model.ID == id {
				found = true
			}
			if model.ID == "gpt-5.6-"+tier {
				predecessor = model
			}
		}
		if found || predecessor == nil {
			continue
		}
		model := cloneModelInfo(predecessor)
		model.ID, model.Version = id, id
		model.DisplayName = "GPT-6 " + map[string]string{"sol": "Sol", "luna": "Luna"}[tier]
		model.Description = model.DisplayName
		model.Created = 1790100733
		model.ContextLength = 1050000
		model.MaxCompletionTokens = 128000
		model.Thinking = &ThinkingSupport{Levels: []string{"none", "low", "medium", "high", "xhigh", "max"}}
		models = append(models, model)
	}
	return models
}
