// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Commonly used ConnectionHooks.
package database

// Enter: attach course database.
// Exit: detach course database.
func AttachCourse(path string) ConnectionHook {
	return ConnectionHook{
		Enter: func(c *Connection) error {
			return attach(c.con, "course", path)
		},
		Exit: func(c *Connection) error {
			return detach(c.con, "course")
		},
	}
}
