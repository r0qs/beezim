// TODO: Remove mock:
var XapianAPI = function() {
    this.initXapianIndexReadOnly = function(){};
    this.bSearch = [{"query":"b","title":"ATC code B05BA01","data":"A/ATC_code_B05BA01"},{"query":"b","title":"ATCvet code QB05BA01","data":"A/ATCvet_code_QB05BA01"},{"query":"b","title":"2002 Berlin controversy involving Michael Jackson","data":"A/2002_Berlin_controversy_involving_Michael_Jackson"},{"docid":27,"data":"A/Group_(mathematics)","wordcount":13943,"title":"Group (mathematics)"},{"docid":68,"data":"A/Giraffe","wordcount":12386,"title":"Giraffe"},{"docid":63,"data":"A/Bat","wordcount":16295,"title":"Bat"},{"docid":72,"data":"A/Hippopotamus","wordcount":8965,"title":"Hippopotamus"},{"docid":11,"data":"A/Virus","wordcount":17488,"title":"Virus"},{"docid":23,"data":"A/Lion","wordcount":16522,"title":"Lion"},{"docid":50,"data":"A/Termite","wordcount":15916,"title":"Termite"},{"docid":80,"data":"A/Tiger","wordcount":16018,"title":"Tiger"},{"docid":92,"data":"A/Frog","wordcount":18776,"title":"Frog"},{"docid":83,"data":"A/Raccoon","wordcount":14783,"title":"Raccoon"},{"docid":18,"data":"A/India","wordcount":25803,"title":"India"},{"docid":58,"data":"A/Cretaceous–Paleogene_extinction_event","wordcount":13590,"title":"Cretaceous–Paleogene extinction event"},{"docid":98,"data":"A/Jesus","wordcount":29041,"title":"Jesus"},{"docid":35,"data":"A/Cougar","wordcount":12939,"title":"Cougar"},{"docid":48,"data":"A/Mollusca","wordcount":9011,"title":"Mollusca"},{"docid":4,"data":"A/Beetle","wordcount":16326,"title":"Beetle"},{"docid":97,"data":"A/Arthropod","wordcount":10502,"title":"Arthropod"}];
    this.beSearch = [{"query":"be","title":"Addis Abeba lion","data":"A/Addis_Abeba_lion"},{"query":"be","title":"Accession of Queen Elizabeth II","data":"A/Accession_of_Queen_Elizabeth_II"},{"query":"be","title":"2002 Berlin controversy involving Michael Jackson","data":"A/2002_Berlin_controversy_involving_Michael_Jackson"},{"docid":90,"data":"A/Solar_eclipse","wordcount":10121,"title":"Solar eclipse"},{"docid":17,"data":"A/Operating_system","wordcount":10900,"title":"Operating system"},{"docid":16,"data":"A/Binomial_nomenclature","wordcount":6078,"title":"Binomial nomenclature"},{"docid":55,"data":"A/Myocardial_infarction","wordcount":12629,"title":"Myocardial infarction"},{"docid":57,"data":"A/Bivalvia","wordcount":13483,"title":"Bivalvia"},{"docid":27,"data":"A/Group_(mathematics)","wordcount":13943,"title":"Group (mathematics)"},{"docid":87,"data":"A/Video_game","wordcount":15063,"title":"Video game"},{"docid":92,"data":"A/Frog","wordcount":18776,"title":"Frog"},{"docid":94,"data":"A/The_Beatles","wordcount":29968,"title":"The Beatles"},{"docid":98,"data":"A/Jesus","wordcount":29041,"title":"Jesus"},{"docid":7,"data":"A/Charles_Darwin","wordcount":16832,"title":"Charles Darwin"},{"docid":71,"data":"A/Elvis_Presley","wordcount":28293,"title":"Elvis Presley"},{"docid":78,"data":"A/Catholic_Church","wordcount":34002,"title":"Catholic Church"},{"docid":11,"data":"A/Virus","wordcount":17488,"title":"Virus"},{"docid":76,"data":"A/United_States_Declaration_of_Independence","wordcount":18878,"title":"United States Declaration of Independence"},{"docid":44,"data":"A/Apollo_13","wordcount":14343,"title":"Apollo 13"},{"docid":46,"data":"A/Rodent","wordcount":13283,"title":"Rodent"}];
    this.beaSearch = [{"query":"bea","title":"Beatles, The","data":"A/Beatles,_The"},{"query":"bea","title":"Beatles, the","data":"A/Beatles,_the"},{"query":"bea","title":"Beatles (The)","data":"A/Beatles_(The)"},{"docid":66,"data":"A/Washington,_D.C.","wordcount":25007,"title":"Washington, D.C."},{"docid":94,"data":"A/The_Beatles","wordcount":29968,"title":"The Beatles"},{"docid":7,"data":"A/Charles_Darwin","wordcount":16832,"title":"Charles Darwin"},{"docid":46,"data":"A/Rodent","wordcount":13283,"title":"Rodent"},{"docid":35,"data":"A/Cougar","wordcount":12939,"title":"Cougar"},{"docid":83,"data":"A/Raccoon","wordcount":14783,"title":"Raccoon"},{"docid":52,"data":"A/Michael_Jackson","wordcount":29848,"title":"Michael Jackson"},{"docid":71,"data":"A/Elvis_Presley","wordcount":28293,"title":"Elvis Presley"},{"docid":75,"data":"A/Bobcat","wordcount":8997,"title":"Bobcat"},{"docid":36,"data":"A/Bob_Dylan","wordcount":33345,"title":"Bob Dylan"},{"docid":45,"data":"A/Elizabeth_II","wordcount":16678,"title":"Elizabeth II"},{"docid":13,"data":"A/William_Shakespeare","wordcount":16744,"title":"William Shakespeare"},{"docid":80,"data":"A/Tiger","wordcount":16018,"title":"Tiger"},{"docid":5,"data":"A/List_of_United_States_cities_by_population","wordcount":8764,"title":"List of United States cities by population"},{"docid":30,"data":"A/List_of_national_parks_of_the_United_States","wordcount":7595,"title":"List of national parks of the United States"},{"docid":72,"data":"A/Hippopotamus","wordcount":8965,"title":"Hippopotamus"},{"docid":42,"data":"A/Turtle","wordcount":12137,"title":"Turtle"}];
    this.beatSearch = [{"query":"beat","title":"Beatles, The","data":"A/Beatles,_The"},{"query":"beat","title":"Beatles, the","data":"A/Beatles,_the"},{"query":"beat","title":"Beatles (The)","data":"A/Beatles_(The)"},{"docid":52,"data":"A/Michael_Jackson","wordcount":29848,"title":"Michael Jackson"},{"docid":94,"data":"A/The_Beatles","wordcount":29968,"title":"The Beatles"},{"docid":71,"data":"A/Elvis_Presley","wordcount":28293,"title":"Elvis Presley"},{"docid":39,"data":"A/Midfielder","wordcount":5676,"title":"Midfielder"},{"docid":36,"data":"A/Bob_Dylan","wordcount":33345,"title":"Bob Dylan"},{"docid":87,"data":"A/Video_game","wordcount":15063,"title":"Video game"},{"docid":63,"data":"A/Bat","wordcount":16295,"title":"Bat"},{"docid":45,"data":"A/Elizabeth_II","wordcount":16678,"title":"Elizabeth II"},{"docid":9,"data":"A/Anfield","wordcount":7609,"title":"Anfield"},{"docid":48,"data":"A/Mollusca","wordcount":9011,"title":"Mollusca"},{"docid":64,"data":"A/Human","wordcount":21211,"title":"Human"},{"docid":55,"data":"A/Myocardial_infarction","wordcount":12629,"title":"Myocardial infarction"},{"docid":84,"data":"A/Argentina","wordcount":24748,"title":"Argentina"},{"docid":22,"data":"A/Billboard_(magazine)","wordcount":4666,"title":"Billboard (magazine)"},{"docid":69,"data":"A/Kit_(association_football)","wordcount":6458,"title":"Kit (association football)"},{"docid":68,"data":"A/Giraffe","wordcount":12386,"title":"Giraffe"},{"docid":14,"data":"A/Thoroughbred","wordcount":8370,"title":"Thoroughbred"}];
    this.elviSearch = [{"query":"elvi","title":"Elvis Presly","data":"A/Elvis_Presly"},{"query":"elvi","title":"Elvis Pres;ey","data":"A/Elvis_Pres;ey"},{"query":"elvi","title":"Elvis Presely","data":"A/Elvis_Presely"},{"docid":71,"data":"A/Elvis_Presley","wordcount":28293,"title":"Elvis Presley"},{"docid":94,"data":"A/The_Beatles","wordcount":29968,"title":"The Beatles"},{"docid":36,"data":"A/Bob_Dylan","wordcount":33345,"title":"Bob Dylan"},{"docid":52,"data":"A/Michael_Jackson","wordcount":29848,"title":"Michael Jackson"},{"docid":86,"data":"A/United_States","wordcount":29078,"title":"United States"},{"docid":77,"data":"A/London","wordcount":26097,"title":"London"},{"docid":98,"data":"A/Jesus","wordcount":29041,"title":"Jesus"},{"query":"elvi","title":"Elvis Presley","data":"A/Elvis_Presley"},{"query":"elvi","title":"Elvis Pressley","data":"A/Elvis_Pressley"},{"query":"elvi","title":"Elvis A. Presley","data":"A/Elvis_A._Presley"},{"query":"elvi","title":"Elvis Aron Presley","data":"A/Elvis_Aron_Presley"},{"query":"elvi","title":"Elvis Extravaganza","data":"A/Elvis_Extravaganza"},{"query":"elvi","title":"Elvis Aaron Presley","data":"A/Elvis_Aaron_Presley"},{"query":"elvi","title":"Death of Elvis Presley","data":"A/Death_of_Elvis_Presley"},{"query":"elvi","title":"Elvis Presley Discography","data":"A/Elvis_Presley_Discography"},{"query":"elvi","title":"Elvis Presley discography","data":"A/Elvis_Presley_discography"},{"query":"elvi","title":"Elvis Presley's political beliefs","data":"A/Elvis_Presley's_political_beliefs"}];
}

