package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Crystal Skull Isu Spyglass", NewCrystalSkullIsuSpyglass)
}

// NewCrystalSkullIsuSpyglass creates a Crystal Skull Isu Spyglass
// {2}{U}{U} - ARTIFACT
func NewCrystalSkullIsuSpyglass(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Crystal Skull Isu Spyglass")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability0)
	return card, nil
}
