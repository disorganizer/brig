// Code generated by capnpc-go. DO NOT EDIT.

package capnp

import (
	capnp "zombiezen.com/go/capnproto2"
	text "zombiezen.com/go/capnproto2/encoding/text"
	schemas "zombiezen.com/go/capnproto2/schemas"
)

type User struct{ capnp.Struct }

// User_TypeID is the unique identifier for the type User.
const User_TypeID = 0x861de4463c5a4a22

func NewUser(s *capnp.Segment) (User, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 4})
	return User{st}, err
}

func NewRootUser(s *capnp.Segment) (User, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 4})
	return User{st}, err
}

func ReadRootUser(msg *capnp.Message) (User, error) {
	root, err := msg.RootPtr()
	return User{root.Struct()}, err
}

func (s User) String() string {
	str, _ := text.Marshal(0x861de4463c5a4a22, s.Struct)
	return str
}

func (s User) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s User) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s User) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s User) SetName(v string) error {
	return s.Struct.SetText(0, v)
}

func (s User) PasswordHash() (string, error) {
	p, err := s.Struct.Ptr(1)
	return p.Text(), err
}

func (s User) HasPasswordHash() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s User) PasswordHashBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(1)
	return p.TextBytes(), err
}

func (s User) SetPasswordHash(v string) error {
	return s.Struct.SetText(1, v)
}

func (s User) Salt() (string, error) {
	p, err := s.Struct.Ptr(2)
	return p.Text(), err
}

func (s User) HasSalt() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s User) SaltBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(2)
	return p.TextBytes(), err
}

func (s User) SetSalt(v string) error {
	return s.Struct.SetText(2, v)
}

func (s User) Folders() (capnp.TextList, error) {
	p, err := s.Struct.Ptr(3)
	return capnp.TextList{List: p.List()}, err
}

func (s User) HasFolders() bool {
	p, err := s.Struct.Ptr(3)
	return p.IsValid() || err != nil
}

func (s User) SetFolders(v capnp.TextList) error {
	return s.Struct.SetPtr(3, v.List.ToPtr())
}

// NewFolders sets the folders field to a newly
// allocated capnp.TextList, preferring placement in s's segment.
func (s User) NewFolders(n int32) (capnp.TextList, error) {
	l, err := capnp.NewTextList(s.Struct.Segment(), n)
	if err != nil {
		return capnp.TextList{}, err
	}
	err = s.Struct.SetPtr(3, l.List.ToPtr())
	return l, err
}

// User_List is a list of User.
type User_List struct{ capnp.List }

// NewUser creates a new list of User.
func NewUser_List(s *capnp.Segment, sz int32) (User_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 4}, sz)
	return User_List{l}, err
}

func (s User_List) At(i int) User { return User{s.List.Struct(i)} }

func (s User_List) Set(i int, v User) error { return s.List.SetStruct(i, v.Struct) }

func (s User_List) String() string {
	str, _ := text.MarshalList(0x861de4463c5a4a22, s.List)
	return str
}

// User_Promise is a wrapper for a User promised by a client call.
type User_Promise struct{ *capnp.Pipeline }

func (p User_Promise) Struct() (User, error) {
	s, err := p.Pipeline.Struct()
	return User{s}, err
}

const schema_a0b1c18bd0f965c4 = "x\xda4\xca\xb1J\xc3P\x18\xc5\xf1s\xee\x97X\x90" +
	"\xd2\xf4BF]\x1c\x05\x0d]\xc5A\x1cD\x9c\xfa\x0d" +
	".\xe2r5W\x83\xa4m\xcc\x8d\x14'\x07A\x04_" +
	"\xc2W\xf0\x09D\xd0]\x07\xdf@\xf0\x19\x9c\"-d" +
	";\xe7\xc7\x7fx\xbfgF\xf1+\x01M\xe3\x95v\xe3" +
	"\xe8d\xf7\xe0g\xfd\x01v\x8d\xed\xbb\xff\xfb|z{" +
	"yF\x1c\xf5\x80\xd1\xc7*\xedw\x0f\xb0_\xbf\xd8j" +
	"/]\xe3\xe7\xee6\x93\xfc,;w\xd5\xb4\xcan\x82" +
	"\xaf\xb7\x97s\xe78\xf8\x1a\x18\x93:\x94\x08\x88\x08X" +
	"\xb7\x09\xe8\xa9P\x0bCK\xa6\\\xa0\xbf\x024\x17j" +
	"eh\x8dIi\x00;Y\x94\x85P\x1bC+\x92R" +
	"\x00{\xbd\x0fh)\xd4G\xc3d\xea&\x9e}\x18\xf6" +
	"\xc1\xb6r!\xccgu\x8e\xe4\xd0\x85\xa2\xe3$\xb8\xb2" +
	"\xe9\xce\xdd\xc5\xac\xcc}\x1d8\x00\xc7\xc2%\x0f\xc0\xff" +
	"\x00\x00\x00\xff\xff\xd0\xc75\xe2"

func init() {
	schemas.Register(schema_a0b1c18bd0f965c4,
		0x861de4463c5a4a22)
}
