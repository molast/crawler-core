package buckets

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

// A DB is a bolt database with convenience methods for working with buckets.
//
// A DB embeds the exposed bolt.DB methods.
type DB struct {
	*bolt.DB
	path string
}

// Open 在指定路径创建/打开存储桶数据库。
func Open(path string) (*DB, error) {
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err := bolt.Open(path, 0600, config)
	if err != nil {
		return nil, fmt.Errorf("couldn't open %s: %s", path, err)
	}
	return &DB{db, path}, nil
}

// New 创建/打开一个命名的存储桶。
func (db *DB) New(name []byte) (*Bucket, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Bucket{db, name}, nil
}

// Delete 删除指定的存储桶。
func (db *DB) Delete(name []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(name)
	})
}

/* -- ITEM -- */

// An Item holds a key/value pair.
type Item struct {
	Key   []byte
	Value []byte
}

/* -- BUCKET-- */

// Bucket 数据库内的键/值对的集合。
type Bucket struct {
	db   *DB
	Name []byte
}

// Put inserts value `v` with key `k`.
func (bk *Bucket) Put(k, v []byte) error {
	return bk.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bk.Name).Put(k, v)
	})
}

// PutNX (put-if-not-exists) inserts value `v` with key `k`
// if key doesn't exist.
func (bk *Bucket) PutNX(k, v []byte) error {
	v, err := bk.Get(k)
	if v != nil || err != nil {
		return err
	}
	return bk.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bk.Name).Put(k, v)
	})
}

// Insert 会遍历一组 k/v 对，将每个项目作为单个事务的一部分放入
// 存储桶中。
// 对于大量插入，请务必对项目进行预排序（按键按字节排序），
// 这将显著降低插入时间和存储成本。
func (bk *Bucket) Insert(items []struct{ Key, Value []byte }) error {
	return bk.db.Update(func(tx *bolt.Tx) error {
		for _, item := range items {
			_ = tx.Bucket(bk.Name).Put(item.Key, item.Value)
		}
		return nil
	})
}

// InsertNX 插入数据，如果存在不处理，不存在则插入
func (bk *Bucket) InsertNX(items []struct{ Key, Value []byte }) error {
	return bk.db.Update(func(tx *bolt.Tx) error {
		for _, item := range items {
			v, _ := bk.Get(item.Key)
			if v == nil {
				_ = tx.Bucket(bk.Name).Put(item.Key, item.Value)
			}
		}
		return nil
	})
}

// Delete removes key `k`.
func (bk *Bucket) Delete(k []byte) error {
	return bk.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bk.Name).Delete(k)
	})
}

// Get retrieves the value for key `k`.
func (bk *Bucket) Get(k []byte) (value []byte, err error) {
	err = bk.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(bk.Name).Get(k)
		if v != nil {
			value = make([]byte, len(v))
			copy(value, v)
		}
		return nil
	})
	return value, err
}

// Count 获取桶内数据数量
func (bk *Bucket) Count() (int, error) {
	var count int
	err := bk.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bk.Name)
		if b == nil {
			return fmt.Errorf("bucket %q not found", bk.Name)
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				count++
			}
		}
		return nil
	})
	return count, err
}

// Clear 清除桶内数据
func (bk *Bucket) Clear() error {
	return bk.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bk.Name)
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if err := c.Delete(); err != nil {
				return err
			}
		}
		return nil
	})
}

// ClearAndCompact 清空指定 bucket，并压缩数据库
func (bk *Bucket) ClearAndCompact() error {
	tmpPath := bk.db.path + ".tmp"
	newDB, err := bolt.Open(tmpPath, 0600, nil)
	if err != nil {
		return err
	}
	// 遍历原数据库，复制所有 bucket
	err = bk.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			return newDB.Update(func(newTx *bolt.Tx) error {
				newB, xErr := newTx.CreateBucketIfNotExists(name)
				if xErr != nil {
					return xErr
				}
				// 如果是目标 bucket，跳过数据复制，只保留结构
				if string(name) == string(bk.Name) {
					return nil
				}
				// 否则复制全部数据
				return b.ForEach(func(k, v []byte) error {
					return newB.Put(k, v)
				})
			})
		})
	})
	if err != nil {
		_ = newDB.Close()
		return err
	}
	if err = newDB.Close(); err != nil {
		return err
	}
	if err = bk.db.Close(); err != nil {
		return err
	}
	if err = os.Rename(tmpPath, bk.db.path); err != nil {
		return err
	}
	newBoltDB, err := bolt.Open(bk.db.path, 0600, nil)
	if err != nil {
		return err
	}
	bk.db.DB = newBoltDB
	return nil
}

