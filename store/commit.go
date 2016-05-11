package store

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/disorganizer/brig/id"
	"github.com/disorganizer/brig/store/wire"
	"github.com/gogo/protobuf/proto"
	"github.com/jbenet/go-multihash"
)

const (
	// ChangeInvalid indicates a bug.
	ChangeInvalid = iota

	// ChangeAdd means the file was added (initially or after ChangeRemove)
	ChangeAdd

	// ChangeModify indicates a content modification.
	ChangeModify

	// ChangeMove indicates that a file's path changed.
	ChangeMove

	// ChangeRemove indicates that the file was deleted.
	// Old versions might still be accessible from the history.
	ChangeRemove
)

// ChangeType describes the nature of a change.
type ChangeType byte

var changeTypeToString = map[ChangeType]string{
	ChangeInvalid: "invalid",
	ChangeAdd:     "added",
	ChangeModify:  "modified",
	ChangeRemove:  "removed",
	ChangeMove:    "moved",
}

var stringToChangeType = map[string]ChangeType{
	"invalid":  ChangeInvalid,
	"added":    ChangeAdd,
	"modified": ChangeModify,
	"removed":  ChangeRemove,
	"moved":    ChangeMove,
}

var (
	// ErrNoChange means that nothing changed between two versions (of a file)
	ErrNoChange = fmt.Errorf("Nothing changed between the given versions")
)

// String formats a changetype to a human readable verb in past tense.
func (c *ChangeType) String() string {
	return changeTypeToString[*c]
}

// UnmarshalJSON reads a json string and tries to convert it to a ChangeType.
func (c *ChangeType) Unmarshal(data []byte) error {
	var ok bool
	*c, ok = stringToChangeType[string(data)]
	if !ok {
		return fmt.Errorf("Bad change type: %v", string(data))
	}

	return nil
}

// Commit groups a change set
type Commit struct {
	// Commit message (might be auto-generated)
	Message string

	// Author is the id of the committer.
	Author id.ID

	// Time at this commit was conceived.
	ModTime time.Time

	// Set of files that were changed.
	Changes map[*File]*Checkpoint

	// Hash of this commit (== hash of the root node)
	Hash *Hash

	// Parent commit (only nil for initial commit)
	Parent *Commit

	store *Store
}

func NewEmptyCommit(store *Store, author id.ID) *Commit {
	return &Commit{
		store:   store,
		ModTime: time.Now(),
		Author:  author,
		Changes: make(map[*File]*Checkpoint),
	}
}

func (cm *Commit) FromProto(c *wire.Commit) error {
	author, err := id.Cast(c.GetAuthor())
	if err != nil {
		return err
	}

	modTime := time.Time{}
	if err := modTime.UnmarshalBinary(c.GetModTime()); err != nil {
		return err
	}

	hash, err := multihash.Cast(c.GetHash())
	if err != nil {
		return err
	}

	changes := make(map[*File]*Checkpoint)

	for _, change := range c.GetChanges() {
		file := cm.store.Root.Lookup(change.GetPath())
		if file == nil {
			// TODO: Which file? Make this a more specific error.
			return ErrNoSuchFile
		}

		checkpoint := &Checkpoint{}
		if err := checkpoint.FromProto(change.GetCheckpoint()); err != nil {
			return err
		}

		changes[file] = checkpoint
	}

	var parentCommit *Commit

	if c.GetParentHash() != nil {
		err = cm.store.viewWithBucket(
			"commits",
			func(tx *bolt.Tx, bckt *bolt.Bucket) error {
				parentData := bckt.Get(c.GetParentHash())
				if parentData == nil {
					return ErrNoSuchFile
				}

				protoCommit := &wire.Commit{}
				if err := proto.Unmarshal(parentData, protoCommit); err != nil {
					return err
				}

				return NewEmptyCommit(cm.store, "").FromProto(protoCommit)
			},
		)

		if err != nil {
			return err
		}
	}

	// Set commit data if everything worked:
	cm.Message = c.GetMessage()
	cm.Author = author
	cm.ModTime = modTime
	cm.Changes = changes
	cm.Hash = &Hash{hash}
	cm.Parent = parentCommit
	return nil
}

