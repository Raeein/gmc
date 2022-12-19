package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"github.com/Raeein/gmc"
	"google.golang.org/api/option"
)

type Service struct {
	firestore            *firestore.Client
	sectionsCollectionID string
	usersCollectionID    string
}

type UserEntry struct {
	User      gmc.User `json:"watcher"`
	SectionID string   `json:"sectionID"`
}

func New(ctx context.Context, projectID, sectionsCollectionID, usersCollectionID, credPath string) (Service, error) {

	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credPath))
	if err != nil {
		return Service{}, err
	}
	return Service{client, sectionsCollectionID, usersCollectionID}, nil
}

func (f Service) AddWatcher(ctx context.Context, section gmc.Section, watcher gmc.User) error {
	// Steps:
	// 1. Retrieve the section document, creating it if it doesn't exist
	// 2. Inspect the current watchers. If the new watcher is already a watcher, stop and return a nil error
	// 3. Append the new watcher to the watchers array
	// 4. Update the document in the collection

	documents, err := f.firestore.Collection(f.sectionsCollectionID).Where("Code", "==", section.Code).Where("Term", "==", section.Term).Where("Course.Code", "==", section.Course.Code).Where("Course.Department", "==", section.Course.Department).Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("failed to get matching section documents: %w", err)
	}

	if len(documents) > 1 {
		return errors.New("more than one matching document found, expected 0 or 1")
	}

	var sectionID string
	if len(documents) == 0 {
		ref, _, err := f.firestore.Collection(f.sectionsCollectionID).Add(ctx, section)
		if err != nil {
			return fmt.Errorf("failed to add %+v to collection: %w", section, err)
		}
		sectionID = ref.ID
	} else {
		sectionID = documents[0].Ref.ID
	}

	documents, err = f.firestore.Collection(f.usersCollectionID).Where("SectionID", "==", sectionID).Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("failed to get matching watcher documents: %w", err)
	}

	for _, document := range documents {
		var firestoreWatcher UserEntry
		err := document.DataTo(&firestoreWatcher)
		if err != nil {
			return fmt.Errorf("failed to deserialize watcher: %w", err)
		}

		if firestoreWatcher.User.Email == watcher.Email {
			// Service already watching this section, nothing to do
			return nil
		}
	}

	newWatcher := UserEntry{User: watcher, SectionID: sectionID}
	_, _, err = f.firestore.Collection(f.usersCollectionID).Add(ctx, newWatcher)
	if err != nil {
		return fmt.Errorf("failed to write new watcher to collection: %w", err)
	}

	return nil
}

func (f Service) GetWatchedSections(ctx context.Context) ([]gmc.Section, error) {
	documents, err := f.firestore.Collection(f.sectionsCollectionID).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all documents in sections collection: %w", err)
	}

	var results []gmc.Section
	for _, document := range documents {
		var result gmc.Section
		err = document.DataTo(&result)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize document: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (f Service) GetWatchers(ctx context.Context, section gmc.Section) ([]gmc.User, error) {
	documents, err := f.firestore.Collection(f.sectionsCollectionID).Where("Code", "==", section.Code).Where("Term", "==", section.Term).Where("Course.Code", "==", section.Course.Code).Where("Course.Department", "==", section.Course.Department).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get matching section documents: %w", err)
	}

	// sanity check, we should never have more than one matching document
	if len(documents) > 1 {
		return nil, errors.New("more than one matching document found, expected 0 or 1")
	}

	if len(documents) == 0 {
		return nil, errors.New("section not found in firestore")
	}

	sectionID := documents[0].Ref.ID

	documents, err = f.firestore.Collection(f.usersCollectionID).Where("SectionID", "==", sectionID).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get matching watcher documents: %w", err)
	}

	var results []gmc.User
	for _, document := range documents {
		var result UserEntry
		err = document.DataTo(&result)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize document: %w", err)
		}

		results = append(results, result.User)
	}

	return results, nil
}

func (f Service) RemoveWatchers(ctx context.Context, section gmc.Section) error {
	documents, err := f.firestore.Collection(f.sectionsCollectionID).Where("Code", "==", section.Code).Where("Term", "==", section.Term).Where("Course.Code", "==", section.Course.Code).Where("Course.Department", "==", section.Course.Department).Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("failed to get matching section documents: %w", err)
	}

	// sanity check, we should never have more than one matching document
	if len(documents) > 1 {
		return errors.New("more than one matching document found, expected 0 or 1")
	}

	if len(documents) == 0 {
		return errors.New("section not found in firestore")
	}

	sectionID := documents[0].Ref.ID
	_, err = documents[0].Ref.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete section: %w", err)
	}

	documents, err = f.firestore.Collection(f.usersCollectionID).Where("SectionID", "==", sectionID).Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("failed to get matching watcher documents: %w", err)
	}

	for _, document := range documents {
		_, err = document.Ref.Delete(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete watcher: %w", err)
		}
	}

	return nil
}
