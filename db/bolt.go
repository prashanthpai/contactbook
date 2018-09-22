package db

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
)

const pageLimit = 10

func New(file string) (DB, error) {

	db, err := storm.Open(file,
		storm.Batch(), storm.Codec(protobuf.Codec),
	)
	if err != nil {
		return nil, err
	}

	return &boltDB{
		db:        db,
		pageLimit: pageLimit,
	}, nil
}

type boltDB struct {
	db        *storm.DB
	pageLimit int
}

func (b *boltDB) Store(entry *Entry) error {
	return b.db.Save(entry)
}

func (b *boltDB) Update(email string, newentry *Entry) error {

	txn, err := b.db.Begin(true)
	if err != nil {
		return err
	}
	defer txn.Rollback()

	var entry Entry
	if err := txn.One("Email", email, &entry); err != nil {
		return err
	}

	if newentry.Email != "" {
		entry.Email = newentry.Email
	}
	if newentry.Name != "" {
		entry.Name = newentry.Name
	}
	if newentry.Phone != "" {
		entry.Phone = newentry.Phone
	}

	if err := txn.Update(&entry); err != nil {
		return err
	}

	return txn.Commit()
}

func (b *boltDB) All(pageNum int) ([]*Entry, error) {

	skip := b.pageLimit * (pageNum - 1)

	var entries []*Entry
	if err := b.db.All(&entries,
		storm.Skip(skip), storm.Limit(b.pageLimit)); err != nil {
		return nil, err
	}

	return entries, nil
}

func (b *boltDB) FindByEmail(email string) (*Entry, error) {

	var entry Entry
	if err := b.db.One("Email", email, &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

func (b *boltDB) FindByName(name string, pageNum int) ([]*Entry, error) {

	skip := b.pageLimit * (pageNum - 1)

	var entries []*Entry
	if err := b.db.Find("Name", name, &entries,
		storm.Skip(skip), storm.Limit(b.pageLimit)); err != nil {
		return nil, err
	}

	return entries, nil
}

func (b *boltDB) Delete(email string) error {

	txn, err := b.db.Begin(true)
	if err != nil {
		return err
	}
	defer txn.Rollback()

	var entry Entry
	if err := txn.One("Email", email, &entry); err != nil {
		return err
	}

	if err := txn.DeleteStruct(&entry); err != nil {
		return err
	}

	return txn.Commit()
}

func (b *boltDB) Close() error {
	return b.db.Close()
}
