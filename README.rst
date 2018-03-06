ipfix
=====

ipfix is a Go service (HTTP+RPC) to retrieve geolocation information
about an ip address with freegeoip_ helpers.

We are using bindings from the maxmind_ database.

Installation
============

Build it locally
----------------

1. Make sure you have a Go language compiler >= 1.3 (required) and git installed.
2. Make sure you have the following go system dependencies in your $PATH: bzr, svn, hg, git
3. Ensure your GOPATH_ is properly set.
4. Download it:

::

    git clone https://github.com/ulule/ipfix.git

4. Run ``make build``

You have now a binary version of ipfix in the ``bin`` directory which
fits perfectly with your architecture.

Build it using Docker
---------------------

If you don't want to install Go and Docker_ is installed on your computer

::

    make docker-build

You will have a binary version of ipfix compiled for linux in the ``bin`` directory.

Configuration
=============

Configuration should be stored in a readable file and in JSON format.

A complete example of the configuration file with RPC+HTTP would be:

``config.json``

.. code-block:: json

    {
        "server": {
            "rpc": {
                "port": 33001,
            },
            "http": {
                "port": 3001,
                "cors": {
                    "allowed_origins": ["*.ulule.com"],
                    "allowed_methods": ["GET", "HEAD", "POST"],
                    "allowed_headers": ["Origin", "Accept", "Content-Type", "X-Requested-With"]
                }
            }
        },
        "database_path": "./GeoLite2-City.mmdb.gz"
    }

Be careful, you should download first locally the GeoLite_ database because the service
will be unavailable until it will download the database.

HTTP server
===========

The HTTP server is based on chi_.

It's disabled by default, you can activate it by adding the `http` section to `server`.

``config.json``

.. code-block:: json

    {
        "server": {
            "http": {
                "port": 3001,
            }
        }
    }

CORS
----

ipfix supports CORS headers customization in your config file.

To enable this feature, set ``allowed_origins`` and ``allowed_methods``,
for example:

``config.json``

.. code-block:: json

    {
      "allowed_origins": ["*.ulule.com"],
      "allowed_methods": ["GET", "HEAD"]
    }

RPC server
===========

The RPC server is based on grpc_.

It's disabled by default, you can activate it by adding the `rpc` section to `server`.

``config.json``

.. code-block:: json

    {
        "server": {
            "http": {
                "port": 33001,
            }
        }
    }

You can found a client example in the `repository <examples/client/main.go>`_ and execute it:

::

    go run examples/client/main.go -ip {YOUR_IP_ADDRESS} -server-addr {RPC_ADDRESS}

Usage
=====

When your configuration is done, you can start the service as follow:

::

    ipfix -c config.json

or using an environment variable:

::

    IPFIX_CONF=/path/to/config.json ipfix

By default, this will run the application on port 3001 and can be accessed by visiting:

::

    http://localhost:3001

The port number can be configured with ``port`` option in your config file.

To see a list of all available options, run:

::

    ipfix --help

Development
===========

I recommend to install the live reload utility modd_ to make your life easier.

Install it:

::

    go get github.com/cortesi/modd/cmd/modd

Then launch it in the ipfix directory:

::

    IPFIX_CONF=config.json make live


.. _GOPATH: http://golang.org/doc/code.html#GOPATH
.. _GeoLite: http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz
.. _freegeoip: https://github.com/fiorix/freegeoip
.. _maxmind: https://www.maxmind.com/fr/home
.. _modd: https://github.com/cortesi/modd
.. _chi: https://github.com/go-chi/chi
.. _grpc: https://grpc.io/
.. _Docker: https://docker.com

Dang, what's this name?
=======================

It was an initial proposal from `kyojin <https://github.com/kyojin>`_ based on `Idéfix <https://en.wikipedia.org/wiki/Dogmatix>`_.

.. image:: https://media.giphy.com/media/Ob7p7lDT99cd2/giphy.gif