addOnPostRun = function(fnc){fnc.apply();};

var BeeZIMSearcher = function() {
	this.fulltextIndex = "./X/fulltext/xapian";
	this.titleIndex = "./X/title/xapian"; // TODO: Can not load both indexes at the same time due to variables are overwritten in the global namespace.
	this.dataDir = "./X";
	this.articles = [];
	this.initRan = false;
	this.xapian = new XapianAPI();
	this.maxResults = 20;
	this.titleResults = 3;

	this.init(this.fulltextIndex);
}

BeeZIMSearcher.prototype.init = function(index) {
	var self = this;
	addOnPostRun(function(){
		self.xapian.initXapianIndexReadOnly(self.fulltextIndex);
		self.initRan = true;
	});
	var xhr = new XMLHttpRequest();
	xhr.open("GET","files.json");
	xhr.onload = function() {
		self.parseFiles(xhr.response);
	};
	xhr.send();
}

BeeZIMSearcher.prototype.parseFiles = function(filesResponse) {
	let files = JSON.parse(filesResponse);
	for (const [key, value] of Object.entries(files)){
		if (key.startsWith("A/")){
			this.articles.push(value);
		}
	}
}

BeeZIMSearcher.prototype.GetRandomArticle = function() {
	if (!this.initRan) {
		return "You need to run 'init()' before searching!";
	}
	return this.articles[this.articles.length * Math.random() << 0];
}

