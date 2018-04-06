/*----------------------------------- core -----------------------------------*/
(function (w, doc, module) {
	var noop = function () { },
		undef = undefined;

	function elem(tag) {
		var o = doc.createElement(tag),
			props = arguments.length > 1 && arguments[1] ? arguments[1] : undef,
			attrs = arguments.length > 2 && arguments[2] ? arguments[2] : undef;

		if (attrs) {
			for (var a in attrs) {
				o.setAttribute(a, attrs[a]);
			}
		}

		if (props) {
			for (var p in props) {
				if (p == 'children') {
					if (typeof props.children[Symbol.iterator] === 'function' ) {
						props.children.forEach(function (c) { o.appendChild(c) });
					} else {
						o.appendChild(props.children);
					}
				} else {
					o[p] = props[p];
				}
			}
		}
		return o;
	}

	function clear() {
		Array.prototype.forEach.call(arguments, function (elem) {
			while (elem.firstChild) elem.removeChild(elem.firstChild);
		});
	}

	function setDefaults(x, y) {
		for (var o in y) {
			if (!(o in x)) {
				x[o] = y[o];
			}
		}
	}

	w[module] = {
		elem: elem,
		clear: clear,
		setDefaults: setDefaults
	};

})(window, window.document, '$');

/*----------------------------------- AJAX -----------------------------------*/
(function (module) {
	var defaults = {
			url: '/',
			method: 'GET',
			accept: 'application/json',
			contentType: 'text/plain',
			onWorking: $.noop,
			onSuccess: $.noop,
			onError: $.noop
		};

	function ajax(args) {
		var r = new XMLHttpRequest();//,
			t = setTimeout(args.onWorking, 250);

		$.setDefaults(args, defaults);

		r.open(args.method, args.url, true);
		r.setRequestHeader('Accept', args.accept);
		r.setRequestHeader('Content-Type', args.contentType);
		r.onreadystatechange = function () {
			if (r.readyState != 4) return;
			var o = args.accept == 'application/json' ? JSON.parse(r.responseText) : responseText;

			clearTimeout(t);
			if (r.status != 200) {
				args.onError(o);
				return;
			}
			args.onSuccess(o);
		}
		r.send(args.data);
	}

	$.ajax = ajax;

})($);

/*--------------------------------- overlay ----------------------------------*/
(function (doc, $) {

	function overlay(selector) {
		var elem = doc.querySelector(selector);
		return {
			show: function () { elem.style.display = '' },
			hide: function () { elem.style.display = 'none' }
		}
	}

	$.overlay = overlay;

})(window.document, $);

/*---------------------------------- iquery ----------------------------------*/
(function ($) {
	$.iquery = { };
})(window['$']);

