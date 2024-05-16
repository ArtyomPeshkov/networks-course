import socket

sock = socket.socket(socket.AF_INET6, socket.SOCK_STREAM) 
sock.setsockopt(socket.IPPROTO_IPV6, socket.IPV6_V6ONLY, 1)
sock.bind(("localhost", 8080))    
sock.listen()

while True:
    conn, _ = sock.accept()
    conn.sendall(conn.recv(2048).decode().upper().encode())