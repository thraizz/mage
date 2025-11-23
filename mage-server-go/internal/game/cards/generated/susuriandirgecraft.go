package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Susurian Dirgecraft", NewSusurianDirgecraft)
}

// NewSusurianDirgecraft creates a Susurian Dirgecraft
// {4}{B} - ARTIFACT
// Flying
func NewSusurianDirgecraft(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Susurian Dirgecraft")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"SPACECRAFT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
