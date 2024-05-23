from scapy.all import *

total_data = {}
cnt = 0

def handler(packet):
    global cnt
    global total_data
    if IP in packet:
        cnt += 1
        src_ip = packet[IP].src
        dst_ip = packet[IP].dst
        src_port = packet[IP].sport
        dst_port = packet[IP].dport
        size = len(packet)
        if (src_ip == "192.168.1.107"):
            if not (src_port in total_data):
                total_data[src_port] = [0, 0]
            total_data[src_port][0] += size
        elif (dst_ip == "192.168.1.107"):
            if not (dst_port in total_data):
                total_data[dst_port] = [0, 0]
            total_data[dst_port][1] += size
        if (cnt % 100 == 0):
            for port in total_data:
                print(f"For {port} input size = {total_data[port][0]}, output size = {total_data[port][1]}")
            print()

sniff(prn=handler, store=0)