/*-------------- query ---------------*/
(function (w, doc, $) {

	function clearTable(table) {
		$.clear(table.tHead);
		$.clear(table.tBodies.item(0));
	}

	function tableRow(row, values, cellType) {
		for (var i in values) {
			var cell = doc.createElement(cellType)
			if (values[i] == null) {
				cell.className = 'nil';
			} else {
				cell.textContent = values[i];
			}
			row.appendChild(cell);
		}
	}

	function noRows(table, count) {
		var row = table.tBodies.item(0).insertRow(),
			cell = doc.createElement('td');

		cell.className = 'noResults';
		cell.colSpan = count;
		row.appendChild(cell);
	}

	function header(table, columns) {
		var row = table.tHead.insertRow();
		tableRow(row, columns, 'th')
	}

	function rows(table, rows, count) {
		rows.forEach(function (values) {
			var row = table.tBodies.item(0).insertRow();
			tableRow(row, values, 'td');
		});
	}

	function querySet(result) {
		var table = $.elem('table', { children: [$.elem('thead'), $.elem('tbody')] });

		header(table, result.columns);
		if (result.rows) {
			rows(table, result.rows);
		} else {
			noRows(table, result.columns ? result.columns.length : 0);
		}

		return table;
	}

	function setQuerySetVisible(i, tabs, querySets, visible) {
		if (visible) {
			tabs.childNodes[i].classList.add('selected')
			querySets.childNodes[i].classList.add('selected')
		} else {
			tabs.childNodes[i].classList.remove('selected')
			querySets.childNodes[i].classList.remove('selected')
		}
	}

	function querySets(results) {
		var tabs = $.elem('div', { className: 'querySetTabs' }),
			querySets = $.elem('div', { className: 'querySetResults' }),
			container = $.elem('div', { className: 'querySets', children: [tabs, querySets] });

		results.forEach(function (result) {
			var qs = querySet(result),
				tab = $.elem('a', { href: 'javascript: void(0)', textContent: querySets.childNodes.length+1 });

			tab.addEventListener('click', function () {
				var t = Array.prototype.indexOf.call(tabs.childNodes, this);
				for (var i=0; i<tabs.childNodes.length; i++) {
					setQuerySetVisible(i, tabs, querySets, t == i);
				}
			});
			tabs.appendChild(tab);
			querySets.appendChild(qs);
		});

		if (results.length > 0) {
			setQuerySetVisible(0, tabs, querySets, true);
		}

		return container;
	}

	function error(err) {
		if (err && err.error) {
			alert(err.error);
		} else {
			alert(err);
		}
	}

	/*------------------------------------------------------------------------*/

	function init(input, results, database,  placeholder, theme, overlay) {
		var editor = null,
			limit = 1000;

		function populateResults(data) {
			var qs = Array.isArray(data)
				? querySets(data)
				: querySet(data);

			$.clear(results);
			results.appendChild(qs);
		}

		function getQuery() {
			var rows = editor.getValue().split('\n'),
				s = editor.getSelection(),
				sql;

			if (s.startLineNumber == s.endLineNumber) {
				if (s.startColumn == s.endColumn) {
					// No selection, execute the full value.
					return rows.join('\n');
				}
				// Single line selected.
				return rows[s.startLineNumber-1].substring(s.startColumn-1, s.endColumn-1);
			}

			sql = [rows[s.startLineNumber-1].substring(s.startColumn-1)];
			for (var i=s.startLineNumber; i<s.endLineNumber-1; i++) {
				sql.push(rows[i]);
			}
			sql.push(rows[s.endLineNumber-1].substring(0, s.endColumn-1));
			return sql.join('\n');
		}

		function query() {
			$.ajax({
				url: ['/databases/', database, '?action=query&limit=', limit].join(''),
				method: 'POST',
				data: getQuery(),
				onSuccess: function (results) {
					overlay.hide();
					populateResults(results)
				},
				onError: function (e) {
					overlay.hide();
					error(e);
				},
				onWorking: overlay.show
			});
		}

		require(['vs/editor/editor.main'], function() {
			editor = monaco.editor.create(input, {
				value: placeholder,
				theme: theme,
				language: 'sql'
			});

			input.addEventListener('keydown', function (event) {
				if ((event.metaKey || event.ctrlKey) && event.keyCode == 13) {
					event.preventDefault();
					event.stopPropagation();

					query();
				}
			}, true);

			w.addEventListener('resize', function () { editor.layout() }, false);
		});


		return {
			getQuery: function () { return editor.getValue() },
			setQuery: function (value) { return editor.setValue(value) },
			setDatabase: function (value) { database = value },
			setLimit: function (value) { limit = value },
			setTheme: function (value) { editor.updateOptions({ 'theme': value }) },
			layout: function () { editor.layout() }
		};
	}

	require.config({ paths: { 'vs': '/static/monaco-editor/min/vs' }});
	$.iquery.query = init;

})(window, window.document, window['$']);

