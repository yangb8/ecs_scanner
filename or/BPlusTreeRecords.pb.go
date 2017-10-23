// Code generated by protoc-gen-go. DO NOT EDIT.
// source: BPlusTreeRecords.proto

/*
Package or is a generated protocol buffer package.

It is generated from these files:
	BPlusTreeRecords.proto
	DirectoryRecords.proto
	JournalRecords.proto
	ObjectRecords.proto
	PageTypes.proto
	SchemaKeyRecords.proto
	SegmentLocationRecord.proto

It has these top-level messages:
	PagePair
	BPlusTreePage
	OccupancyPair
	OccupancyMap
	OccupancyEntry_V2
	Occupancy_V2
	BPlusTreeRecord
	ObjectSegmentInfo
	DirectoryUpdateRecord
	JournalHeader
	DirTableInstanceID
	JournalDirTableLogHeader
	VersionInfo
	SystemMetadata
	HeadSystemMetadata
	UserMetadata
	DataRange
	DataIndex
	IndexMetadataRecord
	SystemMetadataSet
	UserMetadataSet
	HeadMetadataSet
	MetadataSet
	CompactionMarker
	VersionHistory
	UpdateMetadataRecord
	CrossReferenceRecord
	PageItem
	SchemaKey
	DefaultSchemaKey
	DTRecordKey
	DTRecordJournalRegionSubKey
	DTRecordBPTreeInfoSubKey
	DTRecordBPTreeBootstrapSubKey
	DTRecordBPTreeBootstrapJournalSubKey
	DTOwnerKey
	TaskOrder
	MTTableRecordKey
	MTStorageStatKey
	MTAggregatedStorageStatKey
	MTBandwidthStatKey
	MTAggregatedBandwidthStatKey
	RMTaskKey
	ObjectTableKey
	RMTaskJournalEntryGeoSendKey
	RMTaskRecoveryPointSendKey
	SSTableRecordKey
	SSTableDeviceEntryKey
	SSTablePartitionEntryKey
	SSTableFreeBlockEntryKey
	SSTableBusyBlockEntryKey
	SSTableBlockBinEntryKey
	SSTableBlockLevelTaskKey
	DTBootstrapTaskKey
	GCRefCollectionKey
	ChunkKey
	CMTaskKey
	CMGeoInfoSendTaskKey
	CMGeoDataSendTaskKey
	CMGeoCopyTaskKey
	CMXorGroupTaskKey
	CMXorEncodeTaskKey
	CMXorDecodeTaskKey
	CMJobKey
	CMProgressKey
	ChunkSequenceKey
	RgReconfigTaskKey
	ListEntryKey
	DeleteJobTableKey
	BtreeReferenceKey
	RepoReferenceKey
	NamespaceKey
	BucketKey
	UserKey
	ConfigKey
	RGKey
	RGUpdateKey
	ResourceTableBootstrapTaskKey
	NKEntryKey
	NKEntryReplicationTaskKey
	ZKConfigKey
	ZoneInfoKey
	CompressInfo
	RangeInfo
	SegmentLocation
*/
package or

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// page type
type BPlusTreePage_PageType int32

const (
	BPlusTreePage_BLANK     BPlusTreePage_PageType = 0
	BPlusTreePage_ROOT      BPlusTreePage_PageType = 1
	BPlusTreePage_INDEX     BPlusTreePage_PageType = 2
	BPlusTreePage_LEAF      BPlusTreePage_PageType = 3
	BPlusTreePage_BLOOM     BPlusTreePage_PageType = 4
	BPlusTreePage_OCCUPANCY BPlusTreePage_PageType = 5
)

var BPlusTreePage_PageType_name = map[int32]string{
	0: "BLANK",
	1: "ROOT",
	2: "INDEX",
	3: "LEAF",
	4: "BLOOM",
	5: "OCCUPANCY",
}
var BPlusTreePage_PageType_value = map[string]int32{
	"BLANK":     0,
	"ROOT":      1,
	"INDEX":     2,
	"LEAF":      3,
	"BLOOM":     4,
	"OCCUPANCY": 5,
}

