import json

base_url = "https://api.trello.com/1/"
conf = {'key': "", 'token': "", 'remote_link': "",
        'backlog_id': "", 'done_id': ""}


def get_config():
    with open('config.json') as json_data:
        d = json.load(json_data)
        return d