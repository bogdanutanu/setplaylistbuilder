package setlist

import "github.com/jm-duarte/setlistfm"

// ExtractMostRecent ... Will return the most recent non-empty setlist from
// a given slice of setlists. If all the setlists in the slice are empty
// (defining empty as having no songs), an empty setlist will be returned
func ExtractMostRecent(setlists []setlistfm.Setlist) setlistfm.Setlist {
	result := setlistfm.Setlist{}
	for i, setlist := range setlists {
		for _, set := range setlist.Sets.Set {
			if len(set.Song) > 0 {
				return setlists[i]
			}
		}
	}
	return result
}
