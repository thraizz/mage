package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ochre Jelly", NewOchreJelly)
}

// NewOchreJelly creates a Ochre Jelly
// {X}{G} - CREATURE
// Trample
func NewOchreJelly(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ochre Jelly")
	card.ManaCost = "{X}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OOZE"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
