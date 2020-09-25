package matcher

import (
	"wallforfry/esiee-api/ade"
	"wallforfry/esiee-api/aurion"
	"wallforfry/esiee-api/pkg/event"
	"wallforfry/esiee-api/utils"
)

func convertToOldFormat(events []event.Event) []ade.OldFormat {
	var olds []ade.OldFormat
	for _, e := range events {
		olds = append(olds, ade.FromNewFormat(e))
	}
	return olds
}

func GetOldFormatEvents(username string) []ade.OldFormat {
	return convertToOldFormat(GetEvents(username))
}

func GetEvents(username string) []event.Event {
	var events []event.Event

	allEvents := ade.GetEvents()
	groups := aurion.GetUserGroups(username)

	for _, e := range allEvents {
		for _, group := range groups {
			if e.Unite != group.Unite {
				continue
			}
			inter := utils.Intersect(e.Trainees, group.Groups)
			if len(inter) == 0 {
				continue
			}
			events = append(events, e)
			//fmt.Printf("%s :: %s :::: %s : %s <== %s len == %d\n", e.Unite, group.Unite, e.Trainees, inter, group.Groups, len(inter))
		}
	}

	return events
}
