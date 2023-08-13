Search.setIndex({"docnames": ["api_reference", "apis", "assertion_chaining", "base_apis", "changelog", "cmd_reference", "conditional_responses", "conditions_reference", "cors", "env_vars", "headers", "index", "install", "middlewares", "mock_package", "mock_vars", "route_params", "shell_scripts", "static_files", "status_codes", "test_assertions"], "filenames": ["api_reference.rst", "apis.rst", "assertion_chaining.rst", "base_apis.rst", "changelog.rst", "cmd_reference.rst", "conditional_responses.rst", "conditions_reference.rst", "cors.rst", "env_vars.rst", "headers.rst", "index.rst", "install.rst", "middlewares.rst", "mock_package.rst", "mock_vars.rst", "route_params.rst", "shell_scripts.rst", "static_files.rst", "status_codes.rst", "test_assertions.rst"], "titles": ["Mock API Reference", "Creating APIs", "Assertion Chaining", "Base APIs", "Changelog", "Command-line Options Reference", "Conditional Response", "Conditions Reference", "Handling CORS", "Reading Environment Variables", "Response with headers", "mock - Language-agnostic API mocking and testing utility", "Installation", "Middlewares", "Test Assertions with <em>mock</em>\u2019s Go package", "Mock Variables", "Route Parameters", "Responses from Shell scripts", "Serving static files", "Response Status Code", "Test Assertions"], "terms": {"besid": [0, 1, 6, 14, 15, 16, 17, 20], "custom": [0, 5, 6, 9, 13, 17], "endpoint": [0, 2, 3, 4, 5, 6, 7, 9, 10, 11, 13, 15, 16, 17, 18, 19, 20], "defin": [0, 3, 4, 5, 6, 7, 11, 13, 15, 16, 17], "your": [0, 1, 3, 4, 5, 6, 8, 11, 12, 13, 14, 17, 20], "configur": [0, 1, 3, 4, 5, 6, 9, 11, 13, 16, 18], "file": [0, 3, 4, 9, 11, 12, 13, 15, 16], "provid": [0, 5, 13, 15, 17, 20], "intern": 0, "ar": [0, 1, 2, 3, 4, 5, 6, 7, 11, 13, 14, 16, 17, 18, 20], "identifi": 0, "have": [0, 1, 3, 4, 6, 7, 13, 15, 18], "rout": [0, 1, 2, 3, 4, 6, 7, 9, 10, 11, 13, 14, 15, 18, 19, 20], "prefix": 0, "which": [0, 1, 3, 4, 5, 6, 13, 14, 15, 17], "exist": [0, 1, 4, 6, 11, 13, 15, 17, 18], "make": [0, 2, 3, 4, 7, 11, 12, 13, 14, 20], "In": [0, 1, 2, 3, 6, 7, 14, 17, 18, 20], "thi": [0, 1, 3, 4, 5, 6, 7, 13, 14, 15, 20], "section": [0, 1, 3, 4, 5, 7, 11, 13, 14, 16, 20], "you": [0, 1, 2, 3, 4, 5, 6, 7, 11, 12, 13, 14, 15, 16, 17, 19, 20], "ll": [0, 1, 6, 12, 14, 17, 20], "find": [0, 3, 5, 7, 12, 15], "out": [0, 3, 4, 11, 13, 15], "about": [0, 3, 5, 13, 15, 16, 17, 18], "each": [0, 5, 11, 13], "avail": [0, 5, 6, 7, 12, 15, 20], "test": [0, 4, 7], "x": [0, 14], "wa": [0, 2, 3, 4, 7, 11, 12, 14, 15, 16, 20], "call": [0, 2, 6, 7, 14, 20], "y": 0, "payload": [0, 2, 3, 4, 7, 11, 20], "The": [0, 1, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 17], "dedic": 0, "explain": 0, "all": [0, 1, 3, 5, 6, 7, 8, 13, 14, 15, 16, 19, 20], "remov": 0, "request": [0, 1, 2, 4, 5, 6, 7, 11, 13, 14, 15, 16], "record": 0, "been": [0, 1, 4, 7, 12, 15, 20], "made": [0, 3, 4, 7, 13, 14, 15, 20], "so": [0, 5, 7, 12, 13, 17], "far": [0, 7, 13, 17], "ha": [0, 4, 6, 12, 20], "same": [0, 1, 3, 4, 6, 7, 9, 14, 16, 17], "effect": 0, "stop": 0, "start": [0, 1, 4, 5, 7, 9, 12, 14, 15], "over": [0, 3, 13], "again": [0, 1, 3], "There": [0, 2, 6, 20], "paramet": [0, 4, 5, 6, 9, 10, 13, 15, 19], "field": [0, 4, 6, 7, 13, 20], "simplest": 1, "we": [1, 2, 3, 6, 9, 13, 14, 17, 18, 20], "can": [1, 2, 3, 4, 5, 6, 7, 8, 9, 13, 14, 15, 16, 18, 19, 20], "look": [1, 6, 14, 17, 20], "like": [1, 4, 6, 7, 14, 17, 20], "foo": [1, 2, 3, 4, 5, 6, 7, 9, 10, 13, 14, 15, 17, 19, 20], "bar": [1, 2, 3, 4, 5, 6, 7, 9, 10, 13, 14, 15, 17, 19, 20], "method": [1, 2, 3, 4, 6, 7, 9, 10, 11, 12, 14, 15, 16, 18, 19, 20], "post": [1, 2, 4, 7, 10, 14, 17, 19, 20], "A": [1, 3, 4, 13, 14, 15, 17], "http": [1, 3, 4, 5, 7, 11, 12, 13, 14, 15, 17, 20], "respond": [1, 3, 17], "seen": [1, 6, 13, 14, 17, 20], "abov": [1, 3, 6, 11, 13, 14, 15, 16, 17, 18], "also": [1, 2, 6, 7, 13, 16, 17], "set": [1, 3, 4, 5, 6, 7, 8, 11, 14, 15, 17, 20], "wildcard": [1, 4], "With": [1, 5, 6, 14, 16, 20], "anyth": [1, 13], "hello": [1, 2, 3, 4, 5, 6, 11, 13, 17, 20], "world": [1, 2, 3, 4, 5, 6, 11, 13, 17, 20], "placehold": [1, 4], "variabl": [1, 3, 4, 16], "well": [1, 6], "some_vari": 1, "order": [1, 3, 4, 13, 17], "read": [1, 3, 4, 6, 13, 15, 16], "do": [1, 4], "someth": [1, 14], "us": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 17, 19, 20], "need": [1, 4, 14], "shell": [1, 4, 9, 12, 13, 15, 16], "script": [1, 3, 4, 13, 14, 15, 16], "act": [1, 3, 5, 17], "handler": [1, 3, 4, 5, 13, 15], "next": [1, 6], "other": [1, 3, 6, 7, 11, 12, 15, 20], "wai": 1, "up": [1, 3, 4, 8, 9, 11, 18, 20], "an": [1, 3, 4, 5, 9, 10, 13, 14, 15, 17, 19], "altern": [1, 3, 13, 17], "let": [1, 3, 7, 9, 13, 14, 17, 18, 20], "": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 15, 16, 17, 18, 20], "mock": [1, 3, 4, 5, 7, 8, 9, 10, 13, 16, 17, 18, 19, 20], "two": [1, 6, 13, 15], "serv": [1, 3, 4, 5, 8, 9, 10, 11, 13, 15, 16, 17, 19], "get": [1, 3, 4, 5, 6, 7, 9, 11, 14, 16, 18, 20], "i": [1, 3, 4, 5, 6, 7, 9, 11, 12, 13, 14, 15, 16, 17, 20], "anoth": [1, 5, 7, 10, 13, 17, 20], "As": [1, 2, 17], "shown": [1, 2, 4, 13, 14], "accomplish": [1, 7, 9, 17], "json": [1, 2, 4, 5, 7, 8, 9, 11], "done": [1, 14, 20], "just": [1, 3, 7, 11, 12, 14, 17], "matter": [1, 14], "prefer": [1, 11], "move": [1, 12], "forward": 1, "manual": 1, "learn": [1, 3, 5, 11, 13, 14], "more": [1, 3, 4, 5, 6, 11, 12, 13, 14, 16, 17], "advanc": [1, 17], "function": [1, 4, 11, 17], "instruct": 1, "how": [1, 2, 3, 6, 7, 14, 18, 20], "achiev": [1, 3, 4, 6, 14], "thing": [1, 11, 14, 20], "both": [1, 7], "onli": [1, 3, 4, 5, 6, 7, 13, 14], "scratch": 1, "surfac": 1, "few": [1, 12, 14, 17], "note": [1, 7, 9, 14, 16], "awar": 1, "while": [1, 16], "togeth": 1, "when": [1, 3, 4, 5, 7, 8, 11, 13, 14, 15, 16], "alreadi": [1, 13], "former": 1, "overwrit": [1, 13], "latter": 1, "word": [1, 3, 13], "alwai": 1, "ones": [1, 13], "config": [1, 4, 8, 9], "combin": [1, 2, 6, 15], "earlier": 1, "exampl": [1, 2, 3, 4, 5, 6, 7, 9, 11, 14, 15, 16, 17, 18], "object": [1, 3, 6, 13], "contain": [1, 3, 4, 7, 13, 14, 15], "howev": [1, 3, 6, 20], "setup": 1, "complex": [1, 6, 13, 14], "larg": 1, "easili": [1, 11, 14, 18], "readabl": [1, 4], "follow": [1, 2, 3, 4, 5, 6, 7, 9, 12, 13, 15, 16, 17], "re": [1, 2, 13, 14, 15, 16, 20], "referenc": [1, 4], "thu": [1, 5], "leav": 1, "path": [1, 4, 5, 8, 9, 12, 13, 15, 18], "some": [1, 5, 6, 7, 10, 11, 13, 17], "To": [1, 10, 16, 17, 19], "header": [1, 3, 7, 8, 11, 15, 17, 20], "refer": [1, 4, 6, 13, 20], "from": [1, 3, 4, 5, 6, 7, 12, 14, 15, 20], "environ": [1, 3, 4], "written": [1, 13, 16], "static": [1, 4, 5, 16], "statu": [1, 4, 13, 17, 20], "code": [1, 4, 11, 13, 14, 17, 20], "condit": [1, 4, 14, 20], "chain": [1, 14, 20], "middlewar": [1, 3, 4, 15], "intercept": 1, "tl": 1, "option": [2, 3, 4, 6, 10, 13, 19, 20], "extend": [2, 6], "tri": 2, "previous": [2, 4, 13], "our": [2, 13], "type": [2, 3, 6, 7, 13, 14, 20], "method_match": [2, 14, 20], "valu": [2, 3, 4, 5, 6, 7, 9, 10, 11, 13, 14, 15, 16, 17, 20], "json_body_match": [2, 4], "key_valu": [2, 7, 13], "plain": [2, 7, 14, 20], "english": [2, 20], "nest": [2, 6], "within": [2, 17], "parent": 2, "limit": [2, 6, 13, 14], "mani": 2, "server": [3, 4, 7, 11, 14, 18], "result": [3, 4, 5, 6, 16, 17], "receiv": [3, 5], "forwad": 3, "take": [3, 8, 14], "prioriti": 3, "see": [3, 5, 14, 17, 20], "simpl": [3, 6, 13, 15, 17, 20], "com": [3, 4, 12, 14, 15], "mai": [3, 5, 6, 17], "ve": [3, 6, 12, 13, 14, 17, 20], "its": [3, 9, 13, 14, 16], "furthermor": 3, "If": [3, 6, 11, 12, 13, 14, 15], "normal": [3, 17], "domain": 3, "proxi": 3, "manipul": [3, 13, 17], "given": [3, 4, 5, 7, 11, 13], "check": [3, 4, 5, 11, 12, 14], "document": [3, 4, 5, 11, 17], "them": [3, 13], "fact": 3, "distinct": 3, "whether": [3, 11, 20], "current": [3, 4, 7, 15, 17], "meant": 3, "actual": [3, 14], "execut": [3, 4, 5, 12, 13, 14, 17], "mock_base_api_respons": 3, "below": [3, 5, 7, 13, 15, 17], "add": [3, 4, 7, 10, 13, 19], "true": 3, "printf": [3, 5, 11], "mock_response_head": [3, 13, 17], "fi": 3, "route_match": [3, 13], "filter": [3, 5, 13], "want": [3, 6, 9, 20], "target": 3, "pattern": 3, "choic": 3, "Then": 3, "cover": [3, 5, 14], "subject": 3, "detail": [3, 4, 5, 17], "sent": [3, 13, 15, 16], "exactli": 3, "client": [3, 8, 11], "who": 3, "etc": [3, 4, 7, 13, 20], "everyth": 3, "unmodifi": 3, "It": [3, 7, 8, 12, 14], "possibl": [3, 6, 7], "modifi": 3, "ani": [3, 4, 6, 9, 11, 13, 17], "properti": [3, 13], "through": [3, 4, 5, 11, 13, 16, 17, 18], "on_request_to_base_api": 3, "abl": 3, "write": [3, 4, 11, 14, 17], "assign": [3, 17], "mock_request_bodi": 3, "mock_request_head": 3, "protocol": 3, "combo": 3, "further": [3, 13, 17], "understand": 3, "differ": [3, 6, 7, 20], "todo": 3, "ad": [4, 6], "new": [4, 5, 12, 20], "captur": [4, 15, 16], "individu": [4, 15], "querystr": [4, 6, 7, 13, 15], "mock_request_querystring_foobar": 4, "enabl": [4, 5, 6, 7, 8, 11, 13, 14, 17, 20], "plaintext": 4, "respons": [4, 7, 8, 9, 11, 14, 15, 16, 18, 20], "url": [4, 11, 14, 15, 18], "mock_request_url": 4, "support": [4, 12], "exec": [4, 11, 13, 17], "command": [4, 9, 10, 11, 12, 13, 16, 17, 19], "includ": [4, 5, 9, 13, 15, 20], "oper": [4, 13, 17], "pipe": [4, 13, 17], "output": [4, 13], "redirect": [4, 13], "querystring_match_regex": 4, "querystring_exact_match_regex": 4, "break": 4, "chang": [4, 13, 19], "now": [4, 5, 6, 9, 11, 14], "mock_response_bodi": [4, 5, 11, 13, 17], "bodi": [4, 5, 7, 11, 15, 17, 20], "instead": [4, 5, 7, 9, 13, 16, 17, 20], "stdout": 4, "interfac": [4, 14], "assert": [4, 7, 11], "better": 4, "renam": 4, "befor": [4, 5, 6, 12, 14], "upgrad": 4, "mock_request_nth": 4, "route_param_match": 4, "nth": [4, 20], "base": [4, 6, 16, 20], "posit": [4, 15], "histori": [4, 7, 15], "line": [4, 9, 10, 12, 13, 16, 17, 19], "201": [4, 17, 19], "longer": 4, "mandatori": 4, "sinc": [4, 15], "without": [4, 5, 8, 14, 18], "fix": 4, "absolut": [4, 5], "fail": [4, 14, 20], "rel": [4, 5], "issu": 4, "had": [4, 20], "abil": 4, "kind": [4, 13, 14, 16], "either": [4, 5], "text": 4, "user": [4, 17], "guid": [4, 16], "host": [4, 14], "mock_request_host": 4, "name": [4, 5, 11, 13, 15, 16, 18], "mock_request_endpoint_param_foo": 4, "mock_route_param_foo": 4, "featur": 4, "enhanc": 4, "string": [4, 5, 7, 9, 14, 15, 16], "could": [4, 17, 20], "book": [4, 16], "book_nam": [4, 16], "txt": [4, 16], "f": [4, 18], "default": [4, 5, 6, 19, 20], "addit": 4, "mock_host": 4, "retriev": [4, 17], "listen": [4, 15], "bug": [4, 11], "try": [4, 11], "didn": 4, "t": [4, 6, 9, 14], "would": [4, 6, 16, 20], "500": 4, "api": [4, 6, 8, 13, 17, 20], "error": [4, 14, 20], "param": [4, 7], "sh": [4, 13, 17], "my_shell_script": [4, 17], "some_param": [4, 13, 17], "another_param": [4, 17], "delai": 4, "simul": 4, "slow": 4, "some_handl": 4, "cor": 4, "facilit": 4, "usag": [4, 9], "webapp": 4, "minor": 4, "stuff": 4, "prevent": 4, "show": 4, "gracefulli": 4, "instal": 4, "go": [4, 13, 20], "librari": [4, 14], "releas": [4, 11, 12], "github": [4, 12, 14], "dhuan": [4, 12, 14], "pkg": [4, 14], "project": 4, "helper": 4, "toreadableerror": [4, 14], "stringifi": 4, "group": 4, "valid": [4, 5, 13, 14, 20], "work": 4, "independ": 4, "case": [4, 6, 14, 20], "sensit": 4, "doe": [4, 5, 6, 13, 14], "doesn": 4, "signific": 4, "querystring_exact_match": 4, "matcher": 4, "desir": [4, 13], "kei": [4, 5, 6, 7, 10, 17], "gener": [4, 5], "improv": 4, "stabil": 4, "log": 4, "messag": 4, "timestamp": 4, "proper": 4, "handl": [4, 15, 17], "unabl": 4, "rang": 4, "panick": 4, "return": [4, 6, 14, 17], "indic": [4, 13, 15, 20], "tabl": 5, "content": [5, 17], "navig": 5, "spec": 5, "run": [5, 8, 11, 12, 14], "argument": 5, "3000": [5, 11, 15, 20], "say_hi": [5, 11], "my": [5, 11, 18], "what_time_is_it": [5, 11], "date": [5, 11], "h": [5, 11], "m": [5, 11], "must": [5, 7], "directori": 5, "locat": [5, 18], "here": [5, 9, 12, 13], "format": [5, 12, 14], "number": [5, 7, 15], "jump": 5, "page": [5, 12, 13, 17], "regular": [5, 7, 13], "express": [5, 7, 13], "against": [5, 7, 13], "necessari": [5, 6, 8, 16], "browser": [5, 8, 11], "complain": 5, "cross": 5, "origin": 5, "amount": [5, 17], "millisecond": 5, "wait": 5, "immedi": 5, "everi": 5, "3": [5, 7], "second": [5, 7], "certain": [6, 13], "response_if": 6, "querystring_match": [6, 13], "sampl": 6, "singl": [6, 9, 16], "not_bar": 6, "galaxi": 6, "even": 6, "though": 6, "still": [6, 14], "where": [6, 12, 15, 18, 20], "match": [6, 7, 13], "previou": [6, 14], "veri": [6, 13, 20], "present": 6, "deep": 6, "response_head": [6, 10], "foobar": [6, 13, 18], "don": [6, 9], "main": 6, "d": [6, 14], "inherit": 6, "response_headers_bas": 6, "resolv": 6, "dispos": 6, "list": [6, 13, 17], "expect": 7, "For": [7, 11, 13, 15, 20], "multipl": 7, "pair": 7, "some_kei": [7, 14], "another_kei": 7, "comparison": 7, "z": [7, 13], "1": [7, 12, 13], "0": [7, 14], "9": 7, "being": [7, 13, 15, 17], "specifi": [7, 13], "que": 7, "form": 7, "encod": 7, "data": 7, "some_param_nam": 7, "some_valu": [7, 14], "occur": [7, 14], "2nd": [7, 14, 15], "2": [7, 15, 20], "subsequ": 7, "after": 7, "plu": 7, "sign": 7, "onward": 7, "flag": 8, "care": 8, "comun": 8, "problem": 8, "c": [8, 9], "And": [9, 11, 14], "quot": [9, 16], "around": [9, 16], "becaus": [9, 14, 16, 20], "process": [9, 13, 16], "program": [9, 11, 12, 13, 14, 16, 17], "export": 9, "quickli": [11, 12, 18], "end": 11, "respect": 11, "easi": 11, "syntax": [11, 15, 16], "pass": [11, 14, 17, 20], "correctli": 11, "port": [11, 14, 15], "localhost": [11, 14, 15, 20], "download": 11, "linux": [11, 12], "maco": [11, 12], "sourc": 11, "report": 11, "core": 11, "creat": [11, 12, 17], "why": 11, "similar": 11, "tool": 11, "thei": 11, "somehow": 11, "requir": [11, 14, 17], "fake": 11, "hand": [11, 13, 14], "under": 11, "mit": 11, "inform": [11, 13, 15, 17], "easiest": 12, "simpli": [12, 14], "built": 12, "standalon": 12, "depend": 12, "choos": 12, "one": [12, 13, 15, 17], "oss": 12, "Or": [12, 20], "wget": 12, "o": 12, "tgz": 12, "var_download_link_linux": 12, "tar": 12, "xzvf": 12, "version": 12, "proceed": 12, "sure": 12, "system": [12, 17], "golang": 12, "18": 12, "recent": 12, "gnu": 12, "git": 12, "clone": 12, "folder": [12, 17, 18], "build": 12, "successfulli": [12, 14], "should": 12, "insid": 12, "bin": 12, "root": 12, "repositori": 12, "logic": 13, "perform": 13, "But": [13, 14], "rememb": 13, "replac": 13, "occurr": 13, "sed": 13, "g": 13, "echo": [13, 17], "One": 13, "observ": 13, "apped": 13, "append": 13, "overwritten": 13, "tmp": 13, "mktemp": 13, "cat": [13, 17], "grep": 13, "v": [13, 14], "behavior": 13, "hold": [13, 15], "customis": 13, "moment": 13, "mock_response_status_cod": [13, 17], "mock_route_param_some_param": 13, "complet": 13, "consult": 13, "By": [13, 19, 20], "mechan": 13, "404": 13, "mean": [14, 20], "design": 14, "languag": 14, "agnost": 14, "what": [14, 17, 20], "e2": 14, "integr": 14, "curl": 14, "4000": 14, "__mock__": [14, 20], "eof": [14, 17], "convert": 14, "my_test": 14, "import": [14, 16], "func": 14, "test_foobarshouldberequest": 14, "mockconfig": 14, "validationerror": 14, "err": 14, "assertopt": 14, "conditiontype_methodmatch": 14, "nil": 14, "len": 14, "approach": 14, "instanc": 14, "suppos": [14, 16, 17, 18], "regard": 14, "snippet": 14, "prior": 14, "tell": [14, 20], "network": 14, "relat": 14, "failur": 14, "otherwis": [14, 15], "seem": 14, "empti": [14, 15], "slice": 14, "basic": 14, "conditiontype_jsonbodymatch": 14, "keyvalu": 14, "map": 14, "print": 15, "variable_nam": 15, "exemplifi": [15, 16], "place": 15, "hostnam": 15, "ex": 15, "full": 15, "extract": 15, "key_nam": 15, "mock_request_querystring_foo": 15, "ever": 15, "segment": 16, "search": 16, "author": 16, "author_nam": 16, "year": 16, "asimov": 16, "1980": 16, "dynam": 16, "own": 16, "appli": 16, "wrap": 16, "NOT": 16, "repons": 16, "liner": 17, "l": 17, "la": 17, "home": 17, "wc": 17, "sort": 17, "give": 17, "These": 17, "user_id": 17, "id": 17, "mock_route_param_user_id": 17, "u": 17, "public": 18, "wish": 18, "sai": [18, 20], "html": 18, "access": 18, "spin": 18, "200": [19, 20], "response_status_cod": 19, "were": 20, "concept": 20, "autom": 20, "put": 20, "never": 20, "particular": 20, "validation_error": 20, "no_cal": 20, "metadata": 20, "attempt": 20, "method_mismatch": 20, "method_expect": 20, "method_request": 20, "success": 20, "400": 20, "skip": 20, "1st": 20, "than": 20, "first": 20, "packag": 20}, "objects": {}, "objtypes": {}, "objnames": {}, "titleterms": {"mock": [0, 11, 12, 14, 15], "api": [0, 1, 3, 5, 11], "refer": [0, 5, 7, 15], "post": 0, "__mock__": 0, "assert": [0, 2, 14, 20], "reset": 0, "creat": 1, "endpoint": 1, "defin": 1, "through": [1, 12], "command": [1, 5], "line": [1, 5], "paramet": [1, 16, 17], "file": [1, 5, 17, 18], "base": [1, 3], "respons": [1, 3, 5, 6, 10, 13, 17, 19], "content": [1, 20], "chain": [2, 6], "intercept": 3, "request": [3, 17, 20], "tl": 3, "changelog": 4, "unreleas": 4, "yet": 4, "1": 4, "0": 4, "8": 4, "7": 4, "6": 4, "5": 4, "4": 4, "3": 4, "2": 4, "option": 5, "standard": 5, "c": 5, "config": 5, "p": 5, "port": 5, "specifi": 5, "rout": [5, 16, 17], "method": 5, "sh": 5, "shell": [5, 17], "script": [5, 17], "exec": 5, "server": 5, "header": [5, 6, 10, 13], "statu": [5, 19], "code": [5, 12, 19], "middlewar": [5, 13], "match": 5, "miscellan": 5, "cor": [5, 8], "d": 5, "delai": 5, "condit": [6, 7, 13], "querystring_match": 7, "querystring_match_regex": 7, "querystring_exact_match": 7, "querystring_exact_match_regex": 7, "json_body_match": 7, "form_match": 7, "header_match": 7, "method_match": 7, "route_param_match": 7, "nth": 7, "handl": 8, "read": [9, 11, 17], "environ": [9, 13, 17], "variabl": [9, 13, 15, 17], "languag": 11, "agnost": 11, "test": [11, 14, 20], "util": 11, "quick": 11, "link": 11, "further": 11, "licens": 11, "instal": 12, "download": 12, "sourc": 12, "set": 13, "up": 13, "exampl": 13, "modifi": 13, "bodi": 13, "ad": 13, "new": 13, "befor": 13, "send": 13, "client": 13, "remov": 13, "from": [13, 17], "": 14, "go": 14, "packag": 14, "mock_host": 15, "mock_request_url": 15, "mock_request_endpoint": 15, "mock_request_host": 15, "mock_request_head": 15, "mock_request_bodi": 15, "mock_request_querystr": 15, "mock_request_querystring_key_nam": 15, "mock_request_method": 15, "mock_request_nth": 15, "handler": 17, "can": 17, "written": 17, "serv": 18, "static": 18, "which": 20, "against": 20}, "envversion": {"sphinx.domains.c": 3, "sphinx.domains.changeset": 1, "sphinx.domains.citation": 1, "sphinx.domains.cpp": 9, "sphinx.domains.index": 1, "sphinx.domains.javascript": 3, "sphinx.domains.math": 2, "sphinx.domains.python": 4, "sphinx.domains.rst": 2, "sphinx.domains.std": 2, "sphinx": 58}, "alltitles": {"Mock API Reference": [[0, "mock-api-reference"]], "POST __mock__/assert": [[0, "post-mock-assert"]], "POST __mock__/reset": [[0, "post-mock-reset"]], "Creating APIs": [[1, "creating-apis"]], "Endpoints defined through command-line parameters": [[1, "endpoints-defined-through-command-line-parameters"]], "File-based response content": [[1, "file-based-response-content"]], "Contents:": [[1, null], [20, null]], "Assertion Chaining": [[2, "assertion-chaining"]], "Base APIs": [[3, "base-apis"]], "Intercepting responses": [[3, "intercepting-responses"]], "Intercepting requests": [[3, "intercepting-requests"]], "Base APIs and TLS": [[3, "base-apis-and-tls"]], "Changelog": [[4, "changelog"]], "Unreleased yet": [[4, "unreleased-yet"]], "1.1.0": [[4, "id1"]], "1.0.0": [[4, "section-1"]], "0.8.1": [[4, "section-2"]], "0.8.0": [[4, "section-3"]], "0.7.0": [[4, "section-4"]], "0.6.0": [[4, "section-5"]], "0.5.0": [[4, "section-6"]], "0.4.0": [[4, "section-7"]], "0.3.0": [[4, "section-8"]], "0.2.0": [[4, "section-9"]], "0.1.4": [[4, "section-10"]], "0.1.3": [[4, "section-11"]], "0.1.2": [[4, "section-12"]], "0.1.1": [[4, "section-13"]], "0.1.0": [[4, "section-14"]], "Command-line Options Reference": [[5, "command-line-options-reference"]], "Standard options": [[5, "standard-options"]], "-c or --config": [[5, "c-or-config"]], "-p or --port": [[5, "p-or-port"]], "Options for specifying APIs": [[5, "options-for-specifying-apis"]], "--route": [[5, "route"]], "--method": [[5, "method"]], "--response": [[5, "response"]], "--response-file": [[5, "response-file"]], "--response-sh or --shell-script": [[5, "response-sh-or-shell-script"]], "--exec, --response-exec": [[5, "exec-response-exec"]], "--file-server, --response-file-server": [[5, "file-server-response-file-server"]], "--header": [[5, "header"]], "--status-code": [[5, "status-code"]], "Options for Middlewares": [[5, "options-for-middlewares"]], "--middleware": [[5, "middleware"]], "--route-match": [[5, "route-match"]], "Miscellaneous": [[5, "miscellaneous"]], "--cors": [[5, "cors"]], "-d or --delay": [[5, "d-or-delay"]], "Conditional Response": [[6, "conditional-response"]], "Condition Chaining": [[6, "condition-chaining"]], "Headers in Conditional Responses": [[6, "headers-in-conditional-responses"]], "Conditions Reference": [[7, "conditions-reference"]], "querystring_match": [[7, "querystring-match"]], "querystring_match_regex": [[7, "querystring-match-regex"]], "querystring_exact_match": [[7, "querystring-exact-match"]], "querystring_exact_match_regex": [[7, "querystring-exact-match-regex"]], "json_body_match": [[7, "json-body-match"]], "form_match": [[7, "form-match"]], "header_match": [[7, "header-match"]], "method_match": [[7, "method-match"]], "route_param_match": [[7, "route-param-match"]], "nth": [[7, "nth"]], "Handling CORS": [[8, "handling-cors"]], "Reading Environment Variables": [[9, "reading-environment-variables"]], "Response with headers": [[10, "response-with-headers"]], "mock - Language-agnostic API mocking and testing utility": [[11, "mock-language-agnostic-api-mocking-and-testing-utility"]], "Quick links": [[11, "quick-links"]], "Read further\u2026": [[11, "read-further"]], "License": [[11, "license"]], "Installation": [[12, "installation"]], "Download mock": [[12, "download-mock"]], "Install through source code": [[12, "install-through-source-code"]], "Middlewares": [[13, "middlewares"]], "Setting up Middlewares": [[13, "setting-up-middlewares"]], "Examples of Middlewares": [[13, "examples-of-middlewares"]], "Modify Response Body": [[13, "modify-response-body"]], "Adding new headers before sending response to client": [[13, "adding-new-headers-before-sending-response-to-client"]], "Removing headers from the response": [[13, "removing-headers-from-the-response"]], "Environment Variables for Middlewares": [[13, "environment-variables-for-middlewares"]], "Conditions for Middlewares": [[13, "conditions-for-middlewares"]], "Test Assertions with mock\u2019s Go package": [[14, "test-assertions-with-mocks-go-package"]], "Mock Variables": [[15, "mock-variables"]], "Variable Reference": [[15, "variable-reference"]], "MOCK_HOST": [[15, "mock-host"]], "MOCK_REQUEST_URL": [[15, "mock-request-url"]], "MOCK_REQUEST_ENDPOINT": [[15, "mock-request-endpoint"]], "MOCK_REQUEST_HOST": [[15, "mock-request-host"]], "MOCK_REQUEST_HEADERS": [[15, "mock-request-headers"]], "MOCK_REQUEST_BODY": [[15, "mock-request-body"]], "MOCK_REQUEST_QUERYSTRING": [[15, "mock-request-querystring"]], "MOCK_REQUEST_QUERYSTRING_KEY_NAME": [[15, "mock-request-querystring-key-name"]], "MOCK_REQUEST_METHOD": [[15, "mock-request-method"]], "MOCK_REQUEST_NTH": [[15, "mock-request-nth"]], "Route Parameters": [[16, "route-parameters"]], "Responses from Shell scripts": [[17, "responses-from-shell-scripts"]], "Environment Variables for Request Handlers": [[17, "environment-variables-for-request-handlers"]], "Route Parameters - Reading from Shell Scripts": [[17, "route-parameters-reading-from-shell-scripts"]], "Response Files that can be written to by shell scripts": [[17, "response-files-that-can-be-written-to-by-shell-scripts"]], "Serving static files": [[18, "serving-static-files"]], "Response Status Code": [[19, "response-status-code"]], "Test Assertions": [[20, "test-assertions"]], "Which Request to assert against?": [[20, "which-request-to-assert-against"]]}, "indexentries": {}})