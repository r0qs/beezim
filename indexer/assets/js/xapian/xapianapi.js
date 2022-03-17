/* Copyright (C) 2022  Beezim

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>

This file is based on Runbox 7 (Runbox Solutions AS).
*/
Object.defineProperty(exports, "__esModule", { value: true });
exports.XapianAPI = void 0;
var emAllocateString = function (str) {
    if (!str) {
        str = "";
    }
    var $str = Module._malloc(str.length * 4 + 1);
    Module.stringToUTF8(str, $str, str.length * 4 + 1);
    return $str;
};
var XapianAPI = /** @class */ (function () {
    function XapianAPI() {
        this.initXapianIndex = Module.cwrap("initXapianIndex", null, ["string"]);
        this.initXapianIndexReadOnly = Module.cwrap("initXapianIndexReadOnly", null, ["string"]);
        this.addSingleFileXapianIndex = Module.cwrap("addSingleFileXapianIndex", null, ["string"]);
        this.compactDatabase = Module.cwrap("compactDatabase", null, []);
        this.compactToWritableDatabase = Module.cwrap("compactToWritableDatabase", null, ["string"]);
        this.addToXapianIndex = Module.cwrap("addToXapianIndex", null, ["string", "string"]);
        this.commitXapianUpdates = Module.cwrap("commitXapianUpdates", null, []);
        this.getXapianDocCount = Module.cwrap("getDocCount", "number", []);
        this.getLastDocid = Module.cwrap("getLastDocid", "number", []);
        this.reloadXapianDatabase = Module.cwrap("reloadDatabase", null, []);
        this.closeXapianDatabase = Module.cwrap("closeDatabase", null, []);
        this.setStringValueRange = Module.cwrap("setStringValueRange", null, ["number", "string"]);
        this.clearValueRange = Module.cwrap("clearValueRange", null, []);
        this.getNumericValue = Module.cwrap("getNumericValue", "number", ["number", "number"]);
        this.termlist = Module.cwrap("termlist", "number", ["string"]);
        this.documentTermList = Module.cwrap("documentTermList", "number", ["number"]);
        this.documentXTermList = Module.cwrap("documentXTermList", "number", ["number"]);
        this.deleteDocumentByUniqueTerm = Module.cwrap("deleteDocumentByUniqueTerm", null, ["string"]);
        this.deleteDocumentFromAddedWritablesByUniqueTerm = Module.cwrap("deleteDocumentFromAddedWritablesByUniqueTerm", "number", [
            "string",
        ]);
        this.setStringValue = Module.cwrap("setStringValue", null, ["number", "number", "string"]);
        this.addTermToDocument = Module.cwrap("addTermToDocument", null, ["string", "string"]);
        this.removeTermFromDocument = Module.cwrap("removeTermFromDocument", null, ["string", "string"]);
        this.addTextToDocument = Module.cwrap("addTextToDocument", null, [
            "string",
            "boolean",
            "string",
        ]);
        this.getDocIdFromUniqueIdTerm = Module.cwrap("getDocIdFromUniqueIdTerm", "number", ["string"]);
    }
    XapianAPI.prototype.getStringValue = function (docid, slot) {
        var $ret = Module._malloc(1024);
        Module._getStringValue(docid, slot, $ret);
        var ret = Module.UTF8ToString($ret);
        Module._free($ret);
        return ret;
    };
    XapianAPI.prototype.queryXapianIndex = function (querystring, offset, maxresults) {
        var $searchResults = Module._malloc(4 * maxresults);
        Module.HEAP8.set(new Uint8Array(maxresults * 4), $searchResults);
        var $queryString = emAllocateString(querystring);
        var $resultIdTerm = Module._malloc(128);
        var hits = Module._queryIndex($queryString, $searchResults, offset, maxresults);
        // console.log(hits);
        var results = new Array();
        for (var n = 0; n < hits; n++) {
            var docid = Module.getValue($searchResults + n * 4, "i32");
            Module._getDocumentData(docid, $resultIdTerm);
            results.push({
                docid: docid,
                data: Module.UTF8ToString($resultIdTerm),
            });
        }
        Module._free($searchResults);
        Module._free($queryString);
        Module._free($resultIdTerm);
        return results;
    };
    XapianAPI.prototype.getDocumentData = function (docid) {
        var $docdata = Module._malloc(1024);
        Module._getDocumentData(docid, $docdata);
        var ret = Module.UTF8ToString($docdata);
        Module._free($docdata);
        return ret;
    };
    return XapianAPI;
}());
exports.XapianAPI = XapianAPI;
//# sourceMappingURL=xapianapi.js.map