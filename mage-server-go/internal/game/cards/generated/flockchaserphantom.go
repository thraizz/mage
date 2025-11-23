package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flockchaser Phantom", NewFlockchaserPhantom)
}

// NewFlockchaserPhantom creates a Flockchaser Phantom
// {4}{W}{U} - CREATURE
// Flying, Vigilance
func NewFlockchaserPhantom(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flockchaser Phantom")
	card.ManaCost = "{4}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	return card, nil
}
