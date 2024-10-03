import json 
import os 
import random



def main():
    directory = os.path.expanduser("~/.dotfiles/assets/quotes")
    list_quotes = []
    for file in os.listdir(directory):
        list_quotes.append(file)

    selection = random.choice(list_quotes)
    choice = f"~/.dotfiles/assets/quotes/{selection}"
    os.chdir(directory)
    with open(selection, 'r') as file:
        data = json.load(file)
        quote = data.get("quote")
        char = data.get("character")
        anime = data.get("anime")

    print(f'"{quote}"\n')
    print(f'\t{char} - "{anime}"')

if __name__ == "__main__":
    main()

