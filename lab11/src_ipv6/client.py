import socket

sock = socket.socket(socket.AF_INET6, socket.SOCK_STREAM)
sock.setsockopt(socket.IPPROTO_IPV6, socket.IPV6_V6ONLY, 1)
sock.connect(("localhost", 8080))

data = input('Request: ').encode()
sock.sendall(data)
resp = sock.recv(len(data))
print('Respopnse: ', resp.decode())