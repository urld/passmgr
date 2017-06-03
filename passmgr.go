// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package passmgr

// Store provides access to stored Credentials.
type Store interface {
	List() []Subject
	Store(Subject)
	Load(Subject) (Subject, bool)
	Delete(Subject) bool
}

// Subject represents contain information on various secrets for
// a given user name.
type Subject struct {
	Description string
	User        string
	Secrets     map[string]string
}
