package data

import bolt "go.etcd.io/bbolt"

var DB *bolt.DB

func OpenDB(path string) error {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func InitDB() error {
	return DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("channels"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("messages"))
		if err != nil {
			return err
		}
		return nil
	})
}

func WriteChannel(key []byte, value []byte) error {
	return DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("channels"))
		if err != nil {
			return err
		}
		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})
}

func ReadChannel(key []byte) (*Channel, error) {
	c := &Channel{}
	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("channels"))
		v := b.Get(key)
		if v == nil {
			c = nil
			return nil
		}
		err := c.Decode(v)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func WriteMessage(key []byte, value []byte) error {
	return DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("messages"))
		if err != nil {
			return err
		}
		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})
}

func ReadMessage(key []byte) (*Message, error) {
	m := &Message{}
	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("messages"))
		v := b.Get(key)
		if v == nil {
			m = nil
			return nil
		}
		err := m.Decode(v)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}
