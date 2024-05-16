
import os
import socket
import struct
from datetime import datetime

import tkinter as tk

def send_packets(host, port, amount):

    udp_sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    udp_sock.sendto(bytearray(amount.to_bytes(4, "little")) + bytearray(struct.pack("d", [datetime.now().timestamp()][0])), (host, port))
    
    for _ in range(amount):
        udp_sock.sendto(bytearray(os.urandom(1024)), (host, port))
    udp_sock.close()

if __name__ == "__main__":
    window = tk.Tk()
    window.title("Client UDP")
    
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

    amount_label = tk.Label(window, text = "Amount")
    amount_label.pack()
    
    amount_entry = tk.Entry(window)
    amount_entry.pack()
    amount_entry.insert(0, "128")
    
    send_button = tk.Button(window, text="Send", command=lambda: send_packets(host_entry.get(), int(port_entry.get()), int(amount_entry.get())))
    send_button.pack()
    
    window.mainloop()