/*------------ navigator -------------*/
(function (w, doc, $) {

	function addSection(elem, name) {
		var ul = $.elem('ul'),
			o = $.elem('div', {
				className: name.toLowerCase(),
				children: [$.elem('label', { textContent: name }), ul]
			});

		elem.appendChild(o);
		return ul;
	}

	function noResults(section) {
		section.appendChild($.elem('li', { className: 'noResults' }));
	}

	function populateSchemas(elem, schemas, onSelect) {
		if (!schemas) noResults(elem);

		schemas.forEach(function (schema) {
			var li = doc.createElement('li');
			li.textContent = schema;
			li.addEventListener('click', function () {
				if (this.classList.contains('selected')) return;

				var current = this.parentNode.querySelector('.selected');
				if (current) current.classList.remove('selected');
				this.classList.add('selected');
				onSelect(schema);
			}, false);
			elem.appendChild(li);
		});
	}

	function populateTables(elem, tables, onSelect) {
		if (!tables) noResults(elem);

		tables.forEach(function (table) {
			var li,
				children = [ $.elem('span', { className: 'name', textContent: table.name }) ];

			if (table.description) {
				children.push( $.elem('span', { className: 'description', textContent: table.description }) );
			}
			li = $.elem('li', { children: children });
			li.addEventListener('click', function () {
				if (this.classList.contains('selected')) return;

				var current = this.parentNode.querySelector('.selected');
				if (current) current.classList.remove('selected');
				this.classList.add('selected');
				onSelect(table.name)
			}, false);
			elem.appendChild(li);
		});
	}

	function populateColumns(elem, columns) {
		if (!columns) noResults(elem);

		columns.forEach(function (column) {
			var children = [
					$.elem('span', { className: 'name', textContent: column.name }),
					$.elem('span', { className: 'type', textContent: column.type }),
					$.elem('span', { className: 'nullable', textContent: column.isNullable }),
					$.elem('span', { className: 'default', textContent: column.default })
				];

			if (column.description) {
				children.splice(1, 0, $.elem('span', { className: 'description', textContent: column.description }));
			}
			elem.appendChild($.elem('li', { children: children }));
		});
	}

	function error(err) {
		if (err && err.error) {
			alert(err.error);
		} else {
			alert(err);
		}
	}

	/*------------------------------------------------------------------------*/

	function init(elem, database, overlay) {
		var schemas = addSection(elem, 'Schemas'),
			tables = addSection(elem, 'Tables'),
			columns = addSection(elem, 'Columns');

		function loadItems(url, section, clearSections, onPopulate, onSelection) {
			$.clear.apply($, clearSections);
			$.ajax({
				url: url,
				method: 'GET',
				onSuccess: function (results) {
					overlay.hide();
					onPopulate(section, results, onSelection);
				},
				onError: function (e) {
					overlay.hide();
					error(e);
				},
				onWorking: overlay.show
			});
		}

		function loadSchemas(db) {
			var url = ['/databases/', db, '/schemas/'].join('');
			database = db;
			loadItems(url, schemas, [schemas, tables, columns], populateSchemas, loadTables);
		}

		function loadTables(schema) {
			var url = ['/databases/', database, '/schemas/', schema, '/tables/'].join('');
			loadItems(url, tables, [tables, columns], populateTables,
				function (table) {
					loadColumns(schema, table);
				});
		}

		function loadColumns(schema, table) {
			var url = ['/databases/', database, '/schemas/', schema, '/tables/', table].join('');
			loadItems(url, columns, [columns], populateColumns);
		}

		loadSchemas(database);

		return {
			setDatabase: function (value) { loadSchemas(value) }
		};
	}

	$.iquery.navigator = init;

})(window, window.document, window['$']);

/*------------- scripts --------------*/
(function (w, doc, $) {

	function init(elem, database, getQuery, onLoad) {
		var name = elem.querySelector('input'),
			button = elem.querySelector('button'),
			names = elem.querySelector('ul'),
			scripts = {};

		function key(name) {
			return database + '|' + name;
		}

		function loadScripts() {
			Object.keys(localStorage).forEach(function (k) {
				var s = k.split('|', 2),
					db = s[0],
					name = s[1];

				if (!(db in scripts)) {
					scripts[db] = []
				}
				scripts[db].push(name);
			});
			for (var db in scripts) scripts[db].sort();

			populate();
		}

		function populate() {
			$.clear(names);
			if (!(database in scripts) || scripts[database].length == 0) {
				names.appendChild($.elem('li', { className: 'noResults' }));
				return;
			}

			scripts[database].forEach(function (n) {
				var a = $.elem('a', { href: 'javascript: void(0)', textContent: n }),
					x = $.elem('a', { href: 'javascript: void(0)', className: 'remove' });

				a.addEventListener('click', function () { load(n) });
				x.addEventListener('click', function () { remove(n) });
				names.appendChild($.elem('li', { children: [a, x] }).appendChild(a).parentNode);
			});
		}

		function load(selected) {
			var sql = localStorage.getItem(key(selected));
			if (!sql) {
				sql = '';
			}
			onLoad(sql);
			name.value = selected;
		}

		function save() {
			var i, value = name.value;

			if (!(database in scripts)) {
				scripts[database] = [];
			}
			i = scripts[database].indexOf(value),

			localStorage.setItem(key(value), getQuery());
			if (i > -1) {
				scripts[database][i] = value;
			} else {
				scripts[database].push(value);
				scripts[database].sort();
			}
			populate();
		}

		function remove(selected) {
			if (!(database in scripts)) return;

			var i = scripts[database].indexOf(selected);
			if (i < 0) return

			scripts[database].splice(i, 1);
			localStorage.removeItem(key(selected));
			populate();
		}

		loadScripts();
		button.addEventListener('click', save);
		name.addEventListener('keyup', function (event) {
			if (event.keyCode == 13) save();
		});

		return {
			setDatabase: function (value) {
				database = value;
				populate();
			},
			load: load,
			save: save,
			remove: remove
		};
	}

	$.iquery.scripts = init;

})(window, window.document, window['$']);
