// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package passmgr implements a secure store for credentials.
package passmgr

// Store provides access to stored credentials.
type Store interface {

	// List retrieves a list of all Subjects known to the store.
	// The Secrets map of the returned Subjects is empty. To retrieve the
	// complete Subject including its secrets, the Load method needs to be
	// used.
	List() []Subject

	// Load looks up a Subject, identified by its User and URL fields.
	// It returns the complete Subject including its secrets and a flag
	// indicating whether the lookup was successful or not.
	Load(Subject) (s Subject, ok bool)

	// Store adds a new Subject to the store, or updates an existing one.
	Store(Subject)

	// Delete removes a subject from the store. It returns false if the
	// Subject to delete could not be found.
	Delete(Subject) bool
}

// Subject contains various secrets for a given user name.
// Usually the URL and User fields are used as unique identifiers.
type Subject struct {
	User    string
	URL     string
	Secrets map[string]string
}
