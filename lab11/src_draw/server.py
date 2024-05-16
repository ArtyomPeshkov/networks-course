import pygame
import socket
import pickle

pygame.init()
screen = pygame.display.set_mode((800, 600))
pygame.display.set_caption("Remote Drawing Host")

server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_socket.bind(('localhost', 8080))
server_socket.listen(1)
client_socket, address = server_socket.accept()

while True:
    pygame.event.pump()
    data = client_socket.recv(4096)
    if not data:
        break

    start, end = pickle.loads(data)
    pygame.draw.line(screen, (255, 0, 0), start, end, 2)

    pygame.display.flip()

pygame.quit()
