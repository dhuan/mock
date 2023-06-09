# Responses from Shell scripts

You can write shell scripts that will act as "handlers" for your API's Requests (or Controllers if you like to think in terms of the MVC pattern.)

```json
{
  "endpoints": [
    {
      "route": "foo/bar",
      "response": "sh:my_shell_script.sh"
    }
  ]
}
```

In the example above, any request to `POST /foo/bar` will result in *mock* executing the `my_shell_script.sh`. Any output produced from that script execution will result in the HTTP Response returned by your API.

To further customize your script handlers, you may also pass parameters, just like you can normally pass parameters in a shell command:

```diff
 {
   "endpoints": [
     {
       "route": "foo/bar",
+      "response": "sh:my_shell_script.sh some_param another_param"
     }
   ]
 }
```

To define responses with shell scripts using command-line parameters, use the following:

```diff
 $ mock serve \
+  --route "foo/bar" \
+  --shell-script my_shell_script.sh
```

Alternatively, shell commands can be set as one-liners with `exec` instead of `sh`, not requiring you to create a shell script file. As an example, the endpoint below responds with a list of files of the current folder (`ls -la`):

```diff
 {
   "endpoints": [
     {
       "route": "foo/bar",
+      "response": "exec:ls -la"
     }
   ]
 }
```

You can use more advanced shell functionalities within `exec`, such as pipes. Let's set an endpoint that returns the amount of files that exist on the home folder:

```diff
 {
   "endpoints": [
     {
       "route": "foo/bar",
+      "response": "exec:ls ~ | wc -l"
     }
   ]
 }
```

The same can be accomplished through command-line parameters:

```diff
 $ mock serve \
+  --route "foo/bar" \
+  --exec 'ls | sort'
```

## Environment Variables for Request Handlers

A set of environment variables can be read from in response shell scripts in order to obtain useful information about the current request. Static responses (such as JSON) also have access to the same variables. Reading them is done through writing the variable name prefixed with a "$" - for example `$MOCK_REQUEST_URL`. The following are the variables avaiable:

- `MOCK_REQUEST_URL`: The full URL. (ex: `http://localhost/foo/bar`)
- `MOCK_REQUEST_ENDPOINT`: The endpoint extracted from the URL. (ex: `foo/bar`)
- `MOCK_REQUEST_HOST`: The hostname + port combination that the request was sent to. (ex: `example.com:3000`)
- `MOCK_REQUEST_HEADERS`: A file path containing all HTTP Headers.
- `MOCK_REQUEST_BODY`: A file path containing the Request's Body (if one exists, otherwise this will be an empty file.)
- `MOCK_REQUEST_QUERYSTRING`: The Request's Querystring if it exists. (ex: `some_key=some_value&another_key=another_value`)
- `MOCK_REQUEST_METHOD`: A string indicating the Request's Method.

The following environment variables provide other general information not related to the current request:

- `MOCK_HOST`: The hostname + port combination to which Mock is currently listening. (ex: `localhost:3000`)

## Route Parameters - Reading from Shell Scripts

Route Parameters can be read from shell scripts. Suppose an endpoint exists as such: `user/{user_id}`. We could then retrieve the User ID parameter by reading the `MOCK_ROUTE_PARAM_USER_ID` environment variable.

## Response Files that can be written to by shell scripts

So far we've seen environment variables that provide us with information about the Request that's being currently handled. The following environment variables enable you to further define the HTTP Response:

- `MOCK_RESPONSE_STATUS_CODE`: A file that your handler can write to, to define the HTTP Status Code. 
- `MOCK_RESPONSE_HEADERS`: A file that your handler can write to, to define the HTTP Headers.

In the following example, we'll see what a Handler looks like, which responds with a simple `Hello world!` body content, a `201` Status Code and a few custom HTTP Headers.

```sh
echo Hello world!

cat <<EOF > $MOCK_RESPONSE_HEADERS
Some-Header-Key: Some Header Value
Another-Header-Key: Another Header Value
EOF

echo 201 > $MOCK_RESPONSE_STATUS_CODE
```

