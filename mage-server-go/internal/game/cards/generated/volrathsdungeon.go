package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Volraths Dungeon", NewVolrathsDungeon)
}

// NewVolrathsDungeon creates a Volraths Dungeon
// {2}{B}{B} - ENCHANTMENT
func NewVolrathsDungeon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Volraths Dungeon")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - DestroySourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
