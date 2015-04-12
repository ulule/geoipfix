ipfix
=====

ipfix is a dead simple Go webservice to retrieve geolocation information
about an ip address with freegeoip_ helpers.

We are using bindings from the maxmind_ database.

Installation
============

Build it
--------

1. Make sure you have a Go language compiler >= 1.3 (required) and git installed.
2. Make sure you have the following go system dependencies in your $PATH: bzr, svn, hg, git
3. Ensure your GOPATH_ is properly set.
4. Download it:

::

    git clone https://github.com/ulule/ipfix.git

4. Run ``make build``

You have now a binary version of picfit in the ``bin`` directory which
fits perfectly with your architecture.


Configuration
=============

Configuration should be stored in a readable file and in JSON format.

``config.json``

.. code-block:: json

    {
        "port": 3001,
        "allowed_origins": ["*.ulule.com"],
        "allowed_methods": ["GET", "HEAD", "POST"],
        "database_path": "./GeoLite2-City.mmdb.gz"
    }

You should download first locally the GeoLite_ database because the service
will be unavailable until it will download the database.

Usage
=====

When your configuration is done, you can start the service as follow:

::

    $ ipfix -c config.json

By default, this will run the application on port 3001 and can be accessed by visiting:

::

    $ http://localhost:3001

The port number can be configured with ``port`` option in your config file.

To see a list of all available options, run:

::

    $ ipfix --help

.. _GOPATH: http://golang.org/doc/code.html#GOPATH
.. _GeoLite: http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz
.. _freegeoip: https://github.com/fiorix/freegeoip
.. _maxmind: https://www.maxmind.com/fr/home