BeeZIMSearcher.prototype.Search = function(query) {
	if (!this.initRan) {
		return "You need to run 'init()' before searching!";
	}
	
	let results = [];
	let queryLower = query.toLowerCase();

	this.xapian.queryXapianIndex(query,0,this.maxResults-this.titleResults).forEach((r) => {
		results.push({
			docid: r.docid,
			data: r.data,
			wordcount: parseInt(this.xapian.getStringValue(r.docid,1)),
			title: this.xapian.getStringValue(r.docid,0)
		});
	});
	let wantedTitleMatch = this.maxResults-results.length;
	let titleResults = [];
	for (let i = 0; i < this.articles.length; i++) {
		if (wantedTitleMatch <= 0)
			break;
		let value = this.articles[i];
		if(value.Metadata.Title.toLowerCase().indexOf(queryLower) > -1) {
			titleResults.push({
				query: query,
				title: value.Metadata.Title,
				data: value.Path
			});
			wantedTitleMatch--;
		}
	}
	titleResults.sort(function(a,b) {
		return a.title.length - b.title.length;
	});

	// Top 3 results are from titleResults
	if (titleResults.length > this.titleResults){
		return titleResults.slice(0,this.titleResults).concat(results,titleResults.slice(this.titleResults))
	}

	return titleResults.concat(results);
}

BeeZIMSearcher.prototype.Search = function(query) {
    let x = new XapianAPI();
    if (query == "b")
        return x.bSearch;
    if (query == "be")
    return x.beSearch;
    if (query == "bea")
    return x.beaSearch;
    if (query == "beat")
    return x.beatSearch;

    return x.elviSearch;
}