func (x BPlusTreePage_PageType) Enum() *BPlusTreePage_PageType {
	p := new(BPlusTreePage_PageType)
	*p = x
	return p
}
func (x BPlusTreePage_PageType) String() string {
	return proto.EnumName(BPlusTreePage_PageType_name, int32(x))
}
func (x *BPlusTreePage_PageType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(BPlusTreePage_PageType_value, data, "BPlusTreePage_PageType")
	if err != nil {
		return err
	}
	*x = BPlusTreePage_PageType(value)
	return nil
}
func (BPlusTreePage_PageType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

// a key -> value pair
//   when a leaf page then key points to a value (if we store value externally because
//       it is huge then we can use address)
//   when an index page then key points to an address
type PagePair struct {
	Key              *PageItem        `protobuf:"bytes,1,req,name=key" json:"key,omitempty"`
	Value            *PageItem        `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	Address          *SegmentLocation `protobuf:"bytes,3,opt,name=address" json:"address,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *PagePair) Reset()                    { *m = PagePair{} }
func (m *PagePair) String() string            { return proto.CompactTextString(m) }
func (*PagePair) ProtoMessage()               {}
func (*PagePair) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PagePair) GetKey() *PageItem {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *PagePair) GetValue() *PageItem {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *PagePair) GetAddress() *SegmentLocation {
	if m != nil {
		return m.Address
	}
	return nil
}

// currently tree pages are either
// 1. the actual tree pages, which are classified as root, index, leaf
// 2. bloom filter pages, containing a serialized bloom filter
// 3. occupancy pages, which contains a list of occupancy entries
type BPlusTreePage struct {
	Type *BPlusTreePage_PageType `protobuf:"varint,1,req,name=type,enum=or.BPlusTreePage_PageType,def=0" json:"type,omitempty"`
	// page entries -- only when page is ROOT, INDEX, or LEAF
	Elements []*PagePair `protobuf:"bytes,2,rep,name=elements" json:"elements,omitempty"`
	// currently unused
	Checkpoint *SegmentLocation `protobuf:"bytes,3,opt,name=checkpoint" json:"checkpoint,omitempty"`
	// bloom filter entry -- only when page is BLOOM
	BloomFilter *PageItem `protobuf:"bytes,4,opt,name=bloomFilter" json:"bloomFilter,omitempty"`
	// occupancy delta entry -- only when page is OCCUPANCY
	Occupancy        []*OccupancyEntry_V2 `protobuf:"bytes,5,rep,name=occupancy" json:"occupancy,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *BPlusTreePage) Reset()                    { *m = BPlusTreePage{} }
func (m *BPlusTreePage) String() string            { return proto.CompactTextString(m) }
func (*BPlusTreePage) ProtoMessage()               {}
func (*BPlusTreePage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

const Default_BPlusTreePage_Type BPlusTreePage_PageType = BPlusTreePage_BLANK

func (m *BPlusTreePage) GetType() BPlusTreePage_PageType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Default_BPlusTreePage_Type
}

func (m *BPlusTreePage) GetElements() []*PagePair {
	if m != nil {
		return m.Elements
	}
	return nil
}

func (m *BPlusTreePage) GetCheckpoint() *SegmentLocation {
	if m != nil {
		return m.Checkpoint
	}
	return nil
}

func (m *BPlusTreePage) GetBloomFilter() *PageItem {
	if m != nil {
		return m.BloomFilter
	}
	return nil
}

func (m *BPlusTreePage) GetOccupancy() []*OccupancyEntry_V2 {
	if m != nil {
		return m.Occupancy
	}
	return nil
}

// Release 1 and 1.1
type OccupancyPair struct {
	SegmentLocation  *SegmentLocation `protobuf:"bytes,1,req,name=segmentLocation" json:"segmentLocation,omitempty"`
	Occupied         *uint64          `protobuf:"varint,2,req,name=occupied,def=0" json:"occupied,omitempty"`
	TimeWhenZero     *uint64          `protobuf:"varint,3,opt,name=timeWhenZero" json:"timeWhenZero,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *OccupancyPair) Reset()                    { *m = OccupancyPair{} }
func (m *OccupancyPair) String() string            { return proto.CompactTextString(m) }
func (*OccupancyPair) ProtoMessage()               {}
func (*OccupancyPair) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

const Default_OccupancyPair_Occupied uint64 = 0

func (m *OccupancyPair) GetSegmentLocation() *SegmentLocation {
	if m != nil {
		return m.SegmentLocation
	}
	return nil
}

func (m *OccupancyPair) GetOccupied() uint64 {
	if m != nil && m.Occupied != nil {
		return *m.Occupied
	}
	return Default_OccupancyPair_Occupied
}

func (m *OccupancyPair) GetTimeWhenZero() uint64 {
	if m != nil && m.TimeWhenZero != nil {
		return *m.TimeWhenZero
	}
	return 0
}

type OccupancyMap struct {
	Map              []*OccupancyPair `protobuf:"bytes,1,rep,name=map" json:"map,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *OccupancyMap) Reset()                    { *m = OccupancyMap{} }
func (m *OccupancyMap) String() string            { return proto.CompactTextString(m) }
func (*OccupancyMap) ProtoMessage()               {}
func (*OccupancyMap) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *OccupancyMap) GetMap() []*OccupancyPair {
	if m != nil {
		return m.Map
	}
	return nil
}

// Release 2
type OccupancyEntry_V2 struct {
	ChunkId          *string          `protobuf:"bytes,1,opt,name=chunkId" json:"chunkId,omitempty"`
	Location         *SegmentLocation `protobuf:"bytes,2,opt,name=location" json:"location,omitempty"`
	Occupancy        *int64           `protobuf:"zigzag64,3,opt,name=occupancy" json:"occupancy,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *OccupancyEntry_V2) Reset()                    { *m = OccupancyEntry_V2{} }
func (m *OccupancyEntry_V2) String() string            { return proto.CompactTextString(m) }
func (*OccupancyEntry_V2) ProtoMessage()               {}
func (*OccupancyEntry_V2) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *OccupancyEntry_V2) GetChunkId() string {
	if m != nil && m.ChunkId != nil {
		return *m.ChunkId
	}
	return ""
}

func (m *OccupancyEntry_V2) GetLocation() *SegmentLocation {
	if m != nil {
		return m.Location
	}
	return nil
}

func (m *OccupancyEntry_V2) GetOccupancy() int64 {
	if m != nil && m.Occupancy != nil {
		return *m.Occupancy
	}
	return 0
}

// the listing of OccupancyEntries could be larger than the physical limit of a page
// thus the listing may need to span multiple OccupancyPages
type Occupancy_V2 struct {
	OccupancyPages   []*SegmentLocation `protobuf:"bytes,1,rep,name=occupancyPages" json:"occupancyPages,omitempty"`
	Timestamp        *uint64            `protobuf:"varint,2,opt,name=timestamp" json:"timestamp,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (m *Occupancy_V2) Reset()                    { *m = Occupancy_V2{} }
func (m *Occupancy_V2) String() string            { return proto.CompactTextString(m) }
func (*Occupancy_V2) ProtoMessage()               {}
func (*Occupancy_V2) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Occupancy_V2) GetOccupancyPages() []*SegmentLocation {
	if m != nil {
		return m.OccupancyPages
	}
	return nil
}

func (m *Occupancy_V2) GetTimestamp() uint64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

// BPlusTreeRecord is a merege of v1 and v2
type BPlusTreeRecord struct {
	IndexRoot    *SegmentLocation   `protobuf:"bytes,1,opt,name=indexRoot" json:"indexRoot,omitempty"`
	OccupancyMap *SegmentLocation   `protobuf:"bytes,2,opt,name=occupancyMap" json:"occupancyMap,omitempty"`
	BloomFilter  []*SegmentLocation `protobuf:"bytes,3,rep,name=bloomFilter" json:"bloomFilter,omitempty"`
	// Release 2 specific
	OccupancyDeltaBTreeChunks *Occupancy_V2 `protobuf:"bytes,4,opt,name=occupancyDeltaBTreeChunks" json:"occupancyDeltaBTreeChunks,omitempty"`
	TreeVersion               *uint64       `protobuf:"varint,5,opt,name=treeVersion" json:"treeVersion,omitempty"`
	OccupancyDeltaRepoChunks  *Occupancy_V2 `protobuf:"bytes,6,opt,name=occupancyDeltaRepoChunks" json:"occupancyDeltaRepoChunks,omitempty"`
	OccupancyFullBTreeChunks  *Occupancy_V2 `protobuf:"bytes,7,opt,name=occupancyFullBTreeChunks" json:"occupancyFullBTreeChunks,omitempty"`
	OccupancyFullRepoChunks   *Occupancy_V2 `protobuf:"bytes,8,opt,name=occupancyFullRepoChunks" json:"occupancyFullRepoChunks,omitempty"`
	XXX_unrecognized          []byte        `json:"-"`
}

func (m *BPlusTreeRecord) Reset()                    { *m = BPlusTreeRecord{} }
func (m *BPlusTreeRecord) String() string            { return proto.CompactTextString(m) }
func (*BPlusTreeRecord) ProtoMessage()               {}
func (*BPlusTreeRecord) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *BPlusTreeRecord) GetIndexRoot() *SegmentLocation {
	if m != nil {
		return m.IndexRoot
	}
	return nil
}

func (m *BPlusTreeRecord) GetOccupancyMap() *SegmentLocation {
	if m != nil {
		return m.OccupancyMap
	}
	return nil
}

func (m *BPlusTreeRecord) GetBloomFilter() []*SegmentLocation {
	if m != nil {
		return m.BloomFilter
	}
	return nil
}

func (m *BPlusTreeRecord) GetOccupancyDeltaBTreeChunks() *Occupancy_V2 {
	if m != nil {
		return m.OccupancyDeltaBTreeChunks
	}
	return nil
}

func (m *BPlusTreeRecord) GetTreeVersion() uint64 {
	if m != nil && m.TreeVersion != nil {
		return *m.TreeVersion
	}
	return 0
}

func (m *BPlusTreeRecord) GetOccupancyDeltaRepoChunks() *Occupancy_V2 {
	if m != nil {
		return m.OccupancyDeltaRepoChunks
	}
	return nil
}

func (m *BPlusTreeRecord) GetOccupancyFullBTreeChunks() *Occupancy_V2 {
	if m != nil {
		return m.OccupancyFullBTreeChunks
	}
	return nil
}

func (m *BPlusTreeRecord) GetOccupancyFullRepoChunks() *Occupancy_V2 {
	if m != nil {
		return m.OccupancyFullRepoChunks
	}
	return nil
}

func init() {
	proto.RegisterType((*PagePair)(nil), "or.PagePair")
	proto.RegisterType((*BPlusTreePage)(nil), "or.BPlusTreePage")
	proto.RegisterType((*OccupancyPair)(nil), "or.OccupancyPair")
	proto.RegisterType((*OccupancyMap)(nil), "or.OccupancyMap")
	proto.RegisterType((*OccupancyEntry_V2)(nil), "or.OccupancyEntry_V2")
	proto.RegisterType((*Occupancy_V2)(nil), "or.Occupancy_V2")
	proto.RegisterType((*BPlusTreeRecord)(nil), "or.BPlusTreeRecord")
	proto.RegisterEnum("or.BPlusTreePage_PageType", BPlusTreePage_PageType_name, BPlusTreePage_PageType_value)
}

func init() { proto.RegisterFile("BPlusTreeRecords.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 633 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0x8e, 0x43, 0x93, 0x69, 0xda, 0xba, 0x83, 0x00, 0x53, 0x7e, 0x14, 0x99, 0x4b, 0x2e,
	0x04, 0x68, 0x55, 0x21, 0x15, 0x71, 0x68, 0xd2, 0x56, 0x2a, 0xa4, 0x49, 0xb4, 0x94, 0xf2, 0x73,
	0x41, 0xc6, 0x19, 0xb5, 0x56, 0x1d, 0xaf, 0xb5, 0xde, 0x20, 0x22, 0xc1, 0x2b, 0x20, 0xf1, 0x7a,
	0x3c, 0x0d, 0xda, 0x4d, 0xe2, 0xd8, 0x6e, 0x5d, 0x8e, 0x99, 0xf9, 0x66, 0xe6, 0xfb, 0x59, 0x07,
	0xee, 0x75, 0x86, 0xe1, 0x24, 0x39, 0x15, 0x44, 0x8c, 0x7c, 0x2e, 0x46, 0x49, 0x3b, 0x16, 0x5c,
	0x72, 0x34, 0xb9, 0xd8, 0x7a, 0xf8, 0x9e, 0xce, 0xc7, 0x14, 0xc9, 0x1e, 0xf7, 0x3d, 0x19, 0xf0,
	0x68, 0x86, 0x98, 0x01, 0xb6, 0x36, 0x86, 0xde, 0x39, 0x9d, 0x4e, 0x63, 0x9a, 0x4f, 0xb8, 0xbf,
	0xa0, 0xa6, 0x4a, 0x43, 0x2f, 0x10, 0xf8, 0x04, 0x2a, 0x97, 0x34, 0x75, 0x8c, 0xa6, 0xd9, 0x5a,
	0xdd, 0x6e, 0xb4, 0xb9, 0x68, 0xab, 0xd6, 0xb1, 0xa4, 0x31, 0x53, 0x0d, 0x74, 0xa1, 0xfa, 0xdd,
	0x0b, 0x27, 0xe4, 0x98, 0x4d, 0xe3, 0x0a, 0x62, 0xd6, 0xc2, 0x67, 0xb0, 0xe2, 0x8d, 0x46, 0x82,
	0x92, 0xc4, 0xa9, 0x68, 0xd4, 0x1d, 0x85, 0x2a, 0x52, 0x5a, 0x60, 0xdc, 0xbf, 0x26, 0xac, 0xa5,
	0x5a, 0xd4, 0x2e, 0xdc, 0x05, 0x4b, 0x4e, 0x63, 0xd2, 0x2c, 0xd6, 0xb7, 0xb7, 0xd4, 0x74, 0x0e,
	0xd0, 0x5e, 0x28, 0xd8, 0xab, 0x76, 0x7a, 0xfb, 0xfd, 0x77, 0x4c, 0xc3, 0xb1, 0x05, 0x35, 0x0a,
	0x49, 0x1d, 0x49, 0x1c, 0xb3, 0x59, 0xc9, 0xd2, 0x53, 0xda, 0x58, 0xda, 0xc5, 0x1d, 0x00, 0xff,
	0x82, 0xfc, 0xcb, 0x98, 0x07, 0x91, 0xbc, 0x89, 0x64, 0x06, 0x86, 0x6d, 0x58, 0xfd, 0x16, 0x72,
	0x3e, 0x3e, 0x0a, 0x42, 0x49, 0xc2, 0xb1, 0xae, 0x31, 0x20, 0x0b, 0xc0, 0x1d, 0xa8, 0x73, 0xdf,
	0x9f, 0xc4, 0x5e, 0xe4, 0x4f, 0x9d, 0xaa, 0xe6, 0x73, 0x57, 0xa1, 0x07, 0x8b, 0xe2, 0x61, 0x24,
	0xc5, 0xf4, 0xeb, 0xd9, 0x36, 0x5b, 0xe2, 0xdc, 0xfe, 0x2c, 0x0b, 0x25, 0x0e, 0xeb, 0x30, 0x93,
	0x67, 0xdf, 0xc2, 0x1a, 0x58, 0x6c, 0x30, 0x38, 0xb5, 0x0d, 0x55, 0x3c, 0xee, 0x1f, 0x1c, 0x7e,
	0xb2, 0x4d, 0x55, 0xec, 0x1d, 0xee, 0x1f, 0xd9, 0x95, 0x19, 0x72, 0x30, 0x38, 0xb1, 0x2d, 0x5c,
	0x83, 0xfa, 0xa0, 0xdb, 0xfd, 0x30, 0xdc, 0xef, 0x77, 0x3f, 0xdb, 0x55, 0xf7, 0x8f, 0x01, 0x6b,
	0xe9, 0x41, 0x9d, 0xf0, 0x1b, 0xd8, 0x48, 0xf2, 0x2a, 0xe7, 0x69, 0x5f, 0x6b, 0x40, 0x11, 0x8b,
	0x8f, 0xa1, 0xa6, 0xd9, 0x06, 0x34, 0x72, 0xcc, 0xa6, 0xd9, 0xb2, 0xf6, 0x8c, 0x17, 0x2c, 0x2d,
	0xa1, 0x0b, 0x0d, 0x19, 0x8c, 0xe9, 0xe3, 0x05, 0x45, 0x5f, 0x48, 0x70, 0xed, 0xad, 0xc5, 0x72,
	0x35, 0x77, 0x07, 0x1a, 0x29, 0xa5, 0x13, 0x2f, 0xc6, 0xa7, 0x50, 0x19, 0x7b, 0xb1, 0x63, 0x68,
	0x8b, 0x36, 0x73, 0x16, 0xe9, 0xdc, 0x54, 0xd7, 0xfd, 0x09, 0x9b, 0x57, 0x8c, 0x43, 0x07, 0x56,
	0xfc, 0x8b, 0x49, 0x74, 0x79, 0x3c, 0x72, 0x8c, 0xa6, 0xd1, 0xaa, 0xb3, 0xc5, 0x4f, 0x7c, 0x0e,
	0xb5, 0x70, 0x21, 0xcf, 0x2c, 0xcf, 0x37, 0x05, 0xe1, 0xa3, 0x6c, 0x5a, 0x8a, 0x35, 0x66, 0x63,
	0x09, 0x32, 0x94, 0xd5, 0xe1, 0xd7, 0xb0, 0xce, 0x97, 0x1c, 0xcf, 0x29, 0x99, 0xb3, 0xbf, 0xf6,
	0x48, 0x01, 0xaa, 0x4e, 0x29, 0x3f, 0x12, 0xe9, 0x8d, 0x63, 0x4d, 0xce, 0x62, 0xcb, 0x82, 0xfb,
	0xdb, 0x82, 0x8d, 0xc2, 0xa7, 0x8d, 0x2f, 0xa1, 0x1e, 0x44, 0x23, 0xfa, 0xc1, 0x38, 0x97, 0x5a,
	0x69, 0xc9, 0xa5, 0x25, 0x0a, 0x5f, 0x41, 0x83, 0x67, 0x4c, 0xbe, 0xc9, 0x84, 0x1c, 0x10, 0x77,
	0xf3, 0xcf, 0xbc, 0x52, 0xae, 0x2b, 0xf7, 0xda, 0xfb, 0xf0, 0x20, 0x5d, 0x73, 0x40, 0xa1, 0xf4,
	0x3a, 0x8a, 0x7f, 0x57, 0xc5, 0x91, 0xcc, 0xbf, 0x15, 0x3b, 0x17, 0xad, 0x7a, 0xf8, 0xe5, 0x23,
	0xd8, 0x84, 0x55, 0x29, 0x88, 0xce, 0x48, 0x24, 0x2a, 0xc3, 0xaa, 0xb6, 0x29, 0x5b, 0xc2, 0x1e,
	0x38, 0xf9, 0x71, 0x46, 0x31, 0x9f, 0x1f, 0xbc, 0x5d, 0x72, 0xb0, 0x74, 0x22, 0xb7, 0xed, 0x68,
	0x12, 0x86, 0x59, 0xfa, 0x2b, 0xff, 0xdd, 0x56, 0x98, 0xc0, 0xb7, 0x70, 0x3f, 0xd7, 0xcb, 0x50,
	0xab, 0x95, 0x2c, 0x2b, 0x1b, 0xe8, 0x60, 0xc7, 0x2e, 0xfe, 0xd5, 0xff, 0x0b, 0x00, 0x00, 0xff,
	0xff, 0xc6, 0x46, 0xb4, 0xd6, 0xfd, 0x05, 0x00, 0x00,
}
