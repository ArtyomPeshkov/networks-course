import random

def crc(data):
    crc = 0xFF
    for bt in data:
        crc ^= bt
        for _ in range(8):
            crc = (crc << 1, (crc << 1) ^ 0x31)[crc & 0x80 == 128] & 255
    return crc

def corruptor(packet):
    sz = random.randint(0, len(packet) - 1)
    packet[sz] ^= 8
    return packet

f = open("data.txt", 'r')
data = bytearray(f.read().encode())
print(data)
print()
cnt = 0
while (len(data)>0):
    cnt += 1
    sz = min(5, len(data))
    packet = data[:sz]
    data = data[sz:]

    print("Data:               ", end="")
    print(packet)
    print("Code:               ", end="")
    print(crc(packet))
    packet.append(crc(packet))
    print("Encoded:            ", end="")
    print(packet)
    if (cnt % 3 == 0):
        corruptor(packet)
        print("Corrupted?: ", end="")
        print(packet)

    if crc(packet[:sz]) != packet[sz:][0]:
        print("Corrupted!")
    print()
