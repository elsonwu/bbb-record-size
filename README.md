# bbb-record-size
Small tool to calculate how much space the BigBlueButton record use

# Usage
> ./bbb-record-size

# Parameters

	-host string
		the host to listen to (default ":1234")
		
	-published_path string
		the base path of published folder (default "/var/bigbluebutton/published/presentation/")
		
	-raw_path string
		the base path of raw folder (default "/var/bigbluebutton/recording/raw/")    	


# API

	http://bbb.yourhost.com:1234/record/<BBB record internal ID>

# Response
	{
   		error: "string, error string",
    	id: "string, record internal ID",
    	size: "float, how many kb"
	}
