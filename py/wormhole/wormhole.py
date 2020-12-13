import socket
import threading
import sys
import argparse
import ipaddress

arguments   = None
source      = None
Target      = None

def getEndpoint(endpoint: str):
    parts = endpoint.split(':')
    if len(parts) != 2:
        raise Exception('Malformormed endpoint: ' + endpoint)
    # this statement raises an exception if the ip is malformed
    address = ipaddress.ip_address(parts[0])
    port = int(parts[1])
    if port < 0 or port > 65535:
        raise Exception('Invalid port number: ' + str(port))
    return (str(address), port)

def loop(source_socket, target_socket):
    while True:
        buffer = source_socket.recv(0x400)
        if len(buffer) == 0:
            print("No data received! Breaking...")
            break
        target_socket.send(buffer)

class Channel:

    def __init__(self, source_socket, destination_socket):
        self.source_socket          = source_socket
        self.destination_socket     = destination_socket
        self.source_thread          = threading.Thread(target = loop, args = (destination_socket, source_socket))
        self.destination_thread     = threading.Thread(target = loop, args = (source_socket, destination_socket))

    def start(self):
        self.source_thread.start()
        self.destination_thread.start()

    def close(self):
        self.source_socket.shutdown(socket.SHUT_RDWR)
        self.source_socket.close()
        self.destination_socket.shutdown(socket.SHUT_RDWR)
        self.destination_socket.close()

def open(arguments):
    source              = getEndpoint(arguments.source)
    max_connections     = 20
    target              = getEndpoint(arguments.target)
    server              = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    server.bind(source)
    server.listen(max_connections)
    connections = []
    print('Starting server...')
    try:
        while True:
            source_socket, _ = server.accept()
            print('Creating channel...')
            destination_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            destination_socket.connect(target)
            channel = Channel(source_socket, destination_socket)
            connections.append(channel)
            channel.start()
            print('Open channels: ' + str(len(connections)))
    except Exception as e:
        print(e)
        server.shutdown(socket.SHUT_RDWR)
        server.close()

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--source', dest = 'source', type = str, required = True, help = 'Source in form of <bind address>:<port>')
    parser.add_argument('-t', '--target', dest = 'target', type = str, help = 'Target in form of <address>:<port>')
    parser.add_argument('-m', '--maxconnections', dest = 'maxconnections', type = int, help = "Maximum number of concurrent connections")
    arguments = parser.parse_args()
    open(arguments)

if __name__ == "__main__":
    main()