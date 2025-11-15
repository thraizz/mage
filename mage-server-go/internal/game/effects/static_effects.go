package effects

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// SimplePTBoostEffect applies a flat power/toughness modifier to creatures controlled by a player.
type SimplePTBoostEffect struct {
	id           string
	sourceID     string
	controllerID string
	powerDelta   int
	toughDelta   int
	includeSelf  bool
}

// NewSimplePTBoostEffect creates a new static buff effect.
func NewSimplePTBoostEffect(sourceID, controllerID string, powerDelta, toughDelta int, includeSelf bool) *SimplePTBoostEffect {
	source := strings.TrimSpace(sourceID)
	controller := strings.TrimSpace(controllerID)
	seed := fmt.Sprintf("%s|%s|%d|%d|%t", source, controller, powerDelta, toughDelta, includeSelf)
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(seed)).String()

	return &SimplePTBoostEffect{
		id:           id,
		sourceID:     source,
		controllerID: controller,
		powerDelta:   powerDelta,
		toughDelta:   toughDelta,
		includeSelf:  includeSelf,
	}
}

// ID returns the unique identifier.
func (e *SimplePTBoostEffect) ID() string {
	return e.id
}

// Layer identifies the layer in which the effect applies.
func (e *SimplePTBoostEffect) Layer() Layer {
	return LayerPowerToughness
}

// AppliesTo determines whether the snapshot should receive the modification.
func (e *SimplePTBoostEffect) AppliesTo(snapshot *Snapshot) bool {
	if snapshot == nil {
		return false
	}
	if snapshot.ControllerID != e.controllerID {
		return false
	}
	if !snapshot.HasType("creature") {
		return false
	}
	if !e.includeSelf && snapshot.CardID == e.sourceID {
		return false
	}
	return snapshot.HasBasePower && snapshot.HasBaseTough
}

// Apply mutates the snapshot.
func (e *SimplePTBoostEffect) Apply(snapshot *Snapshot) {
	if snapshot == nil {
		return
	}
	if snapshot.HasBasePower {
		snapshot.Power += e.powerDelta
	}
	if snapshot.HasBaseTough {
		snapshot.Toughness += e.toughDelta
	}
}