func (cm *Commit) ToProto() (*wire.Commit, error) {
	pcm := &wire.Commit{}
	modTime, err := cm.ModTime.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var changes []*wire.Change
	for file, checkpoint := range cm.Changes {
		protoCheckpoint, err := checkpoint.ToProto()
		if err != nil {
			return nil, err
		}

		changes = append(changes, &wire.Change{
			Path:       proto.String(file.Path()),
			Checkpoint: protoCheckpoint,
		})
	}

	pcm.Message = proto.String(cm.Message)
	pcm.Author = proto.String(string(cm.Author))
	pcm.ModTime = modTime
	pcm.Hash = cm.Hash.Bytes()
	pcm.Changes = changes

	if cm.Parent != nil {
		pcm.ParentHash = cm.Parent.Hash.Bytes()
	}

	return pcm, nil
}

func (cm *Commit) MarshalProto() ([]byte, error) {
	protoCmt, err := cm.ToProto()
	if err != nil {
		return nil, err
	}

	return proto.Marshal(protoCmt)
}

func (cm *Commit) UnmarshalProto(data []byte) error {
	protoCmt := &wire.Commit{}
	if err := proto.Unmarshal(data, protoCmt); err != nil {
		return err
	}

	return cm.FromProto(protoCmt)
}

// Checkpoint remembers one state of a single file.
type Checkpoint struct {
	// Hash is the hash of the file at this point.
	// It may, or may not be retrievable from ipfs.
	// For ChangeRemove the hash is the hash of the last existing file.
	Hash *Hash

	// ModTime is the time the checkpoint was made.
	ModTime time.Time

	// Size is the size of the file in bytes at this point.
	Size int64

	// Change is the detailed type of the modification.
	Change ChangeType

	// Author of the file modifications (jabber id)
	// TODO: Make separate Authorship struct.
	Author id.ID
}

// TODO: nice representation
func (c *Checkpoint) String() string {
	return fmt.Sprintf("%-7s %+7s@%s", c.Change.String(), c.Hash.B58String(), c.ModTime.String())
}

func (cp *Checkpoint) ToProto() (*wire.Checkpoint, error) {
	mtimeBin, err := cp.ModTime.MarshalBinary()
	if err != nil {
		return nil, err
	}

	protoCheck := &wire.Checkpoint{
		Hash:     cp.Hash.Bytes(),
		ModTime:  mtimeBin,
		FileSize: proto.Int64(cp.Size),
		Change:   proto.Int32(int32(cp.Change)),
		Author:   proto.String(string(cp.Author)),
	}

	if err != nil {
		return nil, err
	}

	return protoCheck, nil
}

func (cp *Checkpoint) FromProto(msg *wire.Checkpoint) error {
	modTime := time.Time{}
	if err := modTime.UnmarshalBinary(msg.GetModTime()); err != nil {
		return err
	}

	cp.Hash = &Hash{msg.GetHash()}
	cp.ModTime = modTime
	cp.Size = msg.GetFileSize()
	cp.Change = ChangeType(msg.GetChange())

	ID, err := id.Cast(msg.GetAuthor())
	if err != nil {
		cp.Author = ID
	}
	return nil
}

func (cp *Checkpoint) Marshal() ([]byte, error) {
	protoCheck, err := cp.ToProto()
	if err != nil {
		return nil, err
	}

	protoData, err := proto.Marshal(protoCheck)
	if err != nil {
		return nil, err
	}

	return protoData, nil
}

func (cp *Checkpoint) Unmarshal(data []byte) error {
	protoCheck := &wire.Checkpoint{}
	if err := proto.Unmarshal(data, protoCheck); err != nil {
		return err
	}

	return cp.FromProto(protoCheck)
}

