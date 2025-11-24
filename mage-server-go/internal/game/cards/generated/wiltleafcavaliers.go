package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wilt Leaf Cavaliers", NewWiltLeafCavaliers)
}

// NewWiltLeafCavaliers creates a Wilt Leaf Cavaliers
// {G/W}{G/W}{G/W} - CREATURE
// Vigilance
func NewWiltLeafCavaliers(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wilt Leaf Cavaliers")
	card.ManaCost = "{G/W}{G/W}{G/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "KNIGHT"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	return card, nil
}