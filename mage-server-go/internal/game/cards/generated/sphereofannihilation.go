package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sphere Of Annihilation", NewSphereOfAnnihilation)
}

// NewSphereOfAnnihilation creates a Sphere Of Annihilation
// {X}{B} - ARTIFACT
func NewSphereOfAnnihilation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sphere Of Annihilation")
	card.ManaCost = "{X}{B}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
