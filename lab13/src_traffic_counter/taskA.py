from scapy.all import *

total_in = 0
total_out = 0

def handler(packet):
    global total_in
    global total_out
    if IP in packet:
        src_ip = packet[IP].src
        dst_ip = packet[IP].dst
        src_port = packet[IP].sport
        dst_port = packet[IP].dport
        size = len(packet)
        if (src_ip == "192.168.1.107"):
            total_out += size
        else:
            total_in += size
        print(f"Source IP: {src_ip}:{src_port}, Destination IP: {dst_ip}:{dst_port} Size: {size}")
        print(f"Total in: {total_in}, Total out: {total_out}")
        print()

sniff(prn=handler, store=0)
