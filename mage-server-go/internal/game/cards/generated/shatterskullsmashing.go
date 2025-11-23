package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shatterskull Smashing", NewShatterskullSmashing)
}

// NewShatterskullSmashing creates a Shatterskull Smashing
//   - SORCERY
func NewShatterskullSmashing(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shatterskull Smashing")
	card.ManaCost = ""
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability0)
	return card, nil
}
