import requests
import json
import yaml
import pandas as pd
import getsteaminfo as steam

"""
    Downloads a user's Steam library and saves it to a CSV file.

    This function will look up the user's Steam ID and API key from the
    credentials.yaml file, download their Steam library, and save it to
    steam_library_<username>.csv.

    Currently, this script only works for the owner of the Steam API key.
    You can obtain a new API key at https://steamcommunity.com/dev/apikey.
"""


## Helper functions ##

def output_games_to_csv(games, personaname):
    df = pd.DataFrame(games)
    df.to_csv(f'steam_library_{personaname}.csv', index=False)


## Load credentials ##

with open('credentials.yaml') as f:
    credentials = yaml.safe_load(f)

api_key = credentials['api_key']
steam_id = credentials['steam_id']


## Main function ##

def main():
    output_games_to_csv(steam.get_steam_library(steam_id, api_key), steam.get_steam_personaname(steam_id, api_key))

if __name__ == '__main__':
    main()