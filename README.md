ipsaver
=======

This server was written, to save the IP address of the machines queering it via HTTP under the key they use as a URL path

saving
-------
The server stores every request with a new IP address or a new key (the request path) while new means that the same combination is not equal to that of the latest entry for that key

presenting
----------
You can display either every key with a "?show" as a suffix to no key or all entries with the same suffix after a key

technique
---------
I use go as the servers language and sqlite to store the results, data is displayed in pure text format with no links provided

status
------
Not ready yet


version
-------
0.2

dependencies
------------
go-sqlite3 under "github.com/mattn/go-sqlite3"
