import os
import socket
import struct
from datetime import datetime

import tkinter as tk

def send_packets(host, port, amount):
    tcp_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    tcp_sock.connect((host, port))
    
    tcp_sock.sendall(bytearray(amount.to_bytes(4, "little")) + bytearray(struct.pack("d", [datetime.now().timestamp()][0])))
    for _ in range(amount):
        tcp_sock.sendall(bytearray(os.urandom(1024)))
    tcp_sock.close()

if __name__ == "__main__":
    window = tk.Tk()
    window.title("Client TCP")
    
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
