#!/usr/bin/env python
from flask import Flask
from waitress import serve

from function import generated_grpc
from function import handler
from concurrent import futures
import grpc

app = Flask(__name__)

@app.route('/_/health', methods=['GET', 'PUT', 'POST', 'PATCH', 'DELETE'])
def call_handler():
    return ("OK", 200)

if __name__ == '__main__':
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10), maximum_concurrent_rpcs=200)
    generated_grpc.Function_pb2_grpc.add_FunctionServiceServicer_to_server(handler.HandleGRPC(), server)
    server.add_insecure_port('[::]:9000')
    server.start()
    print("started grpc server")
    serve(app, host='0.0.0.0', port=8080)
    server.wait_for_termination()
