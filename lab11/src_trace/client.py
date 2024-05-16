import socket
import struct
import time
import sys

sock = socket.socket(socket.AF_INET, socket.SOCK_RAW, socket.IPPROTO_ICMP)
sock.settimeout(2)
pack_num = 0

def checksum_calc(num):
    checksum = (num >> 16) + (num & 0xFFFF)
    checksum += (checksum >> 16)
    checksum = ~checksum & 0xFFFF
    return checksum

def create_pack():
    global pack_num
    pack_num += 1
    return struct.pack("!BBHHH", 8, 0, checksum_calc(2049 + pack_num), 1, pack_num)

def echo(packet, ttl, dest):
    sock.setsockopt(socket.IPPROTO_IP, socket.IP_TTL, ttl)
    try:
        sock.sendto(packet, (socket.gethostbyname(dest), 0))
        start = time.time()
        data, address = sock.recvfrom(1024)
        rtt = "{:4.3f}".format((time.time() - start) * 1000)
        print(f"    {rtt}", end=" ")
        try:
            name = socket.gethostbyaddr(address[0])
            print(f"{name[0]} ({address[0]})")
        except:
            print(f"{address[0]}")
        return data[20] == 0 & data[21] == 0
    except socket.timeout:
        print(f"    *", end="")
        return False

def tracer(dest, packet_num):
    print(f"Destination: {dest}")
    reached = False
    for ttl in range(30):
        print(f"{ttl + 1}.")
        for _ in range(packet_num):
            reached |= echo(create_pack(), ttl + 1, dest)
        if reached:
            break
    sock.close()
    print("Reached destination")


if __name__ == "__main__":
    tracer(sys.argv[1], int(sys.argv[2]))