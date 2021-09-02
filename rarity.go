package opensea_go

type Trait struct {
	TraitType string  `json:"traitType" bson:"trait_type"`
	Value     string  `json:"value" bson:"value"`
	Percent   float64 `json:"percent,omitempty" bson:"percent,omitempty"`
	Score     float64 `json:"score,omitempty" bson:"score,omitempty"`
}

type Meta struct {
	TokenId    int     `json:"tokenId" bson:"token_id"`
	Attributes []Trait `json:"attributes" bson:"attributes"`
}

type Rarity struct {
	TokenId    int     `json:"tokenId" bson:"token_id"`
	Score      float64 `json:"score" bson:"score"`
	Attributes []Trait `json:"attributes" bson:"attributes"`
}

// RarityScore calculate the rarity score like rarity.tools
// See https://raritytools.medium.com/ranking-rarity-understanding-rarity-calculation-methods-86ceaeb9b98c
//   [Rarity Score for a Trait Value] = 1 / ([Number of Items with that Trait Value] / [Total Number of Items in Collection])
func RarityScore(metas []Meta) []Rarity {
	sumOfType := make(map[string]int)
	sumOfTraits := make(map[string]map[string]int)
	for _, m := range metas {
		for _, t := range m.Attributes {
			sumOfType[t.TraitType] += 1
			if mt, ok := sumOfTraits[t.TraitType]; ok {
				mt[t.Value] += 1
			} else {
				sumOfTraits[t.TraitType] = make(map[string]int)
				sumOfTraits[t.TraitType][t.Value] += 1
			}
		}
	}
	var rarities []Rarity
	for i := 0; i < len(metas); i++ {
		var score float64
		for j := 0; j < len(metas[i].Attributes); j++ {
			traitType := metas[i].Attributes[j].TraitType
			traitValue := metas[i].Attributes[j].Value
			metas[i].Attributes[j].Percent = float64(sumOfTraits[traitType][traitValue]) / float64(sumOfType[traitType])
			metas[i].Attributes[j].Score = float64(sumOfType[traitType]) / float64(sumOfTraits[traitType][traitValue])
			score += metas[i].Attributes[j].Score
		}
		rarities = append(rarities, Rarity{
			TokenId:    metas[i].TokenId,
			Score:      score,
			Attributes: metas[i].Attributes,
		})
	}
	return rarities
}