// Items returns a slice of key/value pairs.  Each k/v pair in the slice
// is of type Item (`struct{ Key, Value []byte }`).
func (bk *Bucket) Items() (items []Item, err error) {
	return items, bk.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bk.Name).Cursor()
		var key, value []byte
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				key = make([]byte, len(k))
				copy(key, k)
				value = make([]byte, len(v))
				copy(value, v)
				items = append(items, Item{key, value})
			}
		}
		return nil
	})
}

// ListItems 分页获取数据
func (bk *Bucket) ListItems(offset int, limit int) (items []Item, err error) {
	err = bk.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bk.Name).Cursor()
		count := 0
		skipped := 0
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if skipped < offset {
				skipped++
				continue
			}
			if count >= limit {
				break
			}
			if v != nil {
				key := make([]byte, len(k))
				copy(key, k)
				value := make([]byte, len(v))
				copy(value, v)
				items = append(items, Item{key, value})
				count++
			}
		}
		return nil
	})
	return items, err
}

func (bk *Bucket) MapListItems(do func(k, v []byte) error, offset int, limit int) error {
	return bk.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bk.Name).Cursor()
		count := 0
		skipped := 0
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if skipped < offset {
				skipped++
				continue
			}
			if count >= limit {
				break
			}
			_ = do(k, v)
			count++
		}
		return nil
	})
}

// PrefixItems returns a slice of key/value pairs for all keys with
// a given prefix.  Each k/v pair in the slice is of type Item
// (`struct{ Key, Value []byte }`).
func (bk *Bucket) PrefixItems(pre []byte) (items []Item, err error) {
	err = bk.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bk.Name).Cursor()
		var key, value []byte
		for k, v := c.Seek(pre); bytes.HasPrefix(k, pre); k, v = c.Next() {
			if v != nil {
				key = make([]byte, len(k))
				copy(key, k)
				value = make([]byte, len(v))
				copy(value, v)
				items = append(items, Item{key, value})
			}
		}
		return nil
	})
	return items, err
}

// RangeItems returns a slice of key/value pairs for all keys within
// a given range.  Each k/v pair in the slice is of type Item
// (`struct{ Key, Value []byte }`).
func (bk *Bucket) RangeItems(min []byte, max []byte) (items []Item, err error) {
	err = bk.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bk.Name).Cursor()
		var key, value []byte
		for k, v := c.Seek(min); isBefore(k, max); k, v = c.Next() {
			if v != nil {
				key = make([]byte, len(k))
				copy(key, k)
				value = make([]byte, len(v))
				copy(value, v)
				items = append(items, Item{key, value})
			}
		}
		return nil
	})
	return items, err
}

// Map applies `do` on each key/value pair.
func (bk *Bucket) Map(do func(k, v []byte) error) error {
	return bk.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bk.Name).ForEach(do)
	})
}

// MapPrefix applies `do` on each k/v pair of keys with prefix.
func (bk *Bucket) MapPrefix(do func(k, v []byte) error, pre []byte) error {
	return bk.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bk.Name).Cursor()
		for k, v := c.Seek(pre); bytes.HasPrefix(k, pre); k, v = c.Next() {
			_ = do(k, v)
		}
		return nil
	})
}

// MapRange applies `do` on each k/v pair of keys within range.
func (bk *Bucket) MapRange(do func(k, v []byte) error, min, max []byte) error {
	return bk.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bk.Name).Cursor()
		for k, v := c.Seek(min); isBefore(k, max); k, v = c.Next() {
			_ = do(k, v)
		}
		return nil
	})
}

// NewPrefixScanner initializes a new prefix scanner.
func (bk *Bucket) NewPrefixScanner(pre []byte) *PrefixScanner {
	return &PrefixScanner{bk.db, bk.Name, pre}
}

// NewRangeScanner initializes a new range scanner.  It takes a `min` and a
// `max` key for specifying the range paramaters.
func (bk *Bucket) NewRangeScanner(min, max []byte) *RangeScanner {
	return &RangeScanner{bk.db, bk.Name, min, max}
}
