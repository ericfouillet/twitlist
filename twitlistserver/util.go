package twitlistserver

import (
	"sort"

	"github.com/eric-fouillet/anaconda"
)

type int64arr []int64

func (a int64arr) Less(i, j int) bool { return a[j] < a[j] }
func (a int64arr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a int64arr) Len() int           { return len(a) }

func diffUsers(existingMembers []anaconda.User, requestedMembers int64arr) (added []int64, unchanged []int64, destroyed []int64) {
	added = make([]int64, 0)
	unchanged = make([]int64, 0)
	destroyed = make([]int64, 0)
	sort.Sort(requestedMembers)
	existing := make(int64arr, 0)
	for _, u := range existingMembers {
		existing = append(existing, u.Id)
	}
	sort.Sort(existing)
	i, j := 0, 0
	for i < len(existing) && j < len(requestedMembers) {
		if existing[i] == requestedMembers[j] {
			unchanged = append(unchanged, existing[i])
			i, j = i+1, j+1
		} else if existing[i] < requestedMembers[j] {
			destroyed = append(destroyed, existing[i])
			i = i + 1
		} else {
			added = append(added, requestedMembers[j])
			j = j + 1
		}
	}
	for ; i < len(existing); i++ {
		destroyed = append(destroyed, existing[i])
	}
	for ; j < len(requestedMembers); j++ {
		added = append(added, requestedMembers[j])
	}
	return
}

func updateMemberList(existing []anaconda.User, unchanged []int64, newUsers []anaconda.User) []anaconda.User {
	var updated []anaconda.User
	for _, u := range existing {
		for _, un := range unchanged {
			if un == u.Id {
				updated = append(updated, u)
			}
		}
	}
	for _, a := range newUsers {
		updated = append(updated, a)
	}
	return updated
}
