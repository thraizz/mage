package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Arcane Artisan", NewArcaneArtisan)
}

// NewArcaneArtisan creates a Arcane Artisan
// {2}{U} - CREATURE
func NewArcaneArtisan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arcane Artisan")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Power = "0"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect(player.getId())
	// card.AddAbility(ability0)
	return card, nil
}
