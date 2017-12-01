// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/containerd/containerd/api/services/introspection/v1/introspection.proto

/*
	Package introspection is a generated protocol buffer package.

	It is generated from these files:
		github.com/containerd/containerd/api/services/introspection/v1/introspection.proto

	It has these top-level messages:
		Plugin
		PluginsRequest
		PluginsResponse
*/
package introspection

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import containerd_types "github.com/containerd/containerd/api/types"
import google_rpc "github.com/containerd/containerd/protobuf/google/rpc"

// skipping weak import gogoproto "github.com/gogo/protobuf/gogoproto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

import strings "strings"
import reflect "reflect"
import github_com_gogo_protobuf_sortkeys "github.com/gogo/protobuf/sortkeys"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Plugin struct {
	// Type defines the type of plugin.
	//
	// See package plugin for a list of possible values. Non core plugins may
	// define their own values during registration.
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// ID identifies the plugin uniquely in the system.
	ID string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	// Requires lists the plugin types required by this plugin.
	Requires []string `protobuf:"bytes,3,rep,name=requires" json:"requires,omitempty"`
	// Platforms enumerates the platforms this plugin will support.
	//
	// If values are provided here, the plugin will only be operable under the
	// provided platforms.
	//
	// If this is empty, the plugin will work across all platforms.
	//
	// If the plugin prefers certain platforms over others, they should be
	// listed from most to least preferred.
	Platforms []containerd_types.Platform `protobuf:"bytes,4,rep,name=platforms" json:"platforms"`
	// Exports allows plugins to provide values about state or configuration to
	// interested parties.
	//
	// One example is exposing the configured path of a snapshotter plugin.
	Exports map[string]string `protobuf:"bytes,5,rep,name=exports" json:"exports,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Capabilities allows plugins to communicate feature switches to allow
	// clients to detect features that may not be on be default or may be
	// different from version to version.
	//
	// Use this sparingly.
	Capabilities []string `protobuf:"bytes,6,rep,name=capabilities" json:"capabilities,omitempty"`
	// InitErr will be set if the plugin fails initialization.
	//
	// This means the plugin may have been registered but a non-terminal error
	// was encountered during initialization.
	//
	// Plugins that have this value set cannot be used.
	InitErr *google_rpc.Status `protobuf:"bytes,7,opt,name=init_err,json=initErr" json:"init_err,omitempty"`
}

func (m *Plugin) Reset()                    { *m = Plugin{} }
func (*Plugin) ProtoMessage()               {}
func (*Plugin) Descriptor() ([]byte, []int) { return fileDescriptorIntrospection, []int{0} }

type PluginsRequest struct {
	// Filters contains one or more filters using the syntax defined in the
	// containerd filter package.
	//
	// The returned result will be those that match any of the provided
	// filters. Expanded, plugins that match the following will be
	// returned:
	//
	//   filters[0] or filters[1] or ... or filters[n-1] or filters[n]
	//
	// If filters is zero-length or nil, all items will be returned.
	Filters []string `protobuf:"bytes,1,rep,name=filters" json:"filters,omitempty"`
}

func (m *PluginsRequest) Reset()                    { *m = PluginsRequest{} }
func (*PluginsRequest) ProtoMessage()               {}
func (*PluginsRequest) Descriptor() ([]byte, []int) { return fileDescriptorIntrospection, []int{1} }

type PluginsResponse struct {
	Plugins []Plugin `protobuf:"bytes,1,rep,name=plugins" json:"plugins"`
}

func (m *PluginsResponse) Reset()                    { *m = PluginsResponse{} }
func (*PluginsResponse) ProtoMessage()               {}
func (*PluginsResponse) Descriptor() ([]byte, []int) { return fileDescriptorIntrospection, []int{2} }

func init() {
	proto.RegisterType((*Plugin)(nil), "containerd.services.introspection.v1.Plugin")
	proto.RegisterType((*PluginsRequest)(nil), "containerd.services.introspection.v1.PluginsRequest")
	proto.RegisterType((*PluginsResponse)(nil), "containerd.services.introspection.v1.PluginsResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Introspection service

type IntrospectionClient interface {
	// Plugins returns a list of plugins in containerd.
	//
	// Clients can use this to detect features and capabilities when using
	// containerd.
	Plugins(ctx context.Context, in *PluginsRequest, opts ...grpc.CallOption) (*PluginsResponse, error)
}

type introspectionClient struct {
	cc *grpc.ClientConn
}

func NewIntrospectionClient(cc *grpc.ClientConn) IntrospectionClient {
	return &introspectionClient{cc}
}

func (c *introspectionClient) Plugins(ctx context.Context, in *PluginsRequest, opts ...grpc.CallOption) (*PluginsResponse, error) {
	out := new(PluginsResponse)
	err := grpc.Invoke(ctx, "/containerd.services.introspection.v1.Introspection/Plugins", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Introspection service

type IntrospectionServer interface {
	// Plugins returns a list of plugins in containerd.
	//
	// Clients can use this to detect features and capabilities when using
	// containerd.
	Plugins(context.Context, *PluginsRequest) (*PluginsResponse, error)
}

func RegisterIntrospectionServer(s *grpc.Server, srv IntrospectionServer) {
	s.RegisterService(&_Introspection_serviceDesc, srv)
}

func _Introspection_Plugins_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PluginsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IntrospectionServer).Plugins(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/containerd.services.introspection.v1.Introspection/Plugins",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IntrospectionServer).Plugins(ctx, req.(*PluginsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Introspection_serviceDesc = grpc.ServiceDesc{
	ServiceName: "containerd.services.introspection.v1.Introspection",
	HandlerType: (*IntrospectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Plugins",
			Handler:    _Introspection_Plugins_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/containerd/containerd/api/services/introspection/v1/introspection.proto",
}

func (m *Plugin) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Plugin) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Type) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintIntrospection(dAtA, i, uint64(len(m.Type)))
		i += copy(dAtA[i:], m.Type)
	}
	if len(m.ID) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintIntrospection(dAtA, i, uint64(len(m.ID)))
		i += copy(dAtA[i:], m.ID)
	}
	if len(m.Requires) > 0 {
		for _, s := range m.Requires {
			dAtA[i] = 0x1a
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	if len(m.Platforms) > 0 {
		for _, msg := range m.Platforms {
			dAtA[i] = 0x22
			i++
			i = encodeVarintIntrospection(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Exports) > 0 {
		for k, _ := range m.Exports {
			dAtA[i] = 0x2a
			i++
			v := m.Exports[k]
			mapSize := 1 + len(k) + sovIntrospection(uint64(len(k))) + 1 + len(v) + sovIntrospection(uint64(len(v)))
			i = encodeVarintIntrospection(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintIntrospection(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintIntrospection(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if len(m.Capabilities) > 0 {
		for _, s := range m.Capabilities {
			dAtA[i] = 0x32
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	if m.InitErr != nil {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintIntrospection(dAtA, i, uint64(m.InitErr.Size()))
		n1, err := m.InitErr.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func (m *PluginsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PluginsRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Filters) > 0 {
		for _, s := range m.Filters {
			dAtA[i] = 0xa
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	return i, nil
}

func (m *PluginsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PluginsResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Plugins) > 0 {
		for _, msg := range m.Plugins {
			dAtA[i] = 0xa
			i++
			i = encodeVarintIntrospection(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeVarintIntrospection(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Plugin) Size() (n int) {
	var l int
	_ = l
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovIntrospection(uint64(l))
	}
	l = len(m.ID)
	if l > 0 {
		n += 1 + l + sovIntrospection(uint64(l))
	}
	if len(m.Requires) > 0 {
		for _, s := range m.Requires {
			l = len(s)
			n += 1 + l + sovIntrospection(uint64(l))
		}
	}
	if len(m.Platforms) > 0 {
		for _, e := range m.Platforms {
			l = e.Size()
			n += 1 + l + sovIntrospection(uint64(l))
		}
	}
	if len(m.Exports) > 0 {
		for k, v := range m.Exports {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovIntrospection(uint64(len(k))) + 1 + len(v) + sovIntrospection(uint64(len(v)))
			n += mapEntrySize + 1 + sovIntrospection(uint64(mapEntrySize))
		}
	}
	if len(m.Capabilities) > 0 {
		for _, s := range m.Capabilities {
			l = len(s)
			n += 1 + l + sovIntrospection(uint64(l))
		}
	}
	if m.InitErr != nil {
		l = m.InitErr.Size()
		n += 1 + l + sovIntrospection(uint64(l))
	}
	return n
}

func (m *PluginsRequest) Size() (n int) {
	var l int
	_ = l
	if len(m.Filters) > 0 {
		for _, s := range m.Filters {
			l = len(s)
			n += 1 + l + sovIntrospection(uint64(l))
		}
	}
	return n
}

func (m *PluginsResponse) Size() (n int) {
	var l int
	_ = l
	if len(m.Plugins) > 0 {
		for _, e := range m.Plugins {
			l = e.Size()
			n += 1 + l + sovIntrospection(uint64(l))
		}
	}
	return n
}

func sovIntrospection(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozIntrospection(x uint64) (n int) {
	return sovIntrospection(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Plugin) String() string {
	if this == nil {
		return "nil"
	}
	keysForExports := make([]string, 0, len(this.Exports))
	for k, _ := range this.Exports {
		keysForExports = append(keysForExports, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForExports)
	mapStringForExports := "map[string]string{"
	for _, k := range keysForExports {
		mapStringForExports += fmt.Sprintf("%v: %v,", k, this.Exports[k])
	}
	mapStringForExports += "}"
	s := strings.Join([]string{`&Plugin{`,
		`Type:` + fmt.Sprintf("%v", this.Type) + `,`,
		`ID:` + fmt.Sprintf("%v", this.ID) + `,`,
		`Requires:` + fmt.Sprintf("%v", this.Requires) + `,`,
		`Platforms:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Platforms), "Platform", "containerd_types.Platform", 1), `&`, ``, 1) + `,`,
		`Exports:` + mapStringForExports + `,`,
		`Capabilities:` + fmt.Sprintf("%v", this.Capabilities) + `,`,
		`InitErr:` + strings.Replace(fmt.Sprintf("%v", this.InitErr), "Status", "google_rpc.Status", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *PluginsRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&PluginsRequest{`,
		`Filters:` + fmt.Sprintf("%v", this.Filters) + `,`,
		`}`,
	}, "")
	return s
}
func (this *PluginsResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&PluginsResponse{`,
		`Plugins:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Plugins), "Plugin", "Plugin", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringIntrospection(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Plugin) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIntrospection
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Plugin: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Plugin: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntrospection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntrospection
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntrospection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntrospection
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Requires", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntrospection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntrospection
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Requires = append(m.Requires, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Platforms", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntrospection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIntrospection
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Platforms = append(m.Platforms, containerd_types.Platform{})
			if err := m.Platforms[len(m.Platforms)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Exports", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntrospection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIntrospection
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Exports == nil {
				m.Exports = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowIntrospection
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowIntrospection
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthIntrospection
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowIntrospection
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthIntrospection
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipIntrospection(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthIntrospection
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Exports[mapkey] = mapvalue
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Capabilities", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntrospection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntrospection
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Capabilities = append(m.Capabilities, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitErr", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntrospection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIntrospection
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.InitErr == nil {
				m.InitErr = &google_rpc.Status{}
			}
			if err := m.InitErr.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIntrospection(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIntrospection
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PluginsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIntrospection
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PluginsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PluginsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Filters", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntrospection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIntrospection
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Filters = append(m.Filters, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIntrospection(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIntrospection
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PluginsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIntrospection
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PluginsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PluginsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Plugins", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIntrospection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIntrospection
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Plugins = append(m.Plugins, Plugin{})
			if err := m.Plugins[len(m.Plugins)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIntrospection(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIntrospection
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipIntrospection(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowIntrospection
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowIntrospection
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowIntrospection
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthIntrospection
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowIntrospection
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipIntrospection(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthIntrospection = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowIntrospection   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/containerd/containerd/api/services/introspection/v1/introspection.proto", fileDescriptorIntrospection)
}

var fileDescriptorIntrospection = []byte{
	// 487 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xcd, 0x3a, 0x69, 0xdc, 0x4c, 0xca, 0x87, 0x56, 0x15, 0x58, 0x3e, 0xb8, 0x51, 0xc4, 0x21,
	0x42, 0xb0, 0x56, 0x03, 0x48, 0xb4, 0x48, 0x1c, 0x22, 0x72, 0xa8, 0xd4, 0x43, 0xe5, 0x5e, 0x10,
	0x97, 0xca, 0x71, 0x36, 0x66, 0x85, 0xeb, 0xdd, 0xee, 0xae, 0x2d, 0x72, 0xe3, 0xc6, 0x5f, 0xcb,
	0x91, 0x23, 0xa7, 0x8a, 0xfa, 0x37, 0xf0, 0x03, 0x90, 0xbd, 0x76, 0x9b, 0xdc, 0x12, 0x71, 0x9b,
	0x79, 0x7e, 0x6f, 0xe6, 0xcd, 0x93, 0x17, 0x82, 0x98, 0xe9, 0xaf, 0xd9, 0x8c, 0x44, 0xfc, 0xda,
	0x8f, 0x78, 0xaa, 0x43, 0x96, 0x52, 0x39, 0x5f, 0x2f, 0x43, 0xc1, 0x7c, 0x45, 0x65, 0xce, 0x22,
	0xaa, 0x7c, 0x96, 0x6a, 0xc9, 0x95, 0xa0, 0x91, 0x66, 0x3c, 0xf5, 0xf3, 0xe3, 0x4d, 0x80, 0x08,
	0xc9, 0x35, 0xc7, 0x2f, 0x1e, 0xd4, 0xa4, 0x51, 0x92, 0x4d, 0x62, 0x7e, 0xec, 0x9e, 0x6c, 0xb5,
	0x59, 0x2f, 0x05, 0x55, 0xbe, 0x48, 0x42, 0xbd, 0xe0, 0xf2, 0xda, 0x2c, 0x70, 0x9f, 0xc7, 0x9c,
	0xc7, 0x09, 0xf5, 0xa5, 0x88, 0x7c, 0xa5, 0x43, 0x9d, 0xa9, 0xfa, 0xc3, 0x61, 0xcc, 0x63, 0x5e,
	0x95, 0x7e, 0x59, 0x19, 0x74, 0xf8, 0xd7, 0x82, 0xee, 0x45, 0x92, 0xc5, 0x2c, 0xc5, 0x18, 0x3a,
	0xe5, 0x44, 0x07, 0x0d, 0xd0, 0xa8, 0x17, 0x54, 0x35, 0x7e, 0x06, 0x16, 0x9b, 0x3b, 0x56, 0x89,
	0x4c, 0xba, 0xc5, 0xed, 0x91, 0x75, 0xf6, 0x29, 0xb0, 0xd8, 0x1c, 0xbb, 0xb0, 0x2f, 0xe9, 0x4d,
	0xc6, 0x24, 0x55, 0x4e, 0x7b, 0xd0, 0x1e, 0xf5, 0x82, 0xfb, 0x1e, 0x7f, 0x84, 0x5e, 0xe3, 0x49,
	0x39, 0x9d, 0x41, 0x7b, 0xd4, 0x1f, 0xbb, 0x64, 0xed, 0xec, 0xca, 0x36, 0xb9, 0xa8, 0x29, 0x93,
	0xce, 0xea, 0xf6, 0xa8, 0x15, 0x3c, 0x48, 0xf0, 0x25, 0xd8, 0xf4, 0xbb, 0xe0, 0x52, 0x2b, 0x67,
	0xaf, 0x52, 0x9f, 0x90, 0x6d, 0x42, 0x23, 0xe6, 0x0c, 0x32, 0x35, 0xda, 0x69, 0xaa, 0xe5, 0x32,
	0x68, 0x26, 0xe1, 0x21, 0x1c, 0x44, 0xa1, 0x08, 0x67, 0x2c, 0x61, 0x9a, 0x51, 0xe5, 0x74, 0x2b,
	0xd3, 0x1b, 0x18, 0x7e, 0x0d, 0xfb, 0x2c, 0x65, 0xfa, 0x8a, 0x4a, 0xe9, 0xd8, 0x03, 0x34, 0xea,
	0x8f, 0x31, 0x31, 0x69, 0x12, 0x29, 0x22, 0x72, 0x59, 0xa5, 0x19, 0xd8, 0x25, 0x67, 0x2a, 0xa5,
	0x7b, 0x0a, 0x07, 0xeb, 0xbb, 0xf0, 0x53, 0x68, 0x7f, 0xa3, 0xcb, 0x3a, 0xbe, 0xb2, 0xc4, 0x87,
	0xb0, 0x97, 0x87, 0x49, 0x46, 0x4d, 0x80, 0x81, 0x69, 0x4e, 0xad, 0xf7, 0x68, 0xf8, 0x12, 0x1e,
	0x1b, 0xbb, 0x2a, 0xa0, 0x37, 0x19, 0x55, 0x1a, 0x3b, 0x60, 0x2f, 0x58, 0xa2, 0xa9, 0x54, 0x0e,
	0xaa, 0xbc, 0x35, 0xed, 0xf0, 0x0a, 0x9e, 0xdc, 0x73, 0x95, 0xe0, 0xa9, 0xa2, 0xf8, 0x1c, 0x6c,
	0x61, 0xa0, 0x8a, 0xdc, 0x1f, 0xbf, 0xda, 0x25, 0xa2, 0x3a, 0xf2, 0x66, 0xc4, 0xf8, 0x27, 0x82,
	0x47, 0x67, 0xeb, 0x54, 0x9c, 0x83, 0x5d, 0xaf, 0xc4, 0x6f, 0x77, 0x99, 0xdc, 0x5c, 0xe3, 0xbe,
	0xdb, 0x51, 0x65, 0xee, 0x9a, 0x2c, 0x56, 0x77, 0x5e, 0xeb, 0xf7, 0x9d, 0xd7, 0xfa, 0x51, 0x78,
	0x68, 0x55, 0x78, 0xe8, 0x57, 0xe1, 0xa1, 0x3f, 0x85, 0x87, 0xbe, 0x9c, 0xff, 0xdf, 0x5b, 0xfc,
	0xb0, 0x01, 0x7c, 0xb6, 0x66, 0xdd, 0xea, 0xf7, 0x7f, 0xf3, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xe6,
	0x72, 0xde, 0x35, 0xe4, 0x03, 0x00, 0x00,
}
