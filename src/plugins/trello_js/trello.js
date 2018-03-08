// Get Trello Details
var Trello = require("node-trello");

var client = new Trello(key, token);

var tasks = process.argv[2];
var jsonTasks = JSON.parse(tasks);

// Loop over passed task and add them to trello list
jsonTasks.forEach(function(item){
	var card = {
		idList:"5aa0057dae6f639766e9bff4",
		name: item.TaskName,
		desc: "File: "+item.FileName + "\nLine: "+ item.FileLine+"\nCODE GOES HERE",
	}
	client.post("/1/cards?idList="+card.idList+"&name="+card.name+"&desc="+card.desc, function(err, data){
		if (err) throw err;
		console.log(data);
	});
});


