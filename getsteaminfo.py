import requests
import yaml
import pandas as pd
import json

"""
Module for getting Steam information.

This module provides functions for getting a user's Steam library and
personaname.

"""

def get_steam_personaname(steam_id, api_key):
    url = 'http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key={}&steamids={}&format=json'.format(api_key, steam_id)

    r = requests.get(url)
    data = r.json()

    personaname = data['response']['players'][0]['personaname']
    return personaname

def get_steam_library(steam_id, api_key):
    url = 'http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key={}&steamid={}&format=json&include_appinfo=1'.format(api_key, steam_id)

    r = requests.get(url)
    data = r.json()

    games = data['response']['games']
    return games
