Response Status Code
====================

By default, all responses’ status code will be ``200``. You can change
it using the ``response_status_code`` option:

.. code:: diff

    {
      "endpoints": [
        {
          "route": "foo/bar",
          "method": "POST",
          "response": {
            "foo": "bar"
          },
   +      "response_status_code": 201 
        }
      ]
   }

To add response status codes to an endpoint using command-line
parameters:

.. code:: diff

    $ mock serve \
      --route "foo/bar" \
      --method "POST" \
      --response '{"foo":"bar"}' \
   +  --status-code 201
