enqueuenzb
==========

Submit an NZB file to the SABnzbd download queue.

*enqueuenzb* submits the specified nzb file to the SABnzbd service API according to the configuration provided.

usage
-----
~~~
enqueuenzb <nzb-file> [nzb-file...]
~~~

Where nzb-file is the NZB file you want to submit to SABnzbd. Specifying more than one NZB file is allowed.

config
------
Create a JSON-formatted config file in the following location. The config file contains the server to which to submit the nzb file to and the API key that should be used to acquire authorization. Finally it is possible to specify whether or not to automatically delete the NZB file after a successful submission.

**$HOME/.config/enqueuenzb.conf**:
~~~
{
	"url": "http://localhost:8080/sabnzbd/api",
	"key": "YourSABnzbdApiOrNzbKey",
	"delete": true,
	"verbose": true
}
~~~

In case 'verbose' is enabled, feedback will be given for every NZB file that was successfully submitted.
