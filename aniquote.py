import click 
from download_quotes import main as update_func
from random_quote import main as rand_quote 


@click.command()
@click.option('--update', '-u', is_flag=True, help="update list of quotes")
@click.option('--random', '-r', is_flag=True,  help="display a random quote from the quotes folder")
def main(random, update):
    if random:
        rand_quote()
    elif update:
        update_func()
    else:
        click.echo("Please provide and argument, see --help")


if __name__ == "__main__":
    main()
