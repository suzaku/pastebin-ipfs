package index

import (
	"errors"
	"math/rand"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

var _pref = cid.Prefix{
	Version:  1,
	Codec:    cid.Raw,
	MhType:   mh.SHA2_256,
	MhLength: -1, // default length
}

func TestIndex_Exist(t *testing.T) {
	// And then feed it some data
	os.RemoveAll("/tmp/pst/idx")
	idx, _ := NewIndex("/tmp/pst/idx")

	for j := 0; j < 100; j++ {
		t.Run("parallel-"+cast.ToString(j), func(t *testing.T) {
			t.Parallel()
			keys := make([]string, 0, 1000)
			for i := 0; i < 10000; i++ {
				c, _ := _pref.Sum([]byte(uuid.NewString()))
				cid := c.String()
				idx.SetExist(cid, ObjectTypeFile)
				keys = append(keys, cid)
			}

			err := testIdxCheckExist(idx, keys)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIndex_Filter(t *testing.T) {
	// And then feed it some data
	os.RemoveAll("/tmp/pst/idx")
	idx, _ := NewIndex("/tmp/pst/idx")

	for j := 0; j < 1; j++ {
		t.Run("parallel-"+cast.ToString(j), func(t *testing.T) {
			// t.Parallel()
			keys := make([]string, 0, 1000)

			for i := 0; i < 1000; i++ {
				c, _ := _pref.Sum([]byte(uuid.NewString()))
				cid := c.String()
				if err := idx.SetExist(cid, ObjectTypeFile); err != nil {
                    t.Fatal(err)
                }
				keys = append(keys, cid)
			}

			testIdxCheckFilter(t, idx, keys)
		})
	}
}

func testIdxCheckExist(idx *Index, keys []string) error {
	for _, cid := range keys {
		ok, err := idx.Exist(cid)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("key not exist")
		}
	}
	return nil
}

func testIdxCheckFilter(t *testing.T, idx *Index, keys []string) {
	i := rand.Intn(100) + 1
	keys = keys[len(keys)-i:]
	findKeys, err := idx.FilterFileCid(i)
	if err != nil {
		t.Fatal(err)
	}
    assert.Equal(t, i, len(findKeys))
	for i := range keys {
		assert.Equal(t, keys[i], findKeys[i])
	}
}
