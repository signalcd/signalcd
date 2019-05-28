package boltdb

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/signalcd/signalcd/api"
	"github.com/signalcd/signalcd/signalcd"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/xerrors"
)

type BoltDB struct {
	db *bolt.DB
}

func New() (*BoltDB, func() error, error) {
	// TODO: Make path configurable
	db, err := bolt.Open("./development/data", 0666, nil)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to open bolt db: %w", err)
	}

	err = db.Batch(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(`deployments`))
		if err != nil {
			return err
		}

		for _, d := range api.FakeDeployments {
			key := strconv.Itoa(int(d.Number))
			value, _ := json.Marshal(d)

			if err := bucket.Put([]byte(key), value); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return &BoltDB{db: db}, db.Close, nil
}

func (bdb *BoltDB) List() ([]signalcd.Deployment, error) {
	var ds []signalcd.Deployment

	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`deployments`))
		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var d signalcd.Deployment
			_ = json.Unmarshal(v, &d)
			ds = append(ds, d)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (bdb *BoltDB) Create(pipeline signalcd.Pipeline) (signalcd.Deployment, error) {
	var d signalcd.Deployment

	err := bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`deployments`))
		c := b.Cursor()

		num := 0
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			num++
		}

		d := signalcd.Deployment{
			Number:  int64(num + 1),
			Created: time.Now(),
			Status: signalcd.DeploymentStatus{
				Phase: signalcd.Unknown,
			},
			Pipeline: pipeline,
		}

		key := strconv.Itoa(int(d.Number))
		value, _ := json.Marshal(d)

		return b.Put([]byte(key), value)
	})

	return d, err
}

func (bdb *BoltDB) GetCurrentDeployment() (signalcd.Deployment, error) {
	var value []byte

	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`deployments`))
		c := b.Cursor()

		num := 0
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			i, err := strconv.Atoi(string(k))
			if err != nil {
				return err
			}

			if num < i {
				num = i
				value = v
			}
		}

		return nil
	})
	if err != nil {
		return signalcd.Deployment{}, err
	}

	var d signalcd.Deployment
	err = json.Unmarshal(value, &d)
	return d, err
}
