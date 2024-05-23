from scapy.all import ARP, Ether, srp
import socket

def scan(subnet):
    devices = []
    for _, received in srp(Ether(dst="ff:ff:ff:ff:ff:ff")/ARP(pdst=subnet), timeout=1)[0]:
        name = "*"
        try:
            name = socket.gethostbyaddr(received.psrc)[0]
        except:
            name = "*"
        devices.append({'ip': received.psrc, 'mac': received.hwsrc, 'name': name})
    return devices

if __name__ == "__main__":
    current = socket.gethostbyname(socket.gethostname()) 
    devices = scan(current + "/24")

    print("Devices in the network:")
    for device in devices:
        if (device['ip'] == current):
            print(f"IP: {device['ip']}, MAC: {device['mac']}, Name: {device['name']}")
    for device in devices:
        if (device['ip'] != current):
            print(f"IP: {device['ip']}, MAC: {device['mac']}, Name: {device['name']}")
