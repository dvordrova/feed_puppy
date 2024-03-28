package utils

type Action int64

const (
	ActionFeed Action = iota + 1
	ActionYum
	ActionVaccine
	ActionWeigh
	ActionShitDetected
	ActionShitRemoved
	ActionPeeDetected
	ActionPeeRemoved
	ActionPlayStart
	ActionPlayEnd
	ActionTeachStart
	ActionTeachEnd
	ActionSleepStart
	ActionSleepEnd
	ActionWalkStart
	ActionWalkEnd
)

var ActionsMap = map[string]Action{
	"action_feed":          ActionFeed,
	"action_yum":           ActionYum,
	"action_vaccine":       ActionVaccine,
	"action_weigh":         ActionWeigh,
	"action_shit_detected": ActionShitDetected,
	"action_shit_removed":  ActionShitRemoved,
	"action_pee_detected":  ActionPeeDetected,
	"action_pee_removed":   ActionPeeRemoved,
	"action_play_start":    ActionPlayStart,
	"action_play_end":      ActionPlayEnd,
	"action_teach_start":   ActionTeachStart,
	"action_teach_end":     ActionTeachEnd,
	"action_sleep_start":   ActionSleepStart,
	"action_sleep_end":     ActionSleepEnd,
	"action_walk_start":    ActionWalkStart,
	"action_walk_end":      ActionWalkEnd,
}
