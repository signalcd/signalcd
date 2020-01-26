// Code generated by protoc-gen-go. DO NOT EDIT.
// source: signalcd/proto/ui.proto

package signalcdproto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type ListDeploymentRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListDeploymentRequest) Reset()         { *m = ListDeploymentRequest{} }
func (m *ListDeploymentRequest) String() string { return proto.CompactTextString(m) }
func (*ListDeploymentRequest) ProtoMessage()    {}
func (*ListDeploymentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_df79a9f155088a03, []int{0}
}

func (m *ListDeploymentRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListDeploymentRequest.Unmarshal(m, b)
}
func (m *ListDeploymentRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListDeploymentRequest.Marshal(b, m, deterministic)
}
func (m *ListDeploymentRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListDeploymentRequest.Merge(m, src)
}
func (m *ListDeploymentRequest) XXX_Size() int {
	return xxx_messageInfo_ListDeploymentRequest.Size(m)
}
func (m *ListDeploymentRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListDeploymentRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListDeploymentRequest proto.InternalMessageInfo

type ListDeploymentResponse struct {
	Deployments          []*Deployment `protobuf:"bytes,1,rep,name=deployments,proto3" json:"deployments,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ListDeploymentResponse) Reset()         { *m = ListDeploymentResponse{} }
func (m *ListDeploymentResponse) String() string { return proto.CompactTextString(m) }
func (*ListDeploymentResponse) ProtoMessage()    {}
func (*ListDeploymentResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_df79a9f155088a03, []int{1}
}

func (m *ListDeploymentResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListDeploymentResponse.Unmarshal(m, b)
}
func (m *ListDeploymentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListDeploymentResponse.Marshal(b, m, deterministic)
}
func (m *ListDeploymentResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListDeploymentResponse.Merge(m, src)
}
func (m *ListDeploymentResponse) XXX_Size() int {
	return xxx_messageInfo_ListDeploymentResponse.Size(m)
}
func (m *ListDeploymentResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListDeploymentResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListDeploymentResponse proto.InternalMessageInfo

func (m *ListDeploymentResponse) GetDeployments() []*Deployment {
	if m != nil {
		return m.Deployments
	}
	return nil
}

type GetCurrentDeploymentRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetCurrentDeploymentRequest) Reset()         { *m = GetCurrentDeploymentRequest{} }
func (m *GetCurrentDeploymentRequest) String() string { return proto.CompactTextString(m) }
func (*GetCurrentDeploymentRequest) ProtoMessage()    {}
func (*GetCurrentDeploymentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_df79a9f155088a03, []int{2}
}

func (m *GetCurrentDeploymentRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetCurrentDeploymentRequest.Unmarshal(m, b)
}
func (m *GetCurrentDeploymentRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetCurrentDeploymentRequest.Marshal(b, m, deterministic)
}
func (m *GetCurrentDeploymentRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetCurrentDeploymentRequest.Merge(m, src)
}
func (m *GetCurrentDeploymentRequest) XXX_Size() int {
	return xxx_messageInfo_GetCurrentDeploymentRequest.Size(m)
}
func (m *GetCurrentDeploymentRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetCurrentDeploymentRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetCurrentDeploymentRequest proto.InternalMessageInfo

type GetCurrentDeploymentResponse struct {
	Deployment           *Deployment `protobuf:"bytes,1,opt,name=deployment,proto3" json:"deployment,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *GetCurrentDeploymentResponse) Reset()         { *m = GetCurrentDeploymentResponse{} }
func (m *GetCurrentDeploymentResponse) String() string { return proto.CompactTextString(m) }
func (*GetCurrentDeploymentResponse) ProtoMessage()    {}
func (*GetCurrentDeploymentResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_df79a9f155088a03, []int{3}
}

func (m *GetCurrentDeploymentResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetCurrentDeploymentResponse.Unmarshal(m, b)
}
func (m *GetCurrentDeploymentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetCurrentDeploymentResponse.Marshal(b, m, deterministic)
}
func (m *GetCurrentDeploymentResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetCurrentDeploymentResponse.Merge(m, src)
}
func (m *GetCurrentDeploymentResponse) XXX_Size() int {
	return xxx_messageInfo_GetCurrentDeploymentResponse.Size(m)
}
func (m *GetCurrentDeploymentResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetCurrentDeploymentResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetCurrentDeploymentResponse proto.InternalMessageInfo

func (m *GetCurrentDeploymentResponse) GetDeployment() *Deployment {
	if m != nil {
		return m.Deployment
	}
	return nil
}

type SetCurrentDeploymentRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetCurrentDeploymentRequest) Reset()         { *m = SetCurrentDeploymentRequest{} }
func (m *SetCurrentDeploymentRequest) String() string { return proto.CompactTextString(m) }
func (*SetCurrentDeploymentRequest) ProtoMessage()    {}
func (*SetCurrentDeploymentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_df79a9f155088a03, []int{4}
}

func (m *SetCurrentDeploymentRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetCurrentDeploymentRequest.Unmarshal(m, b)
}
func (m *SetCurrentDeploymentRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetCurrentDeploymentRequest.Marshal(b, m, deterministic)
}
func (m *SetCurrentDeploymentRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetCurrentDeploymentRequest.Merge(m, src)
}
func (m *SetCurrentDeploymentRequest) XXX_Size() int {
	return xxx_messageInfo_SetCurrentDeploymentRequest.Size(m)
}
func (m *SetCurrentDeploymentRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetCurrentDeploymentRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetCurrentDeploymentRequest proto.InternalMessageInfo

func (m *SetCurrentDeploymentRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type SetCurrentDeploymentResponse struct {
	Deployment           *Deployment `protobuf:"bytes,1,opt,name=deployment,proto3" json:"deployment,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *SetCurrentDeploymentResponse) Reset()         { *m = SetCurrentDeploymentResponse{} }
func (m *SetCurrentDeploymentResponse) String() string { return proto.CompactTextString(m) }
func (*SetCurrentDeploymentResponse) ProtoMessage()    {}
func (*SetCurrentDeploymentResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_df79a9f155088a03, []int{5}
}

func (m *SetCurrentDeploymentResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetCurrentDeploymentResponse.Unmarshal(m, b)
}
func (m *SetCurrentDeploymentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetCurrentDeploymentResponse.Marshal(b, m, deterministic)
}
func (m *SetCurrentDeploymentResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetCurrentDeploymentResponse.Merge(m, src)
}
func (m *SetCurrentDeploymentResponse) XXX_Size() int {
	return xxx_messageInfo_SetCurrentDeploymentResponse.Size(m)
}
func (m *SetCurrentDeploymentResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SetCurrentDeploymentResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SetCurrentDeploymentResponse proto.InternalMessageInfo

func (m *SetCurrentDeploymentResponse) GetDeployment() *Deployment {
	if m != nil {
		return m.Deployment
	}
	return nil
}

type ListPipelinesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListPipelinesRequest) Reset()         { *m = ListPipelinesRequest{} }
func (m *ListPipelinesRequest) String() string { return proto.CompactTextString(m) }
func (*ListPipelinesRequest) ProtoMessage()    {}
func (*ListPipelinesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_df79a9f155088a03, []int{6}
}

func (m *ListPipelinesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListPipelinesRequest.Unmarshal(m, b)
}
func (m *ListPipelinesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListPipelinesRequest.Marshal(b, m, deterministic)
}
func (m *ListPipelinesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListPipelinesRequest.Merge(m, src)
}
func (m *ListPipelinesRequest) XXX_Size() int {
	return xxx_messageInfo_ListPipelinesRequest.Size(m)
}
func (m *ListPipelinesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListPipelinesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListPipelinesRequest proto.InternalMessageInfo

type ListPipelinesResponse struct {
	Pipelines            []*Pipeline `protobuf:"bytes,1,rep,name=pipelines,proto3" json:"pipelines,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ListPipelinesResponse) Reset()         { *m = ListPipelinesResponse{} }
func (m *ListPipelinesResponse) String() string { return proto.CompactTextString(m) }
func (*ListPipelinesResponse) ProtoMessage()    {}
func (*ListPipelinesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_df79a9f155088a03, []int{7}
}

func (m *ListPipelinesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListPipelinesResponse.Unmarshal(m, b)
}
func (m *ListPipelinesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListPipelinesResponse.Marshal(b, m, deterministic)
}
func (m *ListPipelinesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListPipelinesResponse.Merge(m, src)
}
func (m *ListPipelinesResponse) XXX_Size() int {
	return xxx_messageInfo_ListPipelinesResponse.Size(m)
}
func (m *ListPipelinesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListPipelinesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListPipelinesResponse proto.InternalMessageInfo

func (m *ListPipelinesResponse) GetPipelines() []*Pipeline {
	if m != nil {
		return m.Pipelines
	}
	return nil
}

type CreatePipelineRequest struct {
	Pipeline             *Pipeline `protobuf:"bytes,1,opt,name=pipeline,proto3" json:"pipeline,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *CreatePipelineRequest) Reset()         { *m = CreatePipelineRequest{} }
func (m *CreatePipelineRequest) String() string { return proto.CompactTextString(m) }
func (*CreatePipelineRequest) ProtoMessage()    {}
func (*CreatePipelineRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_df79a9f155088a03, []int{8}
}

func (m *CreatePipelineRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreatePipelineRequest.Unmarshal(m, b)
}
func (m *CreatePipelineRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreatePipelineRequest.Marshal(b, m, deterministic)
}
func (m *CreatePipelineRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreatePipelineRequest.Merge(m, src)
}
func (m *CreatePipelineRequest) XXX_Size() int {
	return xxx_messageInfo_CreatePipelineRequest.Size(m)
}
func (m *CreatePipelineRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreatePipelineRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreatePipelineRequest proto.InternalMessageInfo

func (m *CreatePipelineRequest) GetPipeline() *Pipeline {
	if m != nil {
		return m.Pipeline
	}
	return nil
}

type CreatePipelineResponse struct {
	Pipeline             *Pipeline `protobuf:"bytes,1,opt,name=pipeline,proto3" json:"pipeline,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *CreatePipelineResponse) Reset()         { *m = CreatePipelineResponse{} }
func (m *CreatePipelineResponse) String() string { return proto.CompactTextString(m) }
func (*CreatePipelineResponse) ProtoMessage()    {}
func (*CreatePipelineResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_df79a9f155088a03, []int{9}
}

func (m *CreatePipelineResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreatePipelineResponse.Unmarshal(m, b)
}
func (m *CreatePipelineResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreatePipelineResponse.Marshal(b, m, deterministic)
}
func (m *CreatePipelineResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreatePipelineResponse.Merge(m, src)
}
func (m *CreatePipelineResponse) XXX_Size() int {
	return xxx_messageInfo_CreatePipelineResponse.Size(m)
}
func (m *CreatePipelineResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreatePipelineResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreatePipelineResponse proto.InternalMessageInfo

func (m *CreatePipelineResponse) GetPipeline() *Pipeline {
	if m != nil {
		return m.Pipeline
	}
	return nil
}

type GetPipelineRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPipelineRequest) Reset()         { *m = GetPipelineRequest{} }
func (m *GetPipelineRequest) String() string { return proto.CompactTextString(m) }
func (*GetPipelineRequest) ProtoMessage()    {}
func (*GetPipelineRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_df79a9f155088a03, []int{10}
}

func (m *GetPipelineRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPipelineRequest.Unmarshal(m, b)
}
func (m *GetPipelineRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPipelineRequest.Marshal(b, m, deterministic)
}
func (m *GetPipelineRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPipelineRequest.Merge(m, src)
}
func (m *GetPipelineRequest) XXX_Size() int {
	return xxx_messageInfo_GetPipelineRequest.Size(m)
}
func (m *GetPipelineRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPipelineRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetPipelineRequest proto.InternalMessageInfo

func (m *GetPipelineRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type GetPipelineResponse struct {
	Pipeline             *Pipeline `protobuf:"bytes,1,opt,name=pipeline,proto3" json:"pipeline,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *GetPipelineResponse) Reset()         { *m = GetPipelineResponse{} }
func (m *GetPipelineResponse) String() string { return proto.CompactTextString(m) }
func (*GetPipelineResponse) ProtoMessage()    {}
func (*GetPipelineResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_df79a9f155088a03, []int{11}
}

func (m *GetPipelineResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPipelineResponse.Unmarshal(m, b)
}
func (m *GetPipelineResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPipelineResponse.Marshal(b, m, deterministic)
}
func (m *GetPipelineResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPipelineResponse.Merge(m, src)
}
func (m *GetPipelineResponse) XXX_Size() int {
	return xxx_messageInfo_GetPipelineResponse.Size(m)
}
func (m *GetPipelineResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPipelineResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetPipelineResponse proto.InternalMessageInfo

func (m *GetPipelineResponse) GetPipeline() *Pipeline {
	if m != nil {
		return m.Pipeline
	}
	return nil
}

func init() {
	proto.RegisterType((*ListDeploymentRequest)(nil), "signalcd.ListDeploymentRequest")
	proto.RegisterType((*ListDeploymentResponse)(nil), "signalcd.ListDeploymentResponse")
	proto.RegisterType((*GetCurrentDeploymentRequest)(nil), "signalcd.GetCurrentDeploymentRequest")
	proto.RegisterType((*GetCurrentDeploymentResponse)(nil), "signalcd.GetCurrentDeploymentResponse")
	proto.RegisterType((*SetCurrentDeploymentRequest)(nil), "signalcd.SetCurrentDeploymentRequest")
	proto.RegisterType((*SetCurrentDeploymentResponse)(nil), "signalcd.SetCurrentDeploymentResponse")
	proto.RegisterType((*ListPipelinesRequest)(nil), "signalcd.ListPipelinesRequest")
	proto.RegisterType((*ListPipelinesResponse)(nil), "signalcd.ListPipelinesResponse")
	proto.RegisterType((*CreatePipelineRequest)(nil), "signalcd.CreatePipelineRequest")
	proto.RegisterType((*CreatePipelineResponse)(nil), "signalcd.CreatePipelineResponse")
	proto.RegisterType((*GetPipelineRequest)(nil), "signalcd.GetPipelineRequest")
	proto.RegisterType((*GetPipelineResponse)(nil), "signalcd.GetPipelineResponse")
}

func init() { proto.RegisterFile("signalcd/proto/ui.proto", fileDescriptor_df79a9f155088a03) }

var fileDescriptor_df79a9f155088a03 = []byte{
	// 496 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x86, 0xe5, 0x20, 0x50, 0x33, 0x69, 0x53, 0x31, 0xb8, 0x49, 0xe5, 0xa6, 0x34, 0x5a, 0x01,
	0xe2, 0x42, 0x8c, 0x0a, 0xe2, 0xd0, 0x23, 0x05, 0x85, 0x4a, 0x1c, 0xaa, 0x18, 0x2e, 0xdc, 0x4c,
	0x3c, 0x8a, 0x56, 0x0a, 0xb6, 0xeb, 0xdd, 0x20, 0x55, 0x88, 0x0b, 0x1c, 0x78, 0x00, 0x5e, 0x89,
	0x37, 0xe0, 0x15, 0x78, 0x10, 0xd4, 0xcd, 0xae, 0xbd, 0xb6, 0xd6, 0x0d, 0x20, 0x4e, 0x89, 0x76,
	0xfe, 0x9d, 0xef, 0xdf, 0xd9, 0x7f, 0x0d, 0x43, 0xc1, 0x17, 0x69, 0xbc, 0x9c, 0x27, 0x61, 0x5e,
	0x64, 0x32, 0x0b, 0x57, 0x7c, 0xa2, 0xfe, 0xe0, 0x96, 0x29, 0x04, 0xa3, 0x45, 0x96, 0x2d, 0x96,
	0x14, 0xc6, 0x39, 0x0f, 0xe3, 0x34, 0xcd, 0x64, 0x2c, 0x79, 0x96, 0x8a, 0xb5, 0x2e, 0x08, 0x1a,
	0x0d, 0xe4, 0x65, 0x4e, 0xba, 0xc6, 0x86, 0xb0, 0xf7, 0x9a, 0x0b, 0xf9, 0x82, 0xf2, 0x65, 0x76,
	0xf9, 0x81, 0x52, 0x39, 0xa3, 0x8b, 0x15, 0x09, 0xc9, 0xce, 0x61, 0xd0, 0x2c, 0x88, 0x3c, 0x4b,
	0x05, 0xe1, 0x33, 0xe8, 0x25, 0xe5, 0xaa, 0xd8, 0xf7, 0xc6, 0x37, 0x1e, 0xf6, 0x8e, 0xfd, 0x89,
	0x81, 0x4c, 0xac, 0x2d, 0xb6, 0x90, 0x1d, 0xc2, 0xc1, 0x94, 0xe4, 0xe9, 0xaa, 0x28, 0x28, 0x75,
	0x00, 0xdf, 0xc0, 0xc8, 0x5d, 0xd6, 0xd8, 0xa7, 0x00, 0x55, 0xb7, 0x7d, 0x6f, 0xec, 0xb5, 0x52,
	0x2d, 0x1d, 0x7b, 0x04, 0x07, 0x51, 0x3b, 0x14, 0xfb, 0xd0, 0xe1, 0x89, 0x6a, 0xd6, 0x9d, 0x75,
	0x78, 0x72, 0x65, 0x22, 0xfa, 0xff, 0x26, 0x06, 0xe0, 0x5f, 0xcd, 0xf2, 0x9c, 0xe7, 0xb4, 0xe4,
	0x29, 0x09, 0x73, 0xe4, 0xb3, 0xf5, 0xf0, 0xad, 0x75, 0x8d, 0x79, 0x0c, 0xdd, 0xdc, 0x2c, 0xea,
	0x01, 0x63, 0x45, 0x31, 0xfa, 0x59, 0x25, 0x62, 0x53, 0xd8, 0x3b, 0x2d, 0x28, 0x96, 0x54, 0x16,
	0xf5, 0x09, 0x27, 0xb0, 0x65, 0x54, 0xda, 0xaf, 0xab, 0x53, 0xa9, 0x61, 0xaf, 0x60, 0xd0, 0x6c,
	0xa4, 0x4d, 0xfd, 0x6d, 0xa7, 0x7b, 0x80, 0x53, 0x92, 0x4d, 0x3f, 0xcd, 0x89, 0xbf, 0x84, 0x3b,
	0x35, 0xd5, 0xbf, 0xc1, 0x8e, 0x7f, 0xdc, 0x84, 0xee, 0xdb, 0xb3, 0x88, 0x8a, 0x8f, 0x7c, 0x4e,
	0xc8, 0xa1, 0x5f, 0x0f, 0x2f, 0x1e, 0x55, 0xbb, 0x9d, 0x79, 0x0f, 0xc6, 0xed, 0x82, 0xb5, 0x25,
	0xe6, 0x7f, 0xf9, 0xf9, 0xeb, 0x7b, 0xa7, 0x8f, 0xdb, 0xa1, 0x95, 0x6a, 0xfc, 0xea, 0x81, 0xef,
	0xca, 0x2d, 0xde, 0xaf, 0x1a, 0x5e, 0x13, 0xfb, 0xe0, 0xc1, 0x26, 0x99, 0xa6, 0x8f, 0x14, 0x7d,
	0x80, 0xbe, 0x4d, 0x0f, 0xe7, 0x6b, 0x3d, 0x7e, 0xf3, 0xc0, 0x8f, 0x36, 0xb8, 0x88, 0xfe, 0xcc,
	0xc5, 0x75, 0xf9, 0x67, 0x63, 0xe5, 0x22, 0x60, 0x4e, 0x17, 0x27, 0x1d, 0x9e, 0x60, 0x02, 0x3b,
	0xb5, 0x4c, 0xe3, 0xdd, 0xfa, 0x60, 0x9b, 0x8f, 0x20, 0x38, 0x6a, 0xad, 0x6b, 0x26, 0x2a, 0xe6,
	0x36, 0x42, 0x58, 0xc6, 0x1d, 0x2f, 0xa0, 0x5f, 0x4f, 0xa9, 0x7d, 0xc1, 0xce, 0x87, 0x60, 0x5f,
	0xb0, 0x3b, 0xe0, 0x66, 0xc4, 0xcc, 0x02, 0x9d, 0x94, 0x09, 0xc3, 0x39, 0xf4, 0xac, 0xa0, 0xe2,
	0xa8, 0x76, 0x6f, 0x4d, 0xd8, 0x61, 0x4b, 0x55, 0x93, 0x86, 0x8a, 0x74, 0x1b, 0x77, 0x2b, 0x52,
	0xf8, 0x89, 0x27, 0x9f, 0x9f, 0xef, 0xbe, 0xdb, 0x31, 0x1b, 0xd5, 0xf7, 0xf9, 0xfd, 0x2d, 0xf5,
	0xf3, 0xe4, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd2, 0xb0, 0x2e, 0x73, 0x05, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// UIServiceClient is the client API for UIService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UIServiceClient interface {
	ListDeployment(ctx context.Context, in *ListDeploymentRequest, opts ...grpc.CallOption) (*ListDeploymentResponse, error)
	GetCurrentDeployment(ctx context.Context, in *GetCurrentDeploymentRequest, opts ...grpc.CallOption) (*GetCurrentDeploymentResponse, error)
	SetCurrentDeployment(ctx context.Context, in *SetCurrentDeploymentRequest, opts ...grpc.CallOption) (*SetCurrentDeploymentResponse, error)
	ListPipelines(ctx context.Context, in *ListPipelinesRequest, opts ...grpc.CallOption) (*ListPipelinesResponse, error)
	CreatePipeline(ctx context.Context, in *CreatePipelineRequest, opts ...grpc.CallOption) (*CreatePipelineResponse, error)
	GetPipeline(ctx context.Context, in *GetPipelineRequest, opts ...grpc.CallOption) (*GetPipelineResponse, error)
}

type uIServiceClient struct {
	cc *grpc.ClientConn
}

func NewUIServiceClient(cc *grpc.ClientConn) UIServiceClient {
	return &uIServiceClient{cc}
}

func (c *uIServiceClient) ListDeployment(ctx context.Context, in *ListDeploymentRequest, opts ...grpc.CallOption) (*ListDeploymentResponse, error) {
	out := new(ListDeploymentResponse)
	err := c.cc.Invoke(ctx, "/signalcd.UIService/ListDeployment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uIServiceClient) GetCurrentDeployment(ctx context.Context, in *GetCurrentDeploymentRequest, opts ...grpc.CallOption) (*GetCurrentDeploymentResponse, error) {
	out := new(GetCurrentDeploymentResponse)
	err := c.cc.Invoke(ctx, "/signalcd.UIService/GetCurrentDeployment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uIServiceClient) SetCurrentDeployment(ctx context.Context, in *SetCurrentDeploymentRequest, opts ...grpc.CallOption) (*SetCurrentDeploymentResponse, error) {
	out := new(SetCurrentDeploymentResponse)
	err := c.cc.Invoke(ctx, "/signalcd.UIService/SetCurrentDeployment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uIServiceClient) ListPipelines(ctx context.Context, in *ListPipelinesRequest, opts ...grpc.CallOption) (*ListPipelinesResponse, error) {
	out := new(ListPipelinesResponse)
	err := c.cc.Invoke(ctx, "/signalcd.UIService/ListPipelines", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uIServiceClient) CreatePipeline(ctx context.Context, in *CreatePipelineRequest, opts ...grpc.CallOption) (*CreatePipelineResponse, error) {
	out := new(CreatePipelineResponse)
	err := c.cc.Invoke(ctx, "/signalcd.UIService/CreatePipeline", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uIServiceClient) GetPipeline(ctx context.Context, in *GetPipelineRequest, opts ...grpc.CallOption) (*GetPipelineResponse, error) {
	out := new(GetPipelineResponse)
	err := c.cc.Invoke(ctx, "/signalcd.UIService/GetPipeline", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UIServiceServer is the server API for UIService service.
type UIServiceServer interface {
	ListDeployment(context.Context, *ListDeploymentRequest) (*ListDeploymentResponse, error)
	GetCurrentDeployment(context.Context, *GetCurrentDeploymentRequest) (*GetCurrentDeploymentResponse, error)
	SetCurrentDeployment(context.Context, *SetCurrentDeploymentRequest) (*SetCurrentDeploymentResponse, error)
	ListPipelines(context.Context, *ListPipelinesRequest) (*ListPipelinesResponse, error)
	CreatePipeline(context.Context, *CreatePipelineRequest) (*CreatePipelineResponse, error)
	GetPipeline(context.Context, *GetPipelineRequest) (*GetPipelineResponse, error)
}

// UnimplementedUIServiceServer can be embedded to have forward compatible implementations.
type UnimplementedUIServiceServer struct {
}

func (*UnimplementedUIServiceServer) ListDeployment(ctx context.Context, req *ListDeploymentRequest) (*ListDeploymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDeployment not implemented")
}
func (*UnimplementedUIServiceServer) GetCurrentDeployment(ctx context.Context, req *GetCurrentDeploymentRequest) (*GetCurrentDeploymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCurrentDeployment not implemented")
}
func (*UnimplementedUIServiceServer) SetCurrentDeployment(ctx context.Context, req *SetCurrentDeploymentRequest) (*SetCurrentDeploymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetCurrentDeployment not implemented")
}
func (*UnimplementedUIServiceServer) ListPipelines(ctx context.Context, req *ListPipelinesRequest) (*ListPipelinesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPipelines not implemented")
}
func (*UnimplementedUIServiceServer) CreatePipeline(ctx context.Context, req *CreatePipelineRequest) (*CreatePipelineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePipeline not implemented")
}
func (*UnimplementedUIServiceServer) GetPipeline(ctx context.Context, req *GetPipelineRequest) (*GetPipelineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPipeline not implemented")
}

func RegisterUIServiceServer(s *grpc.Server, srv UIServiceServer) {
	s.RegisterService(&_UIService_serviceDesc, srv)
}

func _UIService_ListDeployment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDeploymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UIServiceServer).ListDeployment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/signalcd.UIService/ListDeployment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UIServiceServer).ListDeployment(ctx, req.(*ListDeploymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UIService_GetCurrentDeployment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCurrentDeploymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UIServiceServer).GetCurrentDeployment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/signalcd.UIService/GetCurrentDeployment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UIServiceServer).GetCurrentDeployment(ctx, req.(*GetCurrentDeploymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UIService_SetCurrentDeployment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetCurrentDeploymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UIServiceServer).SetCurrentDeployment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/signalcd.UIService/SetCurrentDeployment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UIServiceServer).SetCurrentDeployment(ctx, req.(*SetCurrentDeploymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UIService_ListPipelines_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPipelinesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UIServiceServer).ListPipelines(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/signalcd.UIService/ListPipelines",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UIServiceServer).ListPipelines(ctx, req.(*ListPipelinesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UIService_CreatePipeline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePipelineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UIServiceServer).CreatePipeline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/signalcd.UIService/CreatePipeline",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UIServiceServer).CreatePipeline(ctx, req.(*CreatePipelineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UIService_GetPipeline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPipelineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UIServiceServer).GetPipeline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/signalcd.UIService/GetPipeline",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UIServiceServer).GetPipeline(ctx, req.(*GetPipelineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UIService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "signalcd.UIService",
	HandlerType: (*UIServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListDeployment",
			Handler:    _UIService_ListDeployment_Handler,
		},
		{
			MethodName: "GetCurrentDeployment",
			Handler:    _UIService_GetCurrentDeployment_Handler,
		},
		{
			MethodName: "SetCurrentDeployment",
			Handler:    _UIService_SetCurrentDeployment_Handler,
		},
		{
			MethodName: "ListPipelines",
			Handler:    _UIService_ListPipelines_Handler,
		},
		{
			MethodName: "CreatePipeline",
			Handler:    _UIService_CreatePipeline_Handler,
		},
		{
			MethodName: "GetPipeline",
			Handler:    _UIService_GetPipeline_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "signalcd/proto/ui.proto",
}