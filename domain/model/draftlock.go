// domain/model/draftlock.go
package model

import "sync"

// draftLocks sérialise les lectures-modifications-écritures d'un même
// brouillon : deux requêtes concurrentes sur un module n'écrasent plus la
// modification l'une de l'autre (R-21). Partagé par toutes les instances du
// service d'un même processus — l'API en crée une par requête.
var draftLocks = struct {
	sync.Mutex
	m map[string]*draftLock
}{m: map[string]*draftLock{}}

type draftLock struct {
	sync.Mutex
	users int
}

// lockDraft verrouille le brouillon id et retourne la fonction qui le
// déverrouille ; le verrou est libéré de la table quand plus personne ne
// l'attend.
func lockDraft(id string) (unlock func()) {
	draftLocks.Lock()
	l := draftLocks.m[id]
	if l == nil {
		l = &draftLock{}
		draftLocks.m[id] = l
	}
	l.users++
	draftLocks.Unlock()

	l.Lock()
	return func() {
		l.Unlock()
		draftLocks.Lock()
		if l.users--; l.users == 0 {
			delete(draftLocks.m, id)
		}
		draftLocks.Unlock()
	}
}
