// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: utils/pb/suggestions/suggestions.proto

package suggestions

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	BookSuggestionsService_SuggestBooks_FullMethodName = "/transaction.BookSuggestionsService/suggestBooks"
)

// BookSuggestionsServiceClient is the client API for BookSuggestionsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BookSuggestionsServiceClient interface {
	SuggestBooks(ctx context.Context, in *ItemsBought, opts ...grpc.CallOption) (*BookSuggest, error)
}

type bookSuggestionsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBookSuggestionsServiceClient(cc grpc.ClientConnInterface) BookSuggestionsServiceClient {
	return &bookSuggestionsServiceClient{cc}
}

func (c *bookSuggestionsServiceClient) SuggestBooks(ctx context.Context, in *ItemsBought, opts ...grpc.CallOption) (*BookSuggest, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BookSuggest)
	err := c.cc.Invoke(ctx, BookSuggestionsService_SuggestBooks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookSuggestionsServiceServer is the server API for BookSuggestionsService service.
// All implementations must embed UnimplementedBookSuggestionsServiceServer
// for forward compatibility.
type BookSuggestionsServiceServer interface {
	SuggestBooks(context.Context, *ItemsBought) (*BookSuggest, error)
	mustEmbedUnimplementedBookSuggestionsServiceServer()
}

// UnimplementedBookSuggestionsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBookSuggestionsServiceServer struct{}

func (UnimplementedBookSuggestionsServiceServer) SuggestBooks(context.Context, *ItemsBought) (*BookSuggest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SuggestBooks not implemented")
}
func (UnimplementedBookSuggestionsServiceServer) mustEmbedUnimplementedBookSuggestionsServiceServer() {
}
func (UnimplementedBookSuggestionsServiceServer) testEmbeddedByValue() {}

// UnsafeBookSuggestionsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BookSuggestionsServiceServer will
// result in compilation errors.
type UnsafeBookSuggestionsServiceServer interface {
	mustEmbedUnimplementedBookSuggestionsServiceServer()
}

func RegisterBookSuggestionsServiceServer(s grpc.ServiceRegistrar, srv BookSuggestionsServiceServer) {
	// If the following call pancis, it indicates UnimplementedBookSuggestionsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BookSuggestionsService_ServiceDesc, srv)
}

func _BookSuggestionsService_SuggestBooks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ItemsBought)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookSuggestionsServiceServer).SuggestBooks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookSuggestionsService_SuggestBooks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookSuggestionsServiceServer).SuggestBooks(ctx, req.(*ItemsBought))
	}
	return interceptor(ctx, in, info, handler)
}

// BookSuggestionsService_ServiceDesc is the grpc.ServiceDesc for BookSuggestionsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BookSuggestionsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "transaction.BookSuggestionsService",
	HandlerType: (*BookSuggestionsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "suggestBooks",
			Handler:    _BookSuggestionsService_SuggestBooks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "utils/pb/suggestions/suggestions.proto",
}
