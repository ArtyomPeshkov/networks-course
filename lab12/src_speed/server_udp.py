import socket
import struct
from datetime import datetime

import tkinter as tk

def recieve_packets(host, port, speed_text, loss_text):
    udp_sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    udp_sock.bind((host, port))
    udp_sock.settimeout(1)
    timer_and_amount = udp_sock.recv(12)
    amount = int.from_bytes(timer_and_amount[:4], "little")
    start = struct.unpack_from("d", timer_and_amount[4:])[0]
    recieved = amount

    for _ in range(amount):
        try:
            udp_sock.recv(1024)
        except socket.timeout:
            recieved -= 1

    speed_text['text'] = f"{round(amount * 1024 / ([datetime.now().timestamp()][0] - (amount - recieved) - start))} B/S"
    loss_text['text'] = f"{recieved} of {amount}"
    udp_sock.close()


if __name__ == "__main__":
    window = tk.Tk()
    window.title("Server UDP")
    
    host_label = tk.Label(window, text = "Host")
    host_label.pack()

    host_entry = tk.Entry(window)
    host_entry.pack()
    host_entry.insert(0, "127.0.0.1")

    port_label = tk.Label(window, text = "Port")
    port_label.pack()
    
    port_entry = tk.Entry(window)
    port_entry.pack()
    port_entry.insert(0, "8080")

    speed_label = tk.Label(window, text = "Speed")
    speed_label.pack()
    
    speed_text = tk.Label(window, text = "0 B/S")
    speed_text.pack()

    loss_label = tk.Label(window, text = "Amount of packets")
    loss_label.pack()
    
    loss_text = tk.Label(window, text = "0 of 0")
    loss_text.pack()
    
    send_button = tk.Button(window, text = "Recieve", command=lambda: recieve_packets(host_entry.get(), int(port_entry.get()), speed_text, loss_text))
    send_button.pack()
    
    window.mainloop()
