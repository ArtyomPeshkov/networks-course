from tkinter import *
from ftplib import FTP
from tkinter import messagebox

def ls(dir):
    if dir == "/":
        dir = "."
    files_list.insert(END, dir + "/")
    for filename, data in ftp.mlsd(dir):
        if data["type"] == "dir":
            ls(dir + "/" + filename)
        else:
            files_list.insert(END, dir + "/" + filename)

def connect_to_ftp():
    global ftp
    files_list.delete(0,END)

    ftp = FTP()
    ftp.connect(server_entry.get(), int(port_entry.get()))
    ftp.login(username_entry.get(), password_entry.get())

    ls(".")

def create_file():
    def save_file():
        filename = filename_entry.get()
        content = content_entry.get("1.0", END)
        with open(filename, 'wb') as file:
            file.write(content.encode())
        ftp.storbinary('STOR ' + filename, open(filename, 'rb'))


    create_window = Toplevel(root)
    create_window.title("Create File")

    filename_label = Label(create_window, text="Filename:")
    filename_label.pack()
    filename_entry = Entry(create_window)
    filename_entry.pack()

    content_label = Label(create_window, text="Content:")
    content_label.pack()
    content_entry = Text(create_window)
    content_entry.pack()

    save_button = Button(create_window, text="Save File", command=save_file)
    save_button.pack()


def retrieve_file():
    def get_data(data):
        global content_retr
        print(str(data.decode("utf-8")))
        content_retr = str(data.decode("utf-8"))

    def show_content():
        filename = filename_entry.get()

        ftp.retrbinary('RETR ' + filename, get_data)

        retrieve_window = Toplevel(root)
        retrieve_window.title("Retrieve File")

        content_label = Label(retrieve_window, text="Content:")
        content_label.pack()
        content_text = Text(retrieve_window)
        content_text.insert(END, content_retr)
        content_text.pack()

    retrieve_window = Toplevel(root)
    retrieve_window.title("Retrieve File")

    filename_label = Label(retrieve_window, text="Filename:")
    filename_label.pack()
    filename_entry = Entry(retrieve_window)
    filename_entry.pack()

    retrieve_button = Button(retrieve_window, text="Retrieve File", command=show_content)
    retrieve_button.pack()

def update_file():
    def get_data(data):
        global content_upd
        print(str(data.decode("utf-8")))
        content_upd = str(data.decode("utf-8"))
    def save_updated_file():
        content = content_entry.get("1.0", END)

        with open(filename, 'wb') as file:
            file.write(content.encode())

        ftp.storbinary('STOR ' + filename, open(filename, 'rb'))
    
    filename = filename_upd_entry.get()
    ftp.retrbinary('RETR ' + filename, get_data)
    update_window = Toplevel(root)
    update_window.title("Update File")

    content_label = Label(update_window, text="Content:")
    content_label.pack()
    content_entry = Text(update_window)
    content_entry.insert(END, content_upd)
    content_entry.pack()

    save_button = Button(update_window, text="Save File", command=save_updated_file)
    save_button.pack()

def delete_file():
    def delete():
        filename = filename_entry.get()
        ftp.delete(filename)

    delete_window = Toplevel(root)
    delete_window.title("Delete File")

    filename_label = Label(delete_window, text="Filename:")
    filename_label.pack()
    filename_entry = Entry(delete_window)
    filename_entry.pack()

    delete_button = Button(delete_window, text="Delete File", command=delete)
    delete_button.pack()

root = Tk()
root.title("FTP Client")

server_label = Label(root, text="Server:")
server_label.pack()
server_entry = Entry(root)
server_entry.insert(END, "127.0.0.1")
server_entry.pack()

port_label = Label(root, text="Port:")
port_label.pack()
port_entry = Entry(root)
port_entry.insert(END, "21")
port_entry.pack()

username_label = Label(root, text="Username:")
username_label.pack()
username_entry = Entry(root)
username_entry.insert(END, "TestUser")
username_entry.pack()

password_label = Label(root, text="Password:")
password_label.pack()
password_entry = Entry(root, show="*")
password_entry.pack()

connect_button = Button(root, text="Connect", command=connect_to_ftp)
connect_button.pack()

files_list = Listbox(root, width=50)
files_list.pack()

create_button = Button(root, text="Create", command=create_file)
create_button.pack()

retrieve_button = Button(root, text="Retrieve", command=retrieve_file)
retrieve_button.pack()

filename_upd_entry = Entry(root)
filename_upd_entry.pack()
update_button = Button(root, text="Update", command=update_file)
update_button.pack()

delete_button = Button(root, text="Delete", command=delete_file)
delete_button.pack()

root.mainloop()