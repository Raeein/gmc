package poll

import (
	"context"
	"github.com/Raeein/gmc/email"
	"github.com/Raeein/gmc/firestore"
	"github.com/Raeein/gmc/mongodb"
	"github.com/Raeein/gmc/webadvisor"
	"log"
)

type Trigger struct {
	webadvisorService webadvisor.WebAdvisor
	firestoreService  firestore.Service
	emailService      email.Service
	mongoService      mongodb.Logger
}

func New(ws webadvisor.WebAdvisor, fs firestore.Service, es email.Service, ms mongodb.Logger) Trigger {
	return Trigger{
		webadvisorService: ws,
		firestoreService:  fs,
		emailService:      es,
		mongoService:      ms,
	}
}

func (t Trigger) Trigger(ctx context.Context) {
	// Trigger steps
	// 1. Get all watched sections from the watcher service
	// 2. Loop over the sections, checking the available capacity on each
	// 3. If availability is found, use the notifiers to notify the watchers for that section
	// 4. Remove said watchers once successfully notified

	sections, err := t.firestoreService.GetUserSections(ctx)
	if err != nil {
		return
	}

	if len(sections) == 0 {
		log.Println("No watched sections")
		return
	}

	for _, section := range sections {
		available, err := t.webadvisorService.GetAvailableSeats(ctx, section)
		if err != nil {
			return
		}

		log.Printf("%d available seats found for %+v", available, section)

		if available == 0 {
			continue
		}

		users, err := t.firestoreService.GetUsers(ctx, section)
		if err != nil {
			return
		}

		err = t.emailService.Send(section, users...)
		if err != nil {
			return
		}

		if err := t.firestoreService.RemoveUsers(ctx, section); err != nil {
			return
		}
	}
}
