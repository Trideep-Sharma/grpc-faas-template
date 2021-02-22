#!/usr/bin/env python

from function import generated_grpc
from function import handler
from concurrent import futures
import grpc

if __name__ == '__main__':
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    generated_grpc.Function_pb2_grpc.add_FunctionServiceServicer_to_server(handler.HandleGRPC(), server)
    server.add_insecure_port('[::]:8080')
    server.start()
    print("started grpc server")
    server.wait_for_termination()
