import sys
from . import resources_pb2

input_bytes = sys.stdin.buffer.read()
message = resources_pb2.Person()
message.ParseFromString(input_bytes)
print(message)
