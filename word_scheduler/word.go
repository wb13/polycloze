// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package word_scheduler

type Word struct {
	Word       string
	New        bool
	Difficulty int // Meaningful only if New
}
