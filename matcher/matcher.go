package matcher

import (
	"wallforfry/esiee-api/ade"
	"wallforfry/esiee-api/aurion"
	adeModels "wallforfry/esiee-api/models/ade"
	"wallforfry/esiee-api/utils"
)

func convertToOldFormat(events []adeModels.EventAde) []adeModels.OldFormat {
	var olds []adeModels.OldFormat
	for _, event := range events {
		olds = append(olds, event.ToOldFormat())
	}
	return olds
}

func GetOldFormatEvents(username string) []adeModels.OldFormat {
	return convertToOldFormat(GetEvents(username))
}

func GetEvents(username string) []adeModels.EventAde {
	var events []adeModels.EventAde

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
