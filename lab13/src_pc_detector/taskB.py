import socket
import tkinter as tk
from tkinter import ttk
from scapy.all import ARP, Ether, srp

def scan(subnet):
    global progress_bar
    devices = []
    res = srp(Ether(dst="ff:ff:ff:ff:ff:ff")/ARP(pdst=subnet), timeout=1)[0]
    for _, received in res:
        name = "*"
        progress_bar['value'] += 100 / len(res)
        try:
            name = socket.gethostbyaddr(received.psrc)[0]
        except:
            name = "*"
        devices.append({'ip': received.psrc, 'mac': received.hwsrc, 'name': name})
    return devices

def start_scan():
    global tree
    current = socket.gethostbyname(socket.gethostname()) 
    devices = scan(current + "/24")

    for device in devices:
        if (device['ip'] == current):
            tree.insert("", tk.END, values = (device['ip'], device['mac'], device['name']))
    for device in devices:
        if (device['ip'] != current):
            tree.insert("", tk.END, values = (device['ip'], device['mac'], device['name']))

root = tk.Tk()
root.title("PCs Scanner")

progress_bar = ttk.Progressbar(root, orient="horizontal", length=200, mode="determinate")
progress_bar.grid(row=0, column=2)

scan_button = tk.Button(root, text="Start Scan", command=start_scan)
scan_button.grid(row=1, column=2)

tree = ttk.Treeview(root, columns=("IP", "MAC", "Name"), show="headings")
tree.grid(row=2, column=1, columnspan=3)
tree.heading("IP", text="IP")
tree.heading("MAC", text="MAC")
tree.heading("Name", text="Name")


root.mainloop()
