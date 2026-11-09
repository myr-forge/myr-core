//go:build !unix

// adapters/out/localstorage/lock_other.go — sans verrou entre processus hors
// Unix : myr-core ne s'exécute que sur un serveur Linux, ce cas ne concerne que
// le poste de développement (verrou interne au processus seulement).
package localstorage

import "os"

func lockFile(*os.File, bool) error { return nil }

func unlockFile(*os.File) error { return nil }
