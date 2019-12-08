// Code generated by protoc-gen-go. DO NOT EDIT.
// source: meta.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Meta struct {
	VotedFor             int32    `protobuf:"varint,1,opt,name=VotedFor,json=votedFor,proto3" json:"VotedFor,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Meta) Reset()         { *m = Meta{} }
func (m *Meta) String() string { return proto.CompactTextString(m) }
func (*Meta) ProtoMessage()    {}
func (*Meta) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b5ea8fe65782bcc, []int{0}
}

func (m *Meta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Meta.Unmarshal(m, b)
}
func (m *Meta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Meta.Marshal(b, m, deterministic)
}
func (m *Meta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Meta.Merge(m, src)
}
func (m *Meta) XXX_Size() int {
	return xxx_messageInfo_Meta.Size(m)
}
func (m *Meta) XXX_DiscardUnknown() {
	xxx_messageInfo_Meta.DiscardUnknown(m)
}

var xxx_messageInfo_Meta proto.InternalMessageInfo

func (m *Meta) GetVotedFor() int32 {
	if m != nil {
		return m.VotedFor
	}
	return 0
}

func init() {
	proto.RegisterType((*Meta)(nil), "pb.Meta")
}

func init() { proto.RegisterFile("meta.proto", fileDescriptor_3b5ea8fe65782bcc) }

var fileDescriptor_3b5ea8fe65782bcc = []byte{
	// 76 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xca, 0x4d, 0x2d, 0x49,
	0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0x52, 0xe2, 0x62, 0xf1, 0x4d,
	0x2d, 0x49, 0x14, 0x92, 0xe2, 0xe2, 0x08, 0xcb, 0x2f, 0x49, 0x4d, 0x71, 0xcb, 0x2f, 0x92, 0x60,
	0x54, 0x60, 0xd4, 0x60, 0x0d, 0xe2, 0x28, 0x83, 0xf2, 0x93, 0xd8, 0xc0, 0xca, 0x8d, 0x01, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x62, 0x9f, 0xae, 0xd1, 0x3c, 0x00, 0x00, 0x00,
}
