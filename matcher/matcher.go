package matcher

import (
	"wallforfry/esiee-api/ade"
	"wallforfry/esiee-api/aurion"
	"wallforfry/esiee-api/utils"
)

func convertToOldFormat(events []ade.EventAde) []ade.OldFormat {
	var olds []ade.OldFormat
	for _, event := range events {
		olds = append(olds, event.ToOldFormat())
	}
	return olds
}

func GetOldFormatEvents(username string) []ade.OldFormat {
	return convertToOldFormat(GetEvents(username))
}

func GetEvents(username string) []ade.EventAde {
	var events []ade.EventAde

	allEvents := ade.GetEvents()
	groups := aurion.GetUserGroups(username)

	for _, event := range allEvents {
		for _, group := range groups {
			if event.Unite != group.Unite {
				continue
			}
			inter := utils.Intersect(event.Trainees, group.Groups)
			if len(inter) == 0 {
				continue
			}
			events = append(events, event)
			//fmt.Printf("%s :: %s :::: %s : %s <== %s len == %d\n", event.Unite, group.Unite, event.Trainees, inter, group.Groups, len(inter))
		}
	}

	return events
}
