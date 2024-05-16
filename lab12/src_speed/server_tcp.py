import socket
import struct
from datetime import datetime

import tkinter as tk

def recieve_packets(host, port, speed_text, loss_text):
    tcp_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    tcp_sock.bind((host, port))
    tcp_sock.listen()
    conn, _ = tcp_sock.accept()
    timer_and_amount = conn.recv(12)
    amount = int.from_bytes(timer_and_amount[:4], "little")
    start = struct.unpack_from("d", timer_and_amount[4:])[0]
    for _ in range(amount):
        conn.recv(1024)

    speed_text['text'] = f"{round(amount * 1024 / ([datetime.now().timestamp()][0] - start))} B/S"
    loss_text['text'] = f"{amount} of {amount}"

    conn.close()
    tcp_sock.close()

if __name__ == "__main__":
    window = tk.Tk()
    window.title("Server TCP")
    
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
