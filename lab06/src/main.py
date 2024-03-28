from ftplib import FTP

ftp = FTP("127.0.0.1")
ftp.login("TestUser", "")

def ls(dir):
    print(dir + "/")
    for filename, data in ftp.mlsd(dir):
        if data["type"] == "dir":
            ls(dir + "/" + filename)
        else:
            print(dir + "/" + filename)

def store(client_path, server_path):
    file = open(client_path, "rb")
    ftp.storbinary("STOR " + server_path, file)

def load(server_path, client_path):
    file = open(client_path, "wb")
    ftp.retrbinary("RETR " + server_path, file.write)

while True:
    inp = input().split(" ")
    if inp[0] == "ls":
        ls(".")
    elif inp[0] == "store":
        store(inp[1], inp[2])
    elif inp[0] == "load":
        load(inp[1], inp[2])
    elif inp[0] == "exit":
        break
    else:
        print("Undefined command")

ftp.quit()