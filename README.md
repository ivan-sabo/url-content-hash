# url-content-hash

This command-line tool calculates MD5 hash of the content of Web pages in a concurrent way.

## Parameters:

A list of URLs

## Optional -parallel flag:

A user can specify optional -parallel flag, which specifies how many concurrent requests are permitted. Default value is 10

## Usage:

go run main.go google.com facebook.com yahoo.com<br>
go run main.go -parallel 2 google.com facebook.com yahoo.com

## Implementation details

This tool implements pipeline architecture, meaning that the process is separated into 3 steps.<br><br>
First step: Get content - The Program will fetch the content of provided URLs in a concurrent way.<br>
Second step: Calculate hash - Hash will be calculated from obtained content.<br>
Third step: Print result - At the final stage the calculated hashes will be printed to the standard output.<br>

## Error handling

This program doesn't have any long-running parts; there are no consequences of error occurence, <br>
except the fact that hash won't be calculated in case the content wasn't obtained.<br>
Errors are included in Content and Hash structures, feel free to log/handle them as needed.
