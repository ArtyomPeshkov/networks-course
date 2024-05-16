import pygame
import socket
import pickle

pygame.init()

screen = pygame.display.set_mode((800, 600))
pygame.display.set_caption("Remote Drawing Client")

client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client_socket.connect(('localhost', 8080))

drawing = False
last = None
while True:
    for event in pygame.event.get():
        if event.type == pygame.MOUSEBUTTONDOWN:
            drawing = True
        elif event.type == pygame.MOUSEBUTTONUP:
            drawing = False
            last = None
    
    if drawing:
        mouse_pos = pygame.mouse.get_pos()
        if last != None:
            pygame.draw.line(screen, (255, 0, 0), last, mouse_pos, 2)
            client_socket.send(pickle.dumps((last, mouse_pos)))
        last = mouse_pos

    pygame.display.flip()
