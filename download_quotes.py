from re import I
import requests
import psutil
import socket 
import random 
import json
import os 
from datetime import datetime
from tqdm import tqdm
import time


def is_network_connected():
    addrs = psutil.net_if_addrs()
    stats = psutil.net_if_stats()

    # Check if there's a network interface that is "up" and has an IP address
    for interface, info in stats.items():
        if info.isup and interface != 'lo':  # Skip loopback interface ('lo')
            # Check if the interface has an IPv4 or IPv6 address
            if interface in addrs:
                for addr in addrs[interface]:
                    if addr.family == socket.AF_INET or addr.family == socket.AF_INET6:
                        return True
    return False


def save_quote_to_new_file(quote_data, folder=os.path.expanduser('~/.dotfiles/assets/quotes')):
    # Create folder if it doesn't exist
    if not os.path.exists(folder):
        os.makedirs(folder)

    # Create a unique filename using the current timestamp
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    filename = f"{folder}/quote_{timestamp}.json"

    # Save the quote to a new JSON file
    with open(filename, 'w') as file:
        json.dump(quote_data, file, indent=4)



def download():
    if not is_network_connected():
        print("no wifi")
    else:
        try:
            response = requests.get('https://animechan.io/api/v1/quotes/random')
            data = response.json()

            if "message" in data:
                print("API busy")
            else:
                quote = data["data"].get("content")
                character = data["data"]["character"].get("name")
                anime = data["data"]["anime"].get("name")

                json_quote = {
                        "quote": quote,
                        "character": character,
                        "anime": anime
                        }
                save_quote_to_new_file(json_quote)
        except requests.exceptions.RequestException as e:
            print(f"An error occured: {e}")
        time.sleep(2)


def main():
    for i in tqdm(range(19), desc="downloading quotes"):
        download()


if __name__ == "__main__":
    main()
