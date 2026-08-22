package editor

import (
	"time"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// autoRegenerationDebounce is how long editing must pause before a
// non-preview-affecting change triggers a regeneration.
const autoRegenerationDebounce = 300 * time.Millisecond

// RegenerationDecisionService decides when the live preview regenerates. It is
// stateless: the snapshots and the frame time all arrive as arguments.
type RegenerationDecisionService struct{}

func NewRegenerationDecisionService() IRegenerationDecisionService {
	return &RegenerationDecisionService{}
}

// DecideRegeneration resolves one frame of the regeneration state machine.
//
// Preview-affecting changes (player/zone counts, topology and connection
// settings) regenerate immediately so the preview tracks the control. Every
// other change is debounced and only regenerates once editing has paused for
// autoRegenerationDebounce, so dragging a slider does not regenerate per frame.
func (this *RegenerationDecisionService) DecideRegeneration(
	request dtos.RegenerationDecisionRequestDto) dtos.RegenerationDecisionDto {
	// First generation: populate the preview immediately on startup.
	if request.Previous == nil {
		return dtos.RegenerationDecisionDto{Regenerate: true, NextStateAction: dtos.NextStateLeave}
	}

	// Nothing changed since the last generation -> cancel any pending debounce.
	if request.Previous.EqualsIgnoringManualEdits(request.Current) {
		return dtos.RegenerationDecisionDto{NextStateAction: dtos.NextStateClear}
	}

	if request.Previous.LayoutDefiningOptionsChanged(request.Current) {
		return dtos.RegenerationDecisionDto{Regenerate: true, NextStateAction: dtos.NextStateClear}
	}

	// Still moving: (re)arm the debounce and ask to be woken when it is due.
	if request.Next == nil || !request.Next.EqualsIgnoringManualEdits(request.Current) {
		return dtos.RegenerationDecisionDto{
			NextStateAction: dtos.NextStateSetFromCurrent,
			RedrawAt:        request.Now.Add(autoRegenerationDebounce),
			ScheduleRedraw:  true,
		}
	}

	// Stable since the last frame; keep waiting until due.
	if request.Now.Before(request.DebounceDueAt) {
		return dtos.RegenerationDecisionDto{
			NextStateAction: dtos.NextStateLeave,
			RedrawAt:        request.DebounceDueAt,
			ScheduleRedraw:  true,
		}
	}

	// Editing paused long enough -> regenerate now.
	return dtos.RegenerationDecisionDto{Regenerate: true, NextStateAction: dtos.NextStateClear}
}

// DecideManualEditReapplication compares the live state against the state that
// produced the last generation, so it must be called before that generation is
// snapshotted.
//
// Manual edits are dropped when a layout-defining option changed, because the
// hand-made layout no longer describes the regenerated map.
func (this *RegenerationDecisionService) DecideManualEditReapplication(
	previous, current *editor_state_dto.EditorStateDto) dtos.ManualEditDecisionDto {
	if !current.HasManualEdits() {
		return dtos.ManualEditDecisionDto{}
	}

	if previous == nil {
		return dtos.ManualEditDecisionDto{ReapplyWithCastleChanges: &editor_state_model.CastleSettingChanges{}}
	}

	if previous.LayoutDefiningOptionsChanged(current) {
		return dtos.ManualEditDecisionDto{}
	}

	return dtos.ManualEditDecisionDto{ReapplyWithCastleChanges: new(previous.DiffCastleSettings(current))}
}
