package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Synthesizer Labship", NewSynthesizerLabship)
}

// NewSynthesizerLabship creates a Synthesizer Labship
// {U} - ARTIFACT
// Flying, Vigilance
func NewSynthesizerLabship(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Synthesizer Labship")
	card.ManaCost = "{U}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"SPACECRAFT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	return card, nil
}
