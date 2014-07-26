var cheerio = require("cheerio");
var request = require('request');
var _ = require("underscore");
var http = require('http');
var fs = require('fs');

console.log("running ...");
request("http://somesite.net/", function (error, response, body) {
	if (!error && response.statusCode == 200) {

		console.log("received response");
		var $ = cheerio.load(body);

		var links = [];
		$("div.post-content a[title]").each(function(i, e) {
			links.push($(e).attr("href"));
		});

		var links = _.unique(links);

		console.log("received " + links.length + " links");

		_.each(links, function(link){
			request(link, function (error, response, body) {
				if (!error && response.statusCode == 200) {
					var results = body.match('file":"http:.*mp4');
					if (results) {
						var videolink = results[0].replace('file":"',"");
						
						var crypto = require('crypto');
						var hash = crypto.createHash('md5').update(videolink).digest('hex');
						var savingFileName = hash + ".mp4";

						if (!fs.existsSync(savingFileName)) {
							var file = fs.createWriteStream(savingFileName);
							http.get(videolink, function(response) {
					  			response.pipe(file);
							});

							console.log("saving " + savingFileName + " from " + videolink);
						}
						else {
							console.log(savingFileName + " already exist.");
						}
					}
				}
			});
		});
	}
});
