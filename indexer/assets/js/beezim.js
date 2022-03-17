function asyncFetch(method, url, opts = {}) {
	return new Promise(function (resolve, reject) {
		let xhr = new XMLHttpRequest();
		if (opts.responseType) {
			xhr.responseType = opts.responseType;
		}
		if (opts.mimeType) {
			xhr.overrideMimeType(opts.mimeType);
		}
		xhr.open(method, url);
		xhr.onload = function () {
			if (this.status < 400) {
				resolve(xhr.response);
			} else {
				reject({
					status: this.status,
					statusText: xhr.statusText
				});
			}
		};
		xhr.onerror = function () {
			reject({
				status: this.status,
				statusText: xhr.statusText
			});
		};
		xhr.send();
	});
}

class BeeZIMSearcher {
	#articles = [];
	#initRan = false;
	#indexURL;
	#xapian;
	static #beeZim;

	constructor(indexURL, xapianPath) {
		this.#indexURL = indexURL;
		this.#xapian = new XapianAPI();
		this.#xapian.initXapianIndexReadOnly(xapianPath);
		this.#initRan = true;
	}

	static Init(indexURL) {
		// Note: /data is created and mounted on the pre.js included in the compiled code.
		const xapianIDBFSPath = "/data/xapian";
		return new Promise(async function (resolve, reject) {
			if (BeeZIMSearcher.#beeZim)
				resolve(BeeZIMSearcher.#beeZim);
			try {
				const opts = {
					mimeType: "application/octet-stream+xapian",
					responseType: "blob"
				}
				const response = await asyncFetch("GET", indexURL, opts);
				if (!response) {
					throw ("error retrieving index DB");
				}

				// Convert blob to Uint8Array and write the index DB
				const data = new Uint8Array(await response.arrayBuffer());
				const stream = FS.open(xapianIDBFSPath, 'w+');
				FS.write(stream, data, 0, data.length, 0);
				FS.close(stream);

				// sync from MEMFS to IDBFS
				FS.syncfs(false, function (err) {
					if (err) {
						throw err;
					}
					BeeZIMSearcher.#beeZim = new BeeZIMSearcher(indexURL, xapianIDBFSPath);
					return resolve(BeeZIMSearcher.#beeZim);
				});
			} catch (err) {
				reject(err);
			}
		}).catch(function (err) {
			console.error(err);
		});
	}

	async LoadFiles() {
		if (this.#articles.length == 0) {
			const files = await asyncFetch("GET", "files.json")
			this.#parseFiles(files);
		}
	}

	#parseFiles(filesResponse) {
		let files = JSON.parse(filesResponse);
		for (const [key, value] of Object.entries(files)) {
			if (key.startsWith("A/")) {
				this.#articles.push(value);
			}
		}
	}

	GetRandomArticle() {
		if (!this.#initRan) {
			return "You need to run 'Init()' before searching!";
		}
		return this.#articles[this.#articles.length * Math.random() << 0];
	}

	IndexSearch(query, offset=0, maxResults=1000) {
		if (!query) {
			return [];
		}

		if (!this.#initRan) {
			return "You need to run 'Init()' before searching!";
		}

		let results = [];

		this.#xapian.queryXapianIndex(query, offset, maxResults).forEach((r) => {
			results.push({
				docid: r.docid,
				data: r.data,
				wordcount: parseInt(this.#xapian.getStringValue(r.docid, 1)),
				title: this.#xapian.getStringValue(r.docid, 0)
			});
		});

		return results;
	}

	QuickSearch(query, 	maxResults = 20, titleMatches = 3) {
		if (!query) {
			return [];
		}

		if (!this.#initRan) {
			return "You need to run 'Init()' before searching!";
		}

		let results = [];
		let queryLower = query.toLowerCase();

		this.#xapian.queryXapianIndex(query, 0, maxResults-titleMatches).forEach((r) => {
			results.push({
				docid: r.docid,
				data: r.data,
				wordcount: parseInt(this.#xapian.getStringValue(r.docid, 1)),
				title: this.#xapian.getStringValue(r.docid, 0)
			});
		});

		let wantedTitleMatch = maxResults - results.length;
		let titleResults = [];
		for (let i = 0; i < this.#articles.length; i++) {
			if (wantedTitleMatch <= 0)
				break;
			let value = this.#articles[i];
			if (value.Metadata.Title.toLowerCase().indexOf(queryLower) > -1) {
				titleResults.push({
					query: query,
					title: value.Metadata.Title,
					data: value.Path
				});
				wantedTitleMatch--;
			}
		}
		titleResults.sort(function (a, b) {
			return a.title.length - b.title.length;
		});

		// Top (titleMatches) results are from titleResults
		if (titleResults.length > titleMatches) {
			return titleResults.slice(0, titleMatches).concat(results, titleResults.slice(titleMatches))
		}

		return titleResults.concat(results);
	}

	async GetTextContent(url) {
		const htmlContent = await asyncFetch("GET", url);
		let tmp = document.createElement("DIV");
		tmp.innerHTML = htmlContent;
		let content = tmp.querySelectorAll("#content p");
		let str = "";
		content.forEach((c) => {
			let cStyle = c.querySelector("style");
			let out = c.innerText || c.textContent || "";
			if (cStyle != null)
				out = out.replace(cStyle.innerText,"");
			str += out;
		});
		return str;
	}
}
