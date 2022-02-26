/**
 * @fileoverview gRPC-Web generated client stub for example
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js')
const proto = {};
proto.example = require('./user.api_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.example.UserServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.example.UserServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.example.User,
 *   !proto.example.User>}
 */
const methodDescriptor_UserService_CreateUser = new grpc.web.MethodDescriptor(
  '/example.UserService/CreateUser',
  grpc.web.MethodType.UNARY,
  proto.example.User,
  proto.example.User,
  /**
   * @param {!proto.example.User} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.example.User.deserializeBinary
);


/**
 * @param {!proto.example.User} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.example.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.example.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.example.UserServiceClient.prototype.createUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/example.UserService/CreateUser',
      request,
      metadata || {},
      methodDescriptor_UserService_CreateUser,
      callback);
};


/**
 * @param {!proto.example.User} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.example.User>}
 *     Promise that resolves to the response
 */
proto.example.UserServicePromiseClient.prototype.createUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/example.UserService/CreateUser',
      request,
      metadata || {},
      methodDescriptor_UserService_CreateUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.example.UserId,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_UserService_DeleteUser = new grpc.web.MethodDescriptor(
  '/example.UserService/DeleteUser',
  grpc.web.MethodType.UNARY,
  proto.example.UserId,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.example.UserId} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.example.UserId} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.example.UserServiceClient.prototype.deleteUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/example.UserService/DeleteUser',
      request,
      metadata || {},
      methodDescriptor_UserService_DeleteUser,
      callback);
};


/**
 * @param {!proto.example.UserId} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     Promise that resolves to the response
 */
proto.example.UserServicePromiseClient.prototype.deleteUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/example.UserService/DeleteUser',
      request,
      metadata || {},
      methodDescriptor_UserService_DeleteUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.example.UpdateUserReq,
 *   !proto.example.User>}
 */
const methodDescriptor_UserService_UpdateUser = new grpc.web.MethodDescriptor(
  '/example.UserService/UpdateUser',
  grpc.web.MethodType.UNARY,
  proto.example.UpdateUserReq,
  proto.example.User,
  /**
   * @param {!proto.example.UpdateUserReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.example.User.deserializeBinary
);


/**
 * @param {!proto.example.UpdateUserReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.example.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.example.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.example.UserServiceClient.prototype.updateUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/example.UserService/UpdateUser',
      request,
      metadata || {},
      methodDescriptor_UserService_UpdateUser,
      callback);
};


/**
 * @param {!proto.example.UpdateUserReq} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.example.User>}
 *     Promise that resolves to the response
 */
proto.example.UserServicePromiseClient.prototype.updateUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/example.UserService/UpdateUser',
      request,
      metadata || {},
      methodDescriptor_UserService_UpdateUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.example.UserId,
 *   !proto.example.User>}
 */
const methodDescriptor_UserService_GetUser = new grpc.web.MethodDescriptor(
  '/example.UserService/GetUser',
  grpc.web.MethodType.UNARY,
  proto.example.UserId,
  proto.example.User,
  /**
   * @param {!proto.example.UserId} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.example.User.deserializeBinary
);


/**
 * @param {!proto.example.UserId} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.example.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.example.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.example.UserServiceClient.prototype.getUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/example.UserService/GetUser',
      request,
      metadata || {},
      methodDescriptor_UserService_GetUser,
      callback);
};


/**
 * @param {!proto.example.UserId} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.example.User>}
 *     Promise that resolves to the response
 */
proto.example.UserServicePromiseClient.prototype.getUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/example.UserService/GetUser',
      request,
      metadata || {},
      methodDescriptor_UserService_GetUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.example.ListUsersReq,
 *   !proto.example.ListUsersResp>}
 */
const methodDescriptor_UserService_ListUsers = new grpc.web.MethodDescriptor(
  '/example.UserService/ListUsers',
  grpc.web.MethodType.UNARY,
  proto.example.ListUsersReq,
  proto.example.ListUsersResp,
  /**
   * @param {!proto.example.ListUsersReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.example.ListUsersResp.deserializeBinary
);


/**
 * @param {!proto.example.ListUsersReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.example.ListUsersResp)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.example.ListUsersResp>|undefined}
 *     The XHR Node Readable Stream
 */
proto.example.UserServiceClient.prototype.listUsers =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/example.UserService/ListUsers',
      request,
      metadata || {},
      methodDescriptor_UserService_ListUsers,
      callback);
};


/**
 * @param {!proto.example.ListUsersReq} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.example.ListUsersResp>}
 *     Promise that resolves to the response
 */
proto.example.UserServicePromiseClient.prototype.listUsers =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/example.UserService/ListUsers',
      request,
      metadata || {},
      methodDescriptor_UserService_ListUsers);
};


module.exports = proto.example;

