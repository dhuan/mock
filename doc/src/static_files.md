# Serving static files

Static files can easily be served. Suppose we have a folder named `public` where the static files we wish to serve are located.

```json
{
  "endpoints": [
    {
      "route": "static/*",
      "method": "GET",
      "response": "fs:./public"
    }
  ]
}
```

In the example above, we configured the route "static" to serve files located in the `public` folder. Let's say a file exists located in `public/foobar.html`, then it can be accessed through the URL `/static/foobar.html`.

How about spinning up quickly a static-file server without configuration files?

```diff
 $ mock serve \
+  --route 'static/*' \
+  --file-server /path/to/my/public/files
```

