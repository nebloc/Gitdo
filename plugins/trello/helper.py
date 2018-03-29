import json

base_url = "https://api.trello.com/1/"

def get_config():
    conf = {'key':"", 'token':"", 'remote_link':"", 'backlog_id':"", 'done_id': ""}
    with open('config.json') as json_data:
        d = json.load(json_data)
        for attr in d:
            if attr == "":
                return
        return d