from .generated_grpc import Function_pb2, Function_pb2_grpc
import grpc

class HandleGRPC(Function_pb2_grpc.FunctionServiceServicer):

    def getFunctionResponse(self, request, context):
        return Function_pb2.FunctionResponse(msg=f'You said: {request.msg}')
