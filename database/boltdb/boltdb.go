package boltdb

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/signalcd/signalcd/signalcd"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/xerrors"
)

const (
	bucketDeployments = `deployments`
	bucketPipelines   = `pipelines`
)

// BoltDB has a connection to the database and implements the needed interfaces.
type BoltDB struct {
	db *bolt.DB
}

// New creates a new BoltDB instance
func New(path string) (*BoltDB, func() error, error) {
	// TODO: Make path configurable
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to open bolt db: %w", err)
	}

	err = db.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketDeployments))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(bucketPipelines))
		if err != nil {
			return err
		}

		return nil
	})

	return &BoltDB{db: db}, db.Close, err
}

// ListDeployments lists all Deployments
func (bdb *BoltDB) ListDeployments() ([]signalcd.Deployment, error) {
	var ds []signalcd.Deployment

	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketDeployments))
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

	sort.Slice(ds, func(i, j int) bool {
		return ds[j].Created.Before(ds[i].Created)
	})

	return ds, nil
}

// CreateDeployment creates a new Deployment from a Pipeline
func (bdb *BoltDB) CreateDeployment(pipeline signalcd.Pipeline) (signalcd.Deployment, error) {
	var d signalcd.Deployment

	err := bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketDeployments))
		c := b.Cursor()

		num := 0
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			num++
		}

		d = signalcd.Deployment{
			Number:  int64(num + 1),
			Created: time.Now().UTC(),
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

func (bdb *BoltDB) UpdateDeploymentStatus(deploymentNumber int64, step int64, agent string, phase signalcd.Phase) (signalcd.Deployment, error) {
	var d signalcd.Deployment

	err := bdb.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucketDeployments))
		key := []byte(strconv.Itoa(int(deploymentNumber)))
		value := b.Get(key)

		if value == nil {
			return fmt.Errorf("deployment not found")
		}

		if err := json.Unmarshal(value, &d); err != nil {
			return err
		}

		status := d.Status[agent]
		if status == nil {
			status = &signalcd.Status{}
		}

		if int64(len(status.Steps)) == step {
			status.Steps = append(status.Steps, signalcd.StepStatus{
				Phase:    phase,
				ExitCode: 0,
				Started:  time.Now().UTC(),
				Stopped:  nil,
			})
		} else {
			status.Steps[step].Phase = phase
			if phase == signalcd.Success || phase == signalcd.Failure || phase == signalcd.Killed {
				now := time.Now().UTC()
				status.Steps[step].Stopped = &now
			}
		}

		if d.Status == nil {
			d.Status = map[string]*signalcd.Status{}
		}

		d.Status[agent] = status

		value, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("failed to marshal Deployment after updating status: %w", err)
		}

		return b.Put(key, value)
	})

	return d, err
}

// GetCurrentDeployment gets the current Deployment
func (bdb *BoltDB) GetCurrentDeployment() (signalcd.Deployment, error) {
	var value []byte

	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketDeployments))
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

	if len(value) == 0 {
		return signalcd.Deployment{}, nil
	}

	var d signalcd.Deployment
	err = json.Unmarshal(value, &d)
	return d, err
}

// SaveStepLogs saves the logs for a Deployment step by its number
func (bdb *BoltDB) SaveStepLogs(ctx context.Context, deployment, step int64, logs []byte) error {
	var d signalcd.Deployment

	return bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketDeployments))
		key := []byte(strconv.Itoa(int(deployment)))
		value := b.Get(key)

		if err := json.Unmarshal(value, &d); err != nil {
			return err
		}

		if int64(len(d.Pipeline.Steps)) < step {
			return fmt.Errorf("step %d does not exist", step)
		}

		if d.Pipeline.Steps[step].Status == nil {
			d.Pipeline.Steps[step].Status = &signalcd.Status{}
		}

		d.Pipeline.Steps[step].Status.Logs = logs

		value, err := json.Marshal(d)
		if err != nil {
			return err
		}

		return b.Put(key, value)
	})
}

// GetPipeline gets a Pipeline by its ID
func (bdb *BoltDB) GetPipeline(id string) (signalcd.Pipeline, error) {
	var p signalcd.Pipeline

	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketPipelines))
		bytes := b.Get([]byte(id))

		if err := json.Unmarshal(bytes, &p); err != nil {
			return err
		}
		return nil
	})

	return p, err
}

// ListPipelines returns a list of Pipelines
func (bdb *BoltDB) ListPipelines() ([]signalcd.Pipeline, error) {
	var pipelines []signalcd.Pipeline

	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketPipelines))
		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var p signalcd.Pipeline
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			pipelines = append(pipelines, p)
		}

		return nil
	})

	sort.Slice(pipelines, func(i, j int) bool {
		return pipelines[j].Created.Before(pipelines[i].Created)
	})

	return pipelines, err
}

// CreatePipeline saves a Pipeline and returns the saved Pipeline
func (bdb *BoltDB) CreatePipeline(p signalcd.Pipeline) (signalcd.Pipeline, error) {
	p.ID = uuid.New().String()
	p.Created = time.Now()

	err := bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketPipelines))

		key := p.ID
		value, _ := json.Marshal(p)

		return b.Put([]byte(key), value)
	})

	return p, err
}