// History remembers the changes made to a file.
// New changes get appended to the end.
type History []*Checkpoint

func (hy *History) ToProto() (*wire.History, error) {
	protoHist := &wire.History{}

	for _, ck := range *hy {
		protoCheck, err := ck.ToProto()
		if err != nil {
			return nil, err
		}

		protoHist.Hist = append(protoHist.Hist, protoCheck)
	}

	return protoHist, nil
}

func (hy *History) Marshal() ([]byte, error) {
	protoHist, err := hy.ToProto()
	if err != nil {
		return nil, err
	}

	data, err := proto.Marshal(protoHist)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (hy *History) FromProto(protoHist *wire.History) error {
	for _, protoCheck := range protoHist.Hist {
		ck := &Checkpoint{}
		if err := ck.FromProto(protoCheck); err != nil {
			return err
		}

		*hy = append(*hy, ck)
	}

	return nil
}

func (hy *History) Unmarshal(data []byte) error {
	protoHist := &wire.History{}

	if err := proto.Unmarshal(data, protoHist); err != nil {
		return err
	}

	return hy.FromProto(protoHist)
}

// MakeCheckpoint creates a new checkpoint in the version history of `curr`.
// One of old or curr may be nil (if no old version exists or new version
// does not exist anymore). It is an error to pass nil twice.
//
// If nothing changed between old and curr, ErrNoChange is returned.
func (st *Store) MakeCheckpoint(old, curr *Metadata, oldPath, currPath string) error {
	var change ChangeType
	var hash *Hash
	var path string
	var size int64

	if old == nil {
		change, path, hash, size = ChangeAdd, currPath, curr.hash, curr.size
	} else if curr == nil {
		change, path, hash, size = ChangeRemove, oldPath, old.hash, old.size
	} else if !curr.hash.Equal(old.hash) {
		change, path, hash, size = ChangeModify, currPath, curr.hash, curr.size
	} else if oldPath != currPath {
		change, path, hash, size = ChangeMove, currPath, curr.hash, curr.size
	} else {
		return ErrNoChange
	}

	checkpoint := &Checkpoint{
		Hash:    hash,
		ModTime: time.Now(),
		Size:    size,
		Change:  change,
		Author:  st.ID,
	}

	protoData, err := checkpoint.Marshal()
	if err != nil {
		return err
	}

	mtimeBin, err := checkpoint.ModTime.MarshalBinary()
	if err != nil {
		return err
	}

	dbErr := st.updateWithBucket("checkpoints", func(tx *bolt.Tx, bckt *bolt.Bucket) error {
		histBuck, err := bckt.CreateBucketIfNotExists([]byte(path))
		if err != nil {
			return err
		}

		// On a "move" we need to move the old data to the new path.
		if change == ChangeMove {
			if oldBuck := bckt.Bucket([]byte(oldPath)); oldBuck != nil {
				err = oldBuck.ForEach(func(k, v []byte) error {
					return histBuck.Put(k, v)
				})

				if err != nil {
					return err
				}

				if err := bckt.DeleteBucket([]byte(oldPath)); err != nil {
					return err
				}
			}
		}

		return histBuck.Put(mtimeBin, protoData)
	})

	if dbErr != nil {
		return dbErr
	}

	log.Debugf("created check point: %v", checkpoint)
	return nil
}

// History returns all checkpoints a file has.
// Note: even on error a empty history is returned.
func (s *Store) History(path string) (*History, error) {
	var hist History

	return &hist, s.viewWithBucket("checkpoints", func(tx *bolt.Tx, bckt *bolt.Bucket) error {
		changeBuck := bckt.Bucket([]byte(path))
		if changeBuck == nil {
			return ErrNoSuchFile
		}

		return changeBuck.ForEach(func(k, v []byte) error {
			ck := &Checkpoint{}
			if err := ck.Unmarshal(v); err != nil {
				return err
			}

			hist = append(hist, ck)
			return nil
		})
	})
}
