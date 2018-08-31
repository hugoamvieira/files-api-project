# Files API

## Usage

- Run `dep ensure`
- Inside project folder, run `go run *.go`

---

## Endpoints
**Get an existing file's data**

`GET 127.0.0.1:8000/file?path="/a/path/with/a/file.txt"`

Will error if:
- Path doesn't start with `/` and end with `FILENAME.txt`;
- File not found;

**Delete existing file**

`DELETE 127.0.0.1:8000/file?path="/a/path/with/a/file.txt"`

Will error if:
- Path doesn't start with `/` and end with `FILENAME.txt`
- File not found;
- It can't delete the file for some reason;

**Create new file / Update existing file**

`POST 127.0.0.1:8000/file`

Request body must be JSON and follow the following format:
```
{
	"path": <PATH_TO_FILE> -> string (must end in `<FILE_NAME>.txt`)
	"data": <FILE_DATA> -> string
}
```

Will error if:
- Body doesn't follow above spec;
- It can't create the file for some reason;


**Get folder statistics**

`GET 127.0.0.1:8000/stats?path="/a/path`

Will error if:
- Path doesn't start with `/` or it has a filename at the end;
- Path not found;

Returns the following JSON format:
```
{
	"file_count": int,
	"alphanum_chars_percentage": {
		// key is filename, val is percentage (0-1)
		string: float,
		...
	},
	"alphanum_chars_stddev": float,
	"word_length_avg": float,
	"word_length_stddev": float,
	"bytes_total": int
}
```

A value of `-1` in any field denotes that the particular metric was invalid to calculate (eg: You requested statistics on a folder without files)

---

## Decision / Observations Log

- A very good way to have this API have as less I/O as possible (especially for statistics!) would be to maintain a tree-like structure in memory that reflected the filesystem and then update that tree based on any events on said filesystem - This way we can just check the tree in `O(logN)` time before we even hit the disk. I'm not going to implement this as I've got no time;

- There will be a root folder inside the project called `filesys` where everything will be created. This is so that everything can be contained within it;

- All files must have a `.txt` extension which must be specified to the relevant endpoints. This is so I can keep my sanity and avoid assuming that the last entry in the path is the filename;

- Every path must start with a `/`. This is for consistency and it also helps to avoid some errors;

- I've written the code for the endpoint that creates and writes the file in such a way that it also updates it if it already exists, so that saves some time. I'm not going to argue the correctness of using POST verb to update, though :P

- Even if some metrics are invalid to calculate (think requesting statistics on empty folder), it will still return with the appropriate fields set to `-1`. This is because there may still exist useful data to show, even if there's no files (eg: if you have a folder inside a folder, technically that inner folder occupies space and thus counts towards `bytes_total`).